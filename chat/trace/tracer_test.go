package trace

import (
	"bytes"
	"testing"
)

func TestNG(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("fail")
	} else {
		tracer.Trace("こんにちわ")
		if buf.String() != "こんにちわ\n" {
			t.Errorf("Actual: %s", buf.String())
		}
	}
}