package rateg_test

import (
	"context"
	"testing"
	"time"

	"github.com/Xhofe/rateg"
)

func myFunction(a int) (int, error) {
	// do something
	return a + 1, nil
}

func TestLimitRate(t *testing.T) {
	myLimitedFunction := rateg.LimitRate(myFunction, time.Second)
	start := time.Now()
	result, err := myLimitedFunction(1)
	if err != nil {
		t.Fatal(err)
	}
	if result != 2 {
		t.Fatalf("expect 2, got %d", result)
	}
	result, err = myLimitedFunction(2)
	if err != nil {
		t.Fatal(err)
	}
	if result != 3 {
		t.Fatalf("expect 3, got %d", result)
	}
	elapsed := time.Since(start) // expect elapsed time greater than 1 second
	if elapsed < time.Second {
		t.Fatal("elapsed time less than 1 second")
	}
}

type Test struct {
	limitFn func(string) (string, error)
}

func (t *Test) myFunction(a string) (string, error) {
	// do something
	return a + " world", nil
}

func TestLimitRateStruct(t *testing.T) {
	test := &Test{}
	test.limitFn = rateg.LimitRate(test.myFunction, time.Second)
	start := time.Now()
	result, err := test.limitFn("hello")
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello world" {
		t.Fatalf("expect hello world, got %s", result)
	}
	result, err = test.limitFn("hi")
	if err != nil {
		t.Fatal(err)
	}
	if result != "hi world" {
		t.Fatalf("expect hi world, got %s", result)
	}
	elapsed := time.Since(start) // expect elapsed time greater than 1 second
	if elapsed < time.Second {
		t.Fatal("elapsed time less than 1 second")
	}
}

func myFunctionCtx(ctx context.Context, a int) (int, error) {
	// do something
	return a + 1, nil
}

func TestLimitRateCtx(t *testing.T) {
	myLimitedFunction := rateg.LimitRateCtx(myFunctionCtx, time.Second)
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
