package gameboard

import "math"

func getAbsInt16(val int16) int16 {
	if val < 0 {
		return val * -1
	}
	return val
}

func calculateHypotenuseFromTwoPoints(point1 point64, point2 point64) float64 {
	hypotenuse := math.Sqrt(math.Pow(point2.x-point1.x, 2) + (math.Pow(point2.y-point1.y, 2)))
	return hypotenuse
}

func calculateAngleRadiansUsingSin(opposite float64, hypotenuse float64) float64 {
	radians := math.Asin(opposite / hypotenuse)
	return radians
}

func radiansToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func calculateRotationDegrees(point1 point64, point2 point64) int16 {
	hypotenuse := calculateHypotenuseFromTwoPoints(point1, point2)
	opposite := point1.y - point2.y
	radians := calculateAngleRadiansUsingSin(opposite, hypotenuse)
	degrees64 := radiansToDegrees(radians)
	return int16(degrees64)
}

/*https://www.youtube.com/watch?v=XOk0aGwZYn8*/
func getRotation360(point1 point64, point2 point64) int16 {
	return int16(radiansToDegrees(math.Atan2(point1.y-point2.y, point1.x-point2.x)))
}

func getRotation360Float(point1 point64, point2 point64) float64 {
	return radiansToDegrees(math.Atan2(point1.y-point2.y, point1.x-point2.x))
}

func lengthOfLine(point1 point64, point2 point64) float64 {
	return math.Sqrt(math.Pow(float64(point1.x-point2.x), 2) + math.Pow(float64(point1.y-point2.y), 2))
}
func calculateAdjacentUsingCos(angle int16, hypotenuse int16) int16 {
	return int16(math.Cos(float64(angle)) * float64(hypotenuse))
}

func calculateOppositeUsingSin(angle int16, hypotenuse int16) int16 {
	return int16(math.Sin(float64(angle)) * float64(hypotenuse))
}
