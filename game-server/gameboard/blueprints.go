package gameboard

import (
	c2d "github.com/Tarliton/collision2d"
)

type blueprint struct {
	identifier       uint8
	polygonHitbox    c2d.Polygon
	circleHitbox     c2d.Circle
	hasPolygonHitbox bool
}

type weapon struct {
	damage      int
	projectiles []projectile
}

type projectileBlueprint struct {
	blueprint
	maxSpeed               float64
	damage                 int16
	existenceTime          uint16
	damageMultiplierLight  float64
	damageMultiplierMedium float64
	damageMultiplierHeavy  float64
}

type shipBlueprint struct {
	blueprint
	shipClass  string
	armorClass string
	maxSpeed   float64
	// how much the ship can rotate in a single step
	maxRotationSpeed           int16
	personalSpace              int16
	mass                       int16
	weaponCoolDownFrames       int16
	currentWeaponCoolDownFrame int16
	health                     int16
	creditsCost                int16
}

var shipBlueprints = make(map[pieceName]shipBlueprint)
var projectileBlueprints = make(map[pieceName]projectileBlueprint)
var identifierMap = make(map[pieceName]uint8)
var stringPieceNameMap = make(map[string]pieceName)

type pieceName string

const (
	// SHIPS
	capitalShipName pieceName = "CAPITAL_SHIP"
	interceptorName pieceName = "INTERCEPTOR"
	flackShipName   pieceName = "FLACK_SHIP"
	destroyerName   pieceName = "DESTROYER"
	// PROJECTILES
	laserName      pieceName = "LASER"
	bigLaserName   pieceName = "BIG_LASER"
	flackLaserName pieceName = "FLACK_LASER"
)

func init() {
	stringPieceNameMap["INTERCEPTOR"] = interceptorName
	stringPieceNameMap["CAPITAL_SHIP"] = capitalShipName
	stringPieceNameMap["FLACK_SHIP"] = flackShipName
	stringPieceNameMap["DESTROYER"] = destroyerName
	stringPieceNameMap["LASER"] = laserName
	stringPieceNameMap["BIG_LASER"] = bigLaserName
	stringPieceNameMap["FLACK_LASER"] = flackLaserName
	// SHIP BLUEPRINTS
	identifierMap[capitalShipName] = 0
	identifierMap[interceptorName] = 1
	identifierMap[destroyerName] = 2
	identifierMap[flackShipName] = 3
	identifierMap[flackLaserName] = 253
	identifierMap[bigLaserName] = 254
	identifierMap[laserName] = 255

	shipBlueprints[capitalShipName] = shipBlueprint{
		blueprint: blueprint{
			identifier:       identifierMap[capitalShipName],
			polygonHitbox:    c2d.NewPolygon(c2d.NewVector(0, 0), c2d.NewVector(0, 0), 0, []float64{-5, -5, -5, 5, 5, -5, 5, 5}),
			hasPolygonHitbox: true,
		},
		maxSpeed:                   1,
		maxRotationSpeed:           1,
		shipClass:                  "CAPITAL_SHIP",
		armorClass:                 "HEAVY",
		personalSpace:              60,
		mass:                       50,
		weaponCoolDownFrames:       60,
		currentWeaponCoolDownFrame: 0,
		health:                     1000,
		creditsCost:                1,
	}

	shipBlueprints[interceptorName] = shipBlueprint{
		blueprint: blueprint{
			identifier:       identifierMap[interceptorName],
			polygonHitbox:    c2d.NewPolygon(c2d.NewVector(0, 0), c2d.NewVector(0, 0), 0, []float64{-5, -5, -5, 5, 5, -5, 5, 5}),
			hasPolygonHitbox: true,
		},
		maxSpeed:                   2,
		maxRotationSpeed:           1,
		personalSpace:              20,
		shipClass:                  "INTERCEPTOR",
		armorClass:                 "LIGHT",
		mass:                       5,
		weaponCoolDownFrames:       60,
		currentWeaponCoolDownFrame: 0,
		health:                     50,
		creditsCost:                50,
	}

	shipBlueprints[destroyerName] = shipBlueprint{
		blueprint: blueprint{
			identifier:       identifierMap[destroyerName],
			polygonHitbox:    c2d.NewPolygon(c2d.NewVector(0, 0), c2d.NewVector(0, 0), 0, []float64{-5, -5, -5, 5, 5, -5, 5, 5}),
			hasPolygonHitbox: true,
		},
		maxSpeed:                   1,
		maxRotationSpeed:           1,
		personalSpace:              50,
		shipClass:                  "DESTROYER",
		armorClass:                 "MEDIUM",
		mass:                       5,
		weaponCoolDownFrames:       60,
		currentWeaponCoolDownFrame: 0,
		health:                     150,
		creditsCost:                200,
	}

	shipBlueprints[flackShipName] = shipBlueprint{
		blueprint: blueprint{
			identifier:       identifierMap[flackShipName],
			polygonHitbox:    c2d.NewPolygon(c2d.NewVector(0, 0), c2d.NewVector(0, 0), 0, []float64{-5, -5, -5, 5, 5, -5, 5, 5}),
			hasPolygonHitbox: true,
		},
		maxSpeed:                   1,
		maxRotationSpeed:           1,
		personalSpace:              40,
		shipClass:                  "FLACK_SHIP",
		armorClass:                 "HEAVY",
		mass:                       5,
		weaponCoolDownFrames:       60,
		currentWeaponCoolDownFrame: 0,
		health:                     200,
		creditsCost:                200,
	}

	// PROJECTILE BLUEPRINTS
	projectileBlueprints[laserName] = projectileBlueprint{
		blueprint: blueprint{
			identifier:       identifierMap[laserName],
			polygonHitbox:    c2d.NewPolygon(c2d.NewVector(0, 0), c2d.NewVector(0, 0), 0, []float64{-5, -5, -5, 5, 5, -5, 5, 5}),
			hasPolygonHitbox: true,
		},
		maxSpeed:               5,
		damage:                 10,
		existenceTime:          200,
		damageMultiplierLight:  2,
		damageMultiplierMedium: 2,
		damageMultiplierHeavy:  0.5,
	}

	projectileBlueprints[bigLaserName] = projectileBlueprint{
		blueprint: blueprint{
			identifier:       identifierMap[bigLaserName],
			polygonHitbox:    c2d.NewPolygon(c2d.NewVector(0, 0), c2d.NewVector(0, 0), 0, []float64{-5, -5, -5, 5, 5, -5, 5, 5}),
			hasPolygonHitbox: true,
		},
		maxSpeed:               5,
		damage:                 25,
		existenceTime:          200,
		damageMultiplierLight:  1,
		damageMultiplierMedium: 2,
		damageMultiplierHeavy:  4,
	}

	projectileBlueprints[flackLaserName] = projectileBlueprint{
		blueprint: blueprint{
			identifier:       identifierMap[flackLaserName],
			polygonHitbox:    c2d.NewPolygon(c2d.NewVector(0, 0), c2d.NewVector(0, 0), 0, []float64{-5, -5, -5, 5, 5, -5, 5, 5}),
			hasPolygonHitbox: true,
		},
		maxSpeed:               5,
		damage:                 20,
		existenceTime:          200,
		damageMultiplierLight:  10,
		damageMultiplierMedium: 0.5,
		damageMultiplierHeavy:  0.25,
	}
}
