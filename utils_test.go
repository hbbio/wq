package wq

import "testing"

func TestMin(t *testing.T) {
	var x int
	x = MinInt(-1, 1)
	if x != -1 {
		t.Error("min1")
	}
}
