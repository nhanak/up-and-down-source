package gameboard

import "log"

/***********************
 * GENERAL BOID RULES
 ***********************/

// Fly towards flock center of mass
func boidRule1Ship(gb *GameBoard, ship *ship, logger *log.Logger) velocityVector {
	//logger.Printf("Running boid rule 1 on pieceID: %v", ship.piece.pieceID)
	var centerOfMassX int16
	var centerOfMassXTop int16
	var centerOfMassXBottom int16
	var centerOfMassY int16
	var centerOfMassYTop int16
	var centerOfMassYBottom int16
	var numShipsConsideredInFlock int16
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if spShip.piece.pieceID == ship.piece.pieceID {
			// ignore self when trying to find center of mass

			continue
		}
		if (spShip.piece.player != ship.piece.player) || (spShip.maxSpeed != ship.maxSpeed) {
			// ignore those not in our flock
			continue
		}
		centerOfMassXTop += getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x) * spShip.mass
		centerOfMassYTop += getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y) * spShip.mass
		centerOfMassXBottom += spShip.mass
		centerOfMassYBottom += spShip.mass
		numShipsConsideredInFlock++
	}
	if numShipsConsideredInFlock == 0 {
		// there was nothing in the flock
		//logger.Printf("Returning because no ships were in flock...")
		return velocityVector{0, 0}
	}
	centerOfMassX = centerOfMassXTop / centerOfMassXBottom
	centerOfMassY = centerOfMassYTop / centerOfMassYBottom

	currentPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	currentPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	nextPositionX := (centerOfMassX - currentPositionX) / 100
	nextPositionY := (centerOfMassY - currentPositionY) / 100
	//logger.Printf("Center of Mass Next X: %v", nextPositionX)
	//logger.Printf("Center of Mass Next Y: %v", nextPositionY)
	return velocityVector{nextPositionX, nextPositionY}
}

// Keep distance from boids of different flock
func boidRule2Ship(gb *GameBoard, ship *ship, logger *log.Logger) velocityVector {
	vV := velocityVector{0, 0}
	shipXPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	shipYPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if (spShip.piece.pieceID != ship.piece.pieceID) && (spShip.piece.player != ship.piece.player) {
			spShipXPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			spShipYPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			absXDiff := getAbsInt16(spShipXPos - shipXPos)
			absYDiff := getAbsInt16(spShipYPos - shipYPos)
			if absXDiff < ship.personalSpace {
				vV.x = vV.x - (spShipXPos - shipXPos)
			}
			if absYDiff < ship.personalSpace {
				vV.y = vV.y - (spShipYPos - shipYPos)
			}
		}
	}
	return vV
}

// Try to match velocity of other boids
func boidRule3Ship(gb *GameBoard, ship *ship, logger *log.Logger) velocityVector {
	return velocityVector{0, 0}
}

// Encourage boid to fly within bounds of the gameboard
func boidRule4Ship(gb *GameBoard, ship *ship, logger *log.Logger) velocityVector {
	var xMin, xMax, yMin, yMax int16
	xMax = gb.width
	yMax = gb.height
	currentPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	currentPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	v := velocityVector{0, 0}
	if currentPositionX < xMin {
		v.x = 10
	}
	if currentPositionX > xMax {
		v.x = -10
	}
	if currentPositionY < yMin {
		v.y = 10
	}
	if currentPositionY > yMax {
		v.y = -10
	}
	return v
}

// Fly towards center of mass of other flock
func boidRule5Ship(gb *GameBoard, ship *ship, logger *log.Logger) velocityVector {
	//logger.Printf("Running boid rule 1 on pieceID: %v", ship.piece.pieceID)
	var centerOfMassX int16
	var centerOfMassXTop int16
	var centerOfMassXBottom int16
	var centerOfMassY int16
	var centerOfMassYTop int16
	var centerOfMassYBottom int16
	var numShipsConsideredInFlock int16
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if spShip.piece.pieceID == ship.piece.pieceID {
			// ignore self when trying to find center of mass
			continue
		}
		if spShip.piece.player == ship.piece.player {
			// ignore those in our flock
			continue
		}
		centerOfMassXTop += getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x) * spShip.mass
		centerOfMassYTop += getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y) * spShip.mass
		centerOfMassXBottom += spShip.mass
		centerOfMassYBottom += spShip.mass
		numShipsConsideredInFlock++
	}
	if numShipsConsideredInFlock == 0 {
		// there was nothing in the flock
		return velocityVector{0, 0}
	}
	centerOfMassX = centerOfMassXTop / centerOfMassXBottom
	centerOfMassY = centerOfMassYTop / centerOfMassYBottom

	currentPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	currentPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	nextPositionX := (centerOfMassX - currentPositionX) / 100
	nextPositionY := (centerOfMassY - currentPositionY) / 100
	return velocityVector{nextPositionX, nextPositionY}
}

// Keep distance from boids of same flock
func boidRule6Ship(gb *GameBoard, ship *ship, logger *log.Logger) velocityVector {
	vV := velocityVector{0, 0}
	shipXPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	shipYPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if (spShip.piece.pieceID != ship.piece.pieceID) && (spShip.piece.player == ship.piece.player) && (spShip.name == ship.name) {
			spShipXPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			spShipYPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			absXDiff := getAbsInt16(spShipXPos - shipXPos)
			absYDiff := getAbsInt16(spShipYPos - shipYPos)
			if absXDiff < ship.personalSpace {
				vV.x = vV.x - ((spShipXPos - shipXPos) / 10)
			}
			if absYDiff < ship.personalSpace {
				vV.y = vV.y - ((spShipYPos - shipYPos) / 10)
			}
		}
	}
	return vV
}

// Get in range of enemy with pieceName
func boidRule7(gb *GameBoard, ship *ship, name pieceName, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	var desiredRangeX int16
	desiredRangeX = 100
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if spShip.piece.name == name && spShip.piece.player != ship.player {
			myX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
			myY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
			enemyX := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			enemyY := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			absDifferenceX := getAbsInt16(myX - enemyX)
			if absDifferenceX < desiredRangeX {
				return v
			}
			v.x = (enemyX - myX) / 100
			v.y = (enemyY - myY) / 100
		}
	}
	return v
}

// Keep distance from boids of same flock with intensity (lower intensity value, more intense it is)
func boidRule8Ship(gb *GameBoard, ship *ship, intensity int16, logger *log.Logger) velocityVector {
	vV := velocityVector{0, 0}
	shipXPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	shipYPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if (spShip.piece.pieceID != ship.piece.pieceID) && (spShip.piece.player == ship.piece.player) && (spShip.name == ship.name) {
			spShipXPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			spShipYPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			absXDiff := getAbsInt16(spShipXPos - shipXPos)
			absYDiff := getAbsInt16(spShipYPos - shipYPos)
			if absXDiff < ship.personalSpace {
				vV.x = vV.x - ((spShipXPos - shipXPos) / intensity)
			}
			if absYDiff < ship.personalSpace {
				vV.y = vV.y - ((spShipYPos - shipYPos) / intensity)
			}
		}
	}
	return vV
}

// Keep distance from boids of different flock with intensity (lower intensity value, more intense it is)
func boidRule9Ship(gb *GameBoard, ship *ship, intensity int16, logger *log.Logger) velocityVector {
	vV := velocityVector{0, 0}
	shipXPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	shipYPos := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if (spShip.piece.pieceID != ship.piece.pieceID) && (spShip.piece.player != ship.piece.player) {
			spShipXPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			spShipYPos := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			absXDiff := getAbsInt16(spShipXPos - shipXPos)
			absYDiff := getAbsInt16(spShipYPos - shipYPos)
			if absXDiff < ship.personalSpace {
				vV.x = vV.x - ((spShipXPos - shipXPos) / intensity)
			}
			if absYDiff < ship.personalSpace {
				vV.y = vV.y - ((spShipYPos - shipYPos) / intensity)
			}
		}
	}
	return vV
}

// Get in range of nearest enemy ship with pieceName and intensity (lower intensity value, more intense it is)
func boidRule10(gb *GameBoard, thisShip *ship, name pieceName, intensity int16, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	var nearestInterceptor ship
	var distanceToNearestInterceptor float64
	numEnemyInterceptors := 0
	myX := getMostRecentInt16CurveKeyFrameValue(thisShip.positionCurve.x)
	myY := getMostRecentInt16CurveKeyFrameValue(thisShip.positionCurve.y)
	myPos := point64{float64(myX), float64(myY)}
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if spShip.piece.name == name && spShip.player != thisShip.player {
			enemyX := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			enemyY := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			posOfEnemyInterceptor := point64{float64(enemyX), float64(enemyY)}
			distanceToEnemyInterceptor := lengthOfLine(myPos, posOfEnemyInterceptor)
			if numEnemyInterceptors == 0 {
				distanceToNearestInterceptor = distanceToEnemyInterceptor
				nearestInterceptor = spShip
			} else {
				if distanceToEnemyInterceptor < distanceToNearestInterceptor {
					distanceToNearestInterceptor = distanceToEnemyInterceptor
					nearestInterceptor = spShip
				}
			}
			numEnemyInterceptors++
		}
	}
	if numEnemyInterceptors > 0 {
		enemyX := getMostRecentInt16CurveKeyFrameValue(nearestInterceptor.positionCurve.x)
		enemyY := getMostRecentInt16CurveKeyFrameValue(nearestInterceptor.positionCurve.y)
		v.x = (enemyX - myX) / intensity
		v.y = (enemyY - myY) / intensity
	}
	return v
}

/***********************
 * DESTROYER BOID RULES
 ***********************/

// Get in range of enemy destroyers
func destroyerBoidRule1(gb *GameBoard, ship *ship, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	var desiredRangeX int16
	desiredRangeX = 100
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if spShip.piece.name == capitalShipName && spShip.piece.player != ship.player {
			myX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
			myY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
			enemyX := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			enemyY := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			absDifferenceX := getAbsInt16(myX - enemyX)
			if absDifferenceX < desiredRangeX {
				return v
			}
			v.x = (enemyX - myX) / 100
			v.y = (enemyY - myY) / 100
		}
	}
	return v
}

/***********************
 * INTERCEPTOR BOID RULES
 ***********************/

// Get in range of nearest enemy interceptor
func interceptorBoidRule1(gb *GameBoard, thisShip *ship, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	var nearestInterceptor ship
	var distanceToNearestInterceptor float64
	numEnemyInterceptors := 0
	myX := getMostRecentInt16CurveKeyFrameValue(thisShip.positionCurve.x)
	myY := getMostRecentInt16CurveKeyFrameValue(thisShip.positionCurve.y)
	myPos := point64{float64(myX), float64(myY)}
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if spShip.piece.name == interceptorName && spShip.player != thisShip.player {
			enemyX := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			enemyY := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			posOfEnemyInterceptor := point64{float64(enemyX), float64(enemyY)}
			distanceToEnemyInterceptor := lengthOfLine(myPos, posOfEnemyInterceptor)
			if numEnemyInterceptors == 0 {
				distanceToNearestInterceptor = distanceToEnemyInterceptor
				nearestInterceptor = spShip
			} else {
				if distanceToEnemyInterceptor < distanceToNearestInterceptor {
					distanceToNearestInterceptor = distanceToEnemyInterceptor
					nearestInterceptor = spShip
				}
			}
			numEnemyInterceptors++
		}
	}
	if numEnemyInterceptors > 0 {
		enemyX := getMostRecentInt16CurveKeyFrameValue(nearestInterceptor.positionCurve.x)
		enemyY := getMostRecentInt16CurveKeyFrameValue(nearestInterceptor.positionCurve.y)
		v.x = (enemyX - myX) / 100
		v.y = (enemyY - myY) / 100
	}
	return v
}
