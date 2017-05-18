/*
	Copyright 2017 Brian Bauer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package dateslice

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

func TestRangeSlice(t *testing.T) {
	var rangeTest = []struct {
		BegDate           time.Time
		EndDate           time.Time
		ExpectedCount     int
		ExpectedFirstDate time.Time
		ExpectedLastDate  time.Time
	}{
		{time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), time.Date(2017, time.April, 3, 5, 0, 0, 0, time.UTC), 3, time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), time.Date(2017, time.April, 3, 5, 0, 0, 0, time.UTC)},
		{time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), 1, time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC)},
		{time.Date(2017, time.March, 1, 5, 0, 0, 0, time.UTC), time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), 32, time.Date(2017, time.March, 1, 5, 0, 0, 0, time.UTC), time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC)},
	}

	for _, ex := range rangeTest {
		ds := Range(ex.BegDate, ex.EndDate)
		if len(ds) != ex.ExpectedCount {
			t.Errorf("Unexpected number of elements in range test %d vs %d\n", ex.ExpectedCount, len(ds))
		}
		if ds[0].Year() != ex.ExpectedFirstDate.Year() {
			t.Errorf("Unexpected first year in range test %d vs %d\n", ex.ExpectedFirstDate.Year(), ds[0].Year())
		}
		if ds[len(ds)-1].Year() != ex.ExpectedLastDate.Year() {
			t.Errorf("Unexpected last year in range test %d vs %d\n", ex.ExpectedLastDate.Year(), ds[len(ds)-1].Year())
		}
		if ds[0].YearDay() != ex.ExpectedFirstDate.YearDay() {
			t.Errorf("Unexpected first day in range test %d vs %d\n", ex.ExpectedFirstDate.YearDay(), ds[0].YearDay())
		}
		if ds[len(ds)-1].YearDay() != ex.ExpectedLastDate.YearDay() {
			t.Errorf("Unexpected last day in range test %d vs %d\n", ex.ExpectedLastDate.YearDay(), ds[len(ds)-1].YearDay())
		}
	}
}
