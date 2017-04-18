package dateslice

import (
	"fmt"
	"math"
	"time"
)

/*
Sometimes you need a slice of dates.  Here are some functions that make that
	a little easier.
*/

/*
Today returns a slice containing a single element - the current time
*/
func Today() []time.Time {
	return []time.Time{time.Now()}
}

/*
Yesterday returns a slice containing a single element - yesterday
*/
func Yesterday() []time.Time {
	return []time.Time{time.Now().AddDate(0, 0, -1)}
}

/*
Tomorrow returns a slice containing a single element - tomorrow
*/
func Tomorrow() []time.Time {
	return []time.Time{time.Now().AddDate(0, 0, 1)}
}

func aWeek(baseDate time.Time) []time.Time {
	ds := make([]time.Time, 7)

	dow := baseDate.Weekday()

	// Reset the base date to Sunday
	baseDate = baseDate.AddDate(0, 0, 0-int(dow))

	for i := range ds {
		ds[i] = baseDate.AddDate(0, 0, i)
	}

	return ds
}

/*
WeekOf returns a slice containing all dates that occur during this specific week (Sunday is the first day of the week in Go!)
*/
func WeekOf(date time.Time) []time.Time {
	return aWeek(date)
}

/*
ThisWeek returns a slice containing all dates that occur this week (Sunday is the first day of the week in Go!)
*/
func ThisWeek() []time.Time {
	return aWeek(time.Now())
}

/*
LastWeek returns a slice containing all dates that occured last week (Sunday is the first day of the week in Go!)
*/
func LastWeek() []time.Time {
	return aWeek(time.Now().AddDate(0, 0, -7))
}

/*
NextWeek returns a slice containing all dates that will occur next week (Sunday is the first day of the week in Go!)
*/
func NextWeek() []time.Time {
	return aWeek(time.Now().AddDate(0, 0, 7))
}

func aMonth(baseDate time.Time) []time.Time {
	// This is used for subtraction, so the first day of the month needs to be a 0 instead of a 1
	dom := baseDate.Day() - 1

	//reset the base date to the 1st of the month
	baseDate = baseDate.AddDate(0, 0, 0-int(dom))

	firstOfNextMonth := baseDate.AddDate(0, 1, 0)
	daysInThisMonth := firstOfNextMonth.Sub(baseDate).Hours() / 24.0
	fmt.Printf("%f days in the month\n", math.Ceil(daysInThisMonth))

	ds := make([]time.Time, int(math.Ceil(daysInThisMonth)))

	for i := range ds {
		ds[i] = baseDate.AddDate(0, 0, i)
	}

	return ds
}

/*
ThisMonth returns a slice containing all dates that occur this month
*/
func ThisMonth() []time.Time {
	return aMonth(time.Now())
}

/*
LastMonth returns a slice containing all dates that occured last month
*/
func LastMonth() []time.Time {
	return aMonth(time.Now().AddDate(0, -1, 0))
}

/*
NextMonth returns a slice containing all dates that will occur next month
*/
func NextMonth() []time.Time {
	return aMonth(time.Now().AddDate(0, 1, 0))
}

func aYear(baseDate time.Time) []time.Time {
	// This is used for subtraction, so the first day of the month needs to be a 0 instead of a 1
	dom := baseDate.YearDay() - 1

	//reset the base date to the 1st of the month
	baseDate = baseDate.AddDate(0, 0, 0-int(dom))

	firstOfNextYear := baseDate.AddDate(1, 0, 0)
	daysInThisYear := firstOfNextYear.Sub(baseDate).Hours() / 24.0
	fmt.Printf("%f days in the year\n", math.Ceil(daysInThisYear))

	ds := make([]time.Time, int(math.Ceil(daysInThisYear)))

	for i := range ds {
		ds[i] = baseDate.AddDate(0, 0, i)
	}

	return ds
}
