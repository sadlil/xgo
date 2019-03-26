package xcontext

import (
	"context"
	"fmt"
)

type logTagContext struct {
	context.Context
	logTag string
	uuid   string
}

// WithLogTag returns a copy of ctx and appends [key:value] in the log tag. If uuid is
// available in ctx, it will be used as the prefix of the log tag.
// The soul purpose of this log tag is to carry around identifiers that can be used
// for logging and debugging in systems.
func WithLogTag(ctx context.Context, key, value string) context.Context {
	taggedCtx, ok := ctx.(*logTagContext)
	if !ok {
		taggedCtx = &logTagContext{
			Context: ctx,
		}
	}

	taggedCtx.logTag = fmt.Sprintf("%s[%s:%s]", taggedCtx.logTag, key, value)
	return taggedCtx
}

const (
	emptyString = ""
	uuidTag     = "uuid"
)

// LogTag returns the string representations of log tags to be used for logging.
func LogTag(ctx context.Context) string {
	taggedCtx, ok := ctx.(*logTagContext)
	if !ok {
		// if WithLogTag is not called before but uuid is set only
		// we return uuid as log tag
		uuid, _ := UUID(ctx)
		if len(uuid) > 0 {
			return fmt.Sprintf("[%s:%s]", uuidTag, uuid)
		}
		return emptyString
	}

	if taggedCtx.uuid == emptyString {
		// uuid was not set by any earlier LogTag call
		taggedCtx.uuid, _ = UUID(ctx)
	}

	var logTag = taggedCtx.logTag
	if taggedCtx.uuid != emptyString {
		logTag = fmt.Sprintf("[%s:%s]%s", uuidTag, taggedCtx.uuid, logTag)
	}
	return logTag
}
