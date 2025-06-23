package rateg_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/OpenListTeam/rateg"
)

func myFunction(a int) (int, error) {
	// do something
	return a + 1, nil
}

func TestLimitFn(t *testing.T) {
	myLimitedFunction := rateg.LimitFn(myFunction, rateg.LimitFnOption{
		Limit:  1,
		Bucket: 1,
	})
	start := time.Now()
	wg := sync.WaitGroup{}
	count := 3
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			result, err := myLimitedFunction(1)
			if err != nil {
				t.Fatal(err)
			}
			if result != 2 {
				t.Fatalf("expect 2, got %d", result)
			} else {
				t.Logf("fn called succeed")
			}
			result, err = myLimitedFunction(2)
			if err != nil {
				t.Fatal(err)
			}
			if result != 3 {
				t.Fatalf("expect 3, got %d", result)
			} else {
				t.Logf("fn called succeed")
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	expected := time.Second * time.Duration(count*2-1)
	if elapsed < expected {
		t.Fatalf("elapsed time expected: > %s, actual: %s", expected, elapsed)
	}
}

type Test struct {
	limitFn func(string) (string, error)
}

func (t *Test) myFunction(a string) (string, error) {
	// do something
	return a + " world", nil
}

func TestLimitFnStruct(t *testing.T) {
	test := &Test{}
	test.limitFn = rateg.LimitFn(test.myFunction, rateg.LimitFnOption{
		Limit:  1,
		Bucket: 1,
	})
	start := time.Now()
	wg := sync.WaitGroup{}
	count := 3
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			result, err := test.limitFn("hello")
			if err != nil {
				t.Fatal(err)
			}
			if result != "hello world" {
				t.Fatalf("expect hello world, got %s", result)
			} else {
				t.Logf("fn called succeed")
			}
			result, err = test.limitFn("hi")
			if err != nil {
				t.Fatal(err)
			}
			if result != "hi world" {
				t.Fatalf("expect hi world, got %s", result)
			} else {
				t.Logf("fn called succeed")
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	expected := time.Second * time.Duration(count*2-1)
	if elapsed < expected {
		t.Fatalf("elapsed time expected: > %s, actual: %s", expected, elapsed)
	}
}

func myFunctionCtx(ctx context.Context, a int) (int, error) {
	// do something
	return a + 1, nil
}

func TestLimitFnCtx(t *testing.T) {
	myLimitedFunction := rateg.LimitFnCtx(myFunctionCtx, rateg.LimitFnOption{
		Limit:  1,
		Bucket: 1,
	})
	result, err := myLimitedFunction(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result != 2 {
		t.Fatalf("expect 2, got %d", result)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	result, err = myLimitedFunction(ctx, 2)
	if err != context.Canceled {
		t.Fatalf("expect context.Canceled, got %v", err)
	}
}
