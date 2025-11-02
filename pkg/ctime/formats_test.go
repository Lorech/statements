package ctime_test

import (
	"statements/pkg/ctime"
	"testing"
	"time"
)

type DateTest struct {
	date       string // input string
	yx, mx, dx int    // expected results
}

var dateTests = []DateTest{
	{"03.06.2025", 2025, 6, 3},
	{"31.10.2025", 2025, 10, 31},
}

func TestLittleEndianDateOnly(t *testing.T) {
	for _, dt := range dateTests {
		tp, err := time.Parse(ctime.LittleEndianDateOnly, dt.date)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		y := tp.Year()
		m := tp.Month()
		d := tp.Day()

		if y != dt.yx || m != time.Month(dt.mx) || d != dt.dx {
			t.Errorf("got %d-%d-%d; expected %d-%d-%d", y, m, d, dt.yx, dt.mx, dt.dx)
		}
	}
}
