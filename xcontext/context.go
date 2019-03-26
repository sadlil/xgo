package xcontext

import (
	"context"
)

var (
	// Empty holds an instance of an empty Background context
	Empty = context.Background()
)

type orphanedContext struct {
	context.Context
	detachedParent context.Context
}

// Orphaned returns a copy of the parent without the parent relation.
// Any values passed on parent would be available on Orphaned context
// but it will not be cancelled when the parent context is cancelled.
// Any future update operation will not modify parent.
func Orphaned(parent context.Context) context.Context {
	ctx := context.Background()
	return &orphanedContext{
		Context:        ctx,
		detachedParent: parent,
	}
}

func (c *orphanedContext) Value(key interface{}) interface{} {
	var value = c.Context.Value(key)
	if c.detachedParent != nil && value == nil {
		// if context have a detached parent and the key is not found
		// in the child context itself, we should query its parent for
		// the key to see if the parent have it.
		value = c.detachedParent.Value(key)
	}
	return value
}
