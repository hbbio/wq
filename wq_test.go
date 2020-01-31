package wq

import (
	"testing"
)

func dummyFn(s string) error {
	return nil
}

type testT struct {
	fn        Payload
	nbWorkers int
	list      []string
}

var tests = []testT{
	{dummyFn, 1, []string{}},
	{dummyFn, 1, []string{"foo"}},
	{dummyFn, 1, []string{"foo", "bar"}},
	{dummyFn, 3, []string{}},
	{dummyFn, 3, []string{"foo"}},
	{dummyFn, 3, []string{"foo", "bar"}},
	{dummyFn, 3, []string{"foo", "bar", "foo", "bar"}},
	{dummyFn, -1, []string{}},
	{dummyFn, -1, []string{"foo"}},
}

func TestQueue(t *testing.T) {
	for i, v := range tests {
		res, err := Queue(v.nbWorkers, v.fn, v.list)
		if err == ErrNoJob {
			if !(v.nbWorkers <= 0 || len(v.list) == 0) {
				t.Log(v)
				t.Error("should not fail")
			}
		} else if err != nil {
			t.Error(err)
		}
		if err == nil {
			if res != len(v.list) {
				t.Error("queue", i)
			}
		}
	}
}
