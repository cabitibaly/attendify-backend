package utils

import "math"

const RAYON_TERRESTRE = 6371000

func ConvDegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func HaversineDistance(
	empLatitude,
	empLongitude,
	siteLatitude,
	siteLongitude float64,
) float64 {

	empLatRad := ConvDegreesToRadians(empLatitude)
	empLongRad := ConvDegreesToRadians(empLongitude)
	siteLatRad := ConvDegreesToRadians(siteLatitude)
	siteLongRad := ConvDegreesToRadians(siteLongitude)

	deltaLat := siteLatRad - empLatRad
	deltaLong := siteLongRad - empLongRad

	haversineFormule := math.Pow(math.Sin(deltaLat/2), 2) +
		math.Cos(empLatRad)*math.Cos(siteLatRad)*
			math.Pow(math.Sin(deltaLong/2), 2)

	angleCentral := 2 * math.Atan2(math.Sqrt(haversineFormule), math.Sqrt(1-haversineFormule))

	distance := RAYON_TERRESTRE * angleCentral

	return distance
}

func EstDansLaZone(
	empLatitude,
	empLongitude,
	siteLatitude,
	siteLongitude,
	rayon float64,
) bool {

	distance := HaversineDistance(
		empLatitude,
		empLongitude,
		siteLatitude,
		siteLongitude,
	)

	return distance <= rayon
}
