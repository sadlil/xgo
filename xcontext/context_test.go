package xcontext

import (
	"context"
	"testing"
	"time"
)

func TestOrphanedValue(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")

	orphaned := Orphaned(ctx)
	value := orphaned.Value("key1")
	if value == nil {
		t.Errorf("Expected value, got nil")
	}
	if "value1" != value.(string) {
		t.Errorf("Expected value1, got %s", value.(string))
	}

	// It should Update the value now keeping the parent immutable
	orphaned = context.WithValue(orphaned, "key1", "value2")
	value = orphaned.Value("key1")
	if value == nil {
		t.Errorf("Expected value, got nil")
	}
	if "value2" != value.(string) {
		t.Errorf("Expected value1, got %s", value.(string))
	}

	// Setting new value will effect the child only
	orphaned = context.WithValue(orphaned, "key2", "value2")
	value = orphaned.Value("key2")
	if value == nil {
		t.Errorf("Expected value, got nil")
	}
	if "value2" != value.(string) {
		t.Errorf("Expected value1, got %s", value.(string))
	}
}

func TestOrphanedValueParentCanceled(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "key1", "value1")

	orphaned := Orphaned(ctx)
	value := orphaned.Value("key1")
	if value == nil {
		t.Errorf("Expected value, got nil")
	}
	if "value1" != value.(string) {
		t.Errorf("Expected value1, got %s", value.(string))
	}

	// Parent canceled should still return values
	cancelFunc()
	value = orphaned.Value("key1")
	if value == nil {
		t.Errorf("Expected value, got nil")
	}
	if "value1" != value.(string) {
		t.Errorf("Expected value1, got %s", value.(string))
	}
}

func TestOrphanedNoCancelRelation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	orphaned := Orphaned(parent)

	cancel()
	if parent.Err() == nil {
		t.Errorf("Expected context canceled")
	}

	if orphaned.Err() != nil {
		t.Errorf("Expected context not canceled")
	}
}

func TestOrphanedNoDeadlineRelation(t *testing.T) {
	parent, cancelGoCtx := context.WithDeadline(
		context.Background(),
		time.Now().Add(50*time.Nanosecond),
	)
	defer cancelGoCtx()

	orphaned := Orphaned(parent)

	// wait for context to deadline
	<-parent.Done()

	if parent.Err() == nil {
		t.Errorf("Expected context canceled")
	}

	if orphaned.Err() != nil {
		t.Errorf("Expected context not canceled")
	}
}
