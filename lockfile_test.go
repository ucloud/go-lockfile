package lockfile

import (
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	lf := New("test.lock")

	err := lf.TryLock()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	err = lf.Unlock()
	if err != nil {
		t.Fatal(err)
	}
}
