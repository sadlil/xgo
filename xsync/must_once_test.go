package xsync

import (
	"testing"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, once *MustOnce, o *one, c chan bool) {
	once.Do(func() error { o.Increment(); return nil })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestMustOnce(t *testing.T) {
	o := new(one)
	once := new(MustOnce)
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, once, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestMustOncePanic(t *testing.T) {
	var once MustOnce
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("MustOnce.Do did not panic")
			}
		}()
		once.Do(func() error {
			panic("failed")
			return nil
		})
	}()

	once.Do(func() error {
		t.Fatalf("MustOnce.Do called twice")
		return nil
	})
}

func BenchmarkMustOnce(b *testing.B) {
	var once MustOnce
	f := func() error { return nil }
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do(f)
		}
	})
}
