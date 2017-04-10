package main

import (
	"testing"
	"time"
)

func TestWeekSlice(t *testing.T) {
	var weekTest = []struct {
		Date              time.Time
		ExpectedCount     int
		ExpectedFirstDate time.Time
	}{
		{time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), 7, time.Date(2017, time.March, 26, 5, 0, 0, 0, time.UTC)},
		{time.Date(2016, time.February, 29, 5, 0, 0, 0, time.UTC), 7, time.Date(2016, time.February, 28, 5, 0, 0, 0, time.UTC)},
		{time.Date(2003, time.January, 3, 5, 0, 0, 0, time.UTC), 7, time.Date(2002, time.December, 29, 5, 0, 0, 0, time.UTC)},
	}

	for _, ex := range weekTest {
		ds := WeekOf(ex.Date)
		if len(ds) != ex.ExpectedCount {
			t.Errorf("Unexpected number of elements in weekly test %d vs %d\n", ex.ExpectedCount, len(ds))
		}
		if ds[0].Year() != ex.ExpectedFirstDate.Year() {
			t.Errorf("Unexpected year in weekly test %d vs %d\n", ex.ExpectedFirstDate.Year(), ds[0].Year())
		}
		if ds[0].YearDay() != ex.ExpectedFirstDate.YearDay() {
			t.Errorf("Unexpected day of year in weekly test %d vs %d\n", ex.ExpectedFirstDate.YearDay(), ds[0].YearDay())
		}
	}
}
