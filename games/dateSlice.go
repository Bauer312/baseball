package main

import "time"

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
