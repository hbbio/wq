// (c) Henri Binsztok, 2015
// See: LICENSE

package wq

import (
	"database/sql"
	"testing"
)

var dummyDb *sql.DB

func TestMin(t *testing.T) {
	var x int
	x = Min(-1, 1)
	if x != -1 {
		t.Error("min1")
	}
}

func dummyFn(db *sql.DB, s string) error {
	return nil
}

type testT struct {
	fn        payload
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
	waitBeforeStart = false
	for i, v := range tests {
		res := Queue(dummyDb, v.fn, v.nbWorkers, v.list)
		if res != len(v.list) && !(v.nbWorkers <= 0 && res == 0) {
			t.Error("queue", i)
		}
	}
}
