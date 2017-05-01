// (c) Henri Binsztok, 2015
// See: LICENSE

package wq

import (
	"testing"
	"database/sql"
)

var dummy_db *sql.DB

func TestMin(t *testing.T) {
	var x int
	x = Min(-1, 1)
	if x != -1 {
		t.Error("min1")
	}
}

func dummy_fn(db *sql.DB, s string) error {
	return nil
}

type testT struct {
	fn payload
	nbWorkers int
	list []string
}

var tests = []testT {
	{ dummy_fn, 1, []string{} },
	{ dummy_fn, 1, []string{"foo"} },
	{ dummy_fn, 1, []string{"foo", "bar"} },
	{ dummy_fn, 3, []string{} },
	{ dummy_fn, 3, []string{"foo"} },
	{ dummy_fn, 3, []string{"foo", "bar"} },
	{ dummy_fn, 3, []string{"foo", "bar", "foo", "bar"} },
	{ dummy_fn, -1, []string{} },
	{ dummy_fn, -1, []string{"foo"} },
}

func TestQueue(t *testing.T) {
	waitBeforeStart = false
	for i, v := range tests {
		res := Queue(dummy_db, v.fn, v.nbWorkers, v.list)
		if res != len(v.list) && !(v.nbWorkers <= 0 && res == 0) {
			t.Error("queue", i)
		}
	}
}