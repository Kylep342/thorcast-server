/*
 *
 *
 *
 */

package main

// Location corresponds to a row in the geocodex table
// The only fields that are read/written by the app are below
type Location struct {
	City  string  `db:"city"`
	State string  `db:"state"`
	Lat   float64 `db:"lat"`
	Lng   float64 `db:"lng"`
}

// SetCoords will set the Lat and Lng values for a Location
func (l *Location) SetCoords(o Coordinates) {
	l.Lat = o.Lat
	l.Lng = o.Lng
}
