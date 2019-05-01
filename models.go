/*
 *
 *
 *
 */

package main

type Location struct {
	City  string  `db:"city"`
	State string  `db:"state"`
	Lat   float64 `db:"lat"`
	Lng   float64 `db:"lng"`
}

func (l *Location) SetCoords(o Coordinates) {
	l.Lat = o.Lat
	l.Lng = o.Lng
}
