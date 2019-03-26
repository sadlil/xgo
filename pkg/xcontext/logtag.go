package xcontext

import (
	"context"
	"fmt"
)

type logTagContext struct{
	logTag string
	uuid string
	context.Context
}

// WithLogTag returns a copy of ctx and appends [key:value] in the log tag. If uuid is
// available beforehand in ctx, it will be used as the prefix of the log tag.
// WIthUUID is expected to be called before WithLogTag if UUID is expected in Log Tag.
// The soul purpose of this log tag is to carry around identifiers that can be used
// for logging and debugging in systems.
func WithLogTag(ctx context.Context, key, value string) context.Context {
	v := ctx.Value(contextKeyLogTag)
	if v == nil {
		uuid, _ := UUID(ctx)
		v = ""
		if len(uuid) > 0 {
			// if uuid is available set the uuid as the first element
			v = fmt.Sprintf("[uuid:%s]", uuid)
		}
	}

	logTag := v.(string)
	logTag = fmt.Sprintf("%s[%s:%s]", logTag, key, value)
	return context.WithValue(ctx, contextKeyLogTag, logTag)
}

// LogTag returns the string representations of log tags to be used for logging.
func LogTag(ctx context.Context) string {
	v := ctx.Value(contextKeyLogTag)
	if v == nil {
		uuid, _ := UUID(ctx)
		v = ""
		if len(uuid) > 0 {
			// if uuid is available set the uuid as the first element
			return fmt.Sprintf("[uuid:%s]", uuid)
		}
		return ""
	}

	logTag := v.(string)
	return logTag
}
