package rateg_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/Xhofe/rateg"
)

func TestCopy(t *testing.T) {
	src := bytes.NewBufferString("hello world")
	start := time.Now()
	dst := &bytes.Buffer{}
	_, err := rateg.Copy(dst, src, 5)
	if err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(start) // expect elapsed time greater than 3.x seconds
	if elapsed < 3*time.Second {
		t.Fatal("elapsed time less than 3 second")
	}
	if elapsed > 4*time.Second {
		t.Fatal("elapsed time greater than 4 second")
	}
	t.Log("elapsed time:", elapsed)
}
