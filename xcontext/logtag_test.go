package xcontext

import (
	"context"
	"sync"
	"testing"
)

func TestLogTag(t *testing.T) {
	tests := []struct {
		Name        string
		Setup       func() context.Context
		ExpectedTag string
	}{
		{
			Name: "Happy Path - Empty Context",
			Setup: func() context.Context {
				return context.Background()
			},
			ExpectedTag: "",
		},
		{
			Name: "Happy Path - LogTag Set",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithLogTag(ctx, "key", "value")
				return ctx
			},
			ExpectedTag: "[key:value]",
		},
		{
			Name: "Happy Path - With UUID Set",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithUUID(ctx, "uuid")
				ctx = WithLogTag(ctx, "key", "value")
				return ctx
			},
			ExpectedTag: "[uuid:uuid][key:value]",
		},
		{
			Name: "Happy Path - LogTag Set Multiple times",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithLogTag(ctx, "key", "value")
				ctx = WithLogTag(ctx, "key", "value")
				ctx = WithLogTag(ctx, "key", "value")
				return ctx
			},
			ExpectedTag: "[key:value][key:value][key:value]",
		},
		{
			Name: "Happy Path - With UUID Set",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithUUID(ctx, "uuid")
				ctx = WithUUID(ctx, "uuid2")
				ctx = WithLogTag(ctx, "key", "value")
				ctx = WithLogTag(ctx, "key", "value")
				ctx = WithLogTag(ctx, "key", "value")
				return ctx
			},
			ExpectedTag: "[uuid:uuid][key:value][key:value][key:value]",
		},
		{
			Name: "No LogTag Set",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithUUID(ctx, "uuid")
				return ctx
			},
			ExpectedTag: "[uuid:uuid]",
		},
		{
			Name: "WithValue called after WithLogTag",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithUUID(ctx, "uuid")
				ctx = WithLogTag(ctx, "key", "value")
				ctx = context.WithValue(ctx, "new-key", "new-value")
				ctx = context.WithValue(ctx, "new-key", "new-value")
				return ctx
			},
			ExpectedTag: "[uuid:uuid][key:value]",
		},
		{
			Name: "Child have no effect to parent",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithUUID(ctx, "uuid")
				ctx = WithLogTag(ctx, "key", "value")
				ctx = context.WithValue(ctx, "new-key", "new-value")

				var wg sync.WaitGroup
				go func(c context.Context) {
					wg.Add(1)
					defer wg.Done()

					c = WithLogTag(c, "child", "routine")

				}(ctx)
				wg.Wait()

				return ctx
			},
			ExpectedTag: "[uuid:uuid][key:value]",
		},
	}

	for _, test := range tests {
		ctx := test.Setup()

		if test.ExpectedTag != LogTag(ctx) {
			t.Errorf("%s: Expected %s, Got %s", test.Name, test.ExpectedTag, LogTag(ctx))
		}
	}
}
