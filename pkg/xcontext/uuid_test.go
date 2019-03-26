package xcontext

import (
	"context"
	"testing"
)

func TestUUID(t *testing.T) {
	testes := []struct {
		Name               string
		Setup              func() context.Context
		ExpectedUUID       string
		ExpectedUUIDLength int
		ExpectedErr        error
	}{
		{
			Name: "Happy Path",
			Setup: func() context.Context {
				ctx := context.Background()
				return WithUUID(ctx, "uuid")
			},
			ExpectedUUID:       "uuid",
			ExpectedUUIDLength: 4,
			ExpectedErr:        nil,
		},
		{
			Name: "UUID Overwrite",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithUUID(ctx, "uuid")
				ctx = WithUUID(ctx, "another-uuid")
				ctx = WithUUID(ctx, "other-uuid")
				return ctx
			},
			ExpectedUUID:       "uuid",
			ExpectedUUIDLength: 4,
			ExpectedErr:        nil,
		},
		{
			Name: "No UUID Set",
			Setup: func() context.Context {
				return context.Background()
			},
			ExpectedErr: ErrNoUUIDSet,
		},
		{
			Name: "UUID Not Provided",
			Setup: func() context.Context {
				ctx := context.Background()
				ctx = WithUUID(ctx, "")
				return ctx
			},
			ExpectedUUIDLength: 36,
			ExpectedErr:        nil,
		},
	}

	for _, test := range testes {
		ctx := test.Setup()
		uuid, err := UUID(ctx)
		if err != test.ExpectedErr {
			t.Errorf("%s: Expected %v, Got %v", test.Name, test.ExpectedErr, err)
		}

		if len(test.ExpectedUUID) > 0 {
			if uuid != test.ExpectedUUID {
				t.Errorf("%s: Expected %v, Got %v", test.Name, test.ExpectedUUID, uuid)
			}
		}

		if len(uuid) != test.ExpectedUUIDLength {
			t.Errorf("%s: Expected %v, Got %v", test.Name, test.ExpectedUUIDLength, len(uuid))
		}
	}

}

func TestNextUUID(t *testing.T) {
	uuid := nextUUID()
	if uuid == "" {
		t.Errorf("expected not empty, got empty")
	}

	var idMap = map[string]struct{}{}
	for i := 0; i < 1000; i++ {
		uuid = nextUUID()
		if _, ok := idMap[uuid]; ok {
			t.Errorf("uuid %s found duplicated", uuid)
		}

		if len(uuid) == 0 {
			t.Errorf("uuid length is zero")
		}

		idMap[uuid] = struct{}{}
	}
}

func BenchmarkNextUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nextUUID()
	}

	b.StopTimer()
	b.ReportAllocs()
}
