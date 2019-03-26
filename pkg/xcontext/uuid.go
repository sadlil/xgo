package xcontext

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

type keyUUID struct{}

var (
	// ErrNoUUIDSet is returned when context doesn't have a uuid set to it,
	// but xcontext.UUID is called.
	ErrNoUUIDSet = fmt.Errorf("no uuid is set for context")

	contextKeyUUID = keyUUID{}
)

// WithUUID returns a copy of ctx in which the contextKeyUUID is set for the context.
// WithUUID only sets the contextKeyUUID with the specified value only if the
// contextKeyUUID, ignores resetting the value otherwise.
// UUID is a efficient way to track a context by the id through out its lifespan.
// For a request scoped context it is much easier to track the request scope with this.
// WIthUUID is expected to be called before WithLogTag.
func WithUUID(ctx context.Context, uuid string) context.Context {
	if len(uuid) == 0 {
		uuid = nextUUID()
	}

	id := ctx.Value(contextKeyUUID)
	if id != nil {
		return ctx
	}
	return context.WithValue(ctx, contextKeyUUID, uuid)
}

// UUID returns the value set by WithUUID from context. It returns ErrNoUUIDSet
// if there is no value set with the contextKeyUUID key.
func UUID(ctx context.Context) (string, error) {
	id := ctx.Value(contextKeyUUID)
	if id != nil {
		return id.(string), nil
	}
	return "", ErrNoUUIDSet
}

// nextUUID returns a V4 uuid
func nextUUID() string {
	var uuid [16]byte

	const (
		// ensures we backOff for less than 450ms total. Use the following to
		// select new value, in units of 10ms:
		// 	n*(n+1)/2 = d -> n^2 + n - 2d -> n = (sqrt(8d + 1) - 1)/2
		maxRetries = 9
		backOff    = time.Millisecond * 10
	)

	var (
		retries int
	)

	for {
		// This should never block but the read may fail. Because of this,
		// we just try to read the random number generator until we get
		// something. This is a very rare condition but may happen.
		b := time.Duration(retries) * backOff
		time.Sleep(b)

		_, err := rand.Read(uuid[:])
		if err != nil {
			if retries < maxRetries {
				retries++
				continue
			}

			// Any other errors represent a system problem. What did someone
			// do to /dev/urandom?
			break
		}

		// Set the two most significant bits (bits 6 and 7) of the
		// clock_seq_hi_and_reserved to zero and one, respectively.
		uuid[8] = (uuid[8] | 0x40) & 0x7F

		// Set the four most significant bits (bits 12 through 15) of the
		// time_hi_and_version field to the 4-bit version number.
		uuid[6] = (uuid[6] & 0xF) | (byte(4) << 4)

		return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	}
	return ""
}
