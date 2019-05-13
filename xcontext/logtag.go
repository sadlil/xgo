package xcontext

import (
	"context"
	"fmt"
)

type logTagContext struct {
	context.Context
	logTag string
}

type keyLogTag struct{}

var (
	contextKeyLogTag = keyLogTag{}
)

// Value returns the value from logTagContext. If the provided key is
// contextKeyLogTag it will return the logTag
func (c *logTagContext) Value(key interface{}) interface{} {
	if key == contextKeyLogTag {
		return c.logTag
	}
	return c.Context.Value(key)
}

// WithLogTag returns a copy of ctx and appends [key:value] in the log tag.
// The soul purpose of this log tag is to carry around identifiers that can be used
// for logging and debugging in systems.
func WithLogTag(ctx context.Context, key, value string) context.Context {
	val := ctx.Value(contextKeyLogTag)
	if val == nil {
		val = ""
	}

	logTag := val.(string)
	return &logTagContext{
		logTag:  fmt.Sprintf("%s[%s:%s]", logTag, key, value),
		Context: ctx,
	}
}

const (
	emptyString = ""
	uuidTag     = "uuid"
)

// LogTag returns the string representations of log tags to be used for logging.
func LogTag(ctx context.Context) string {
	val := ctx.Value(contextKeyLogTag)
	if val == nil {
		val = ""
	}
	logTag := val.(string)
	return logTagWithUUID(ctx, logTag)
}

func logTagWithUUID(ctx context.Context, logTag string) string {
	uuid, _ := UUID(ctx)
	if uuid != emptyString {
		return fmt.Sprintf("[%s:%s]%s", uuidTag, uuid, logTag)
	}
	return logTag
}
