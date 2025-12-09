package utils

import "time"

func CalculerNombreJours(debutJour, finJour time.Time) int {
	var nombreJours int

	for debutJour.Before(finJour) {
		if debutJour.Weekday() != time.Sunday {
			nombreJours++
		}
		debutJour = debutJour.AddDate(0, 0, 1)
	}

	return nombreJours
}
