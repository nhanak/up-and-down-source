package gameboard

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"

	"strconv"

	c2d "github.com/Tarliton/collision2d"
)

// Static defenses, like turrets
type structure struct {
	piece
}

// Fields like repulsion fields, healing fields, damage fields etc
type field struct {
	piece
}

type position struct {
	x int16
	y int16
}

type point struct {
	x int16
	y int16
}

type point64 struct {
	x float64
	y float64
}

type velocityVector struct {
	x int16
	y int16
}

type piece struct {
	positionCurve    positionCurve
	rotationCurve    []int16CurveKeyFrame
	player           uint8
	name             pieceName
	pieceID          uint16
	existenceCurve   []uint8CurveKeyFrame
	polygonHitbox    c2d.Polygon
	circleHitbox     c2d.Circle
	hasPolygonHitbox bool
	identifier       uint8
}

type projectile struct {
	piece
	angleShot              float64
	damage                 int16
	maxSpeed               float64
	existenceTime          uint16
	timeExisted            uint16
	damageMultiplierLight  float64
	damageMultiplierMedium float64
	damageMultiplierHeavy  float64
}

type newProjectileInfo struct {
	position position
	name     pieceName
	angle    float64
	pieceID  uint16
	time     uint16
	player   uint8
}

type ship struct {
	piece
	shipClass                  string
	armorClass                 string
	velocity                   velocityVector
	maxSpeed                   float64
	maxRotationSpeed           int16
	personalSpace              int16
	mass                       int16
	weaponCoolDownFrames       int16
	currentWeaponCoolDownFrame int16
	health                     int16
}

type newShipInfo struct {
	position position
	name     pieceName
	player   uint8
	rotation int16
	pieceID  uint16
	time     uint16
}

// GameBoard the gameboard that spacebattles takes place on
type GameBoard struct {
	width                           int16
	height                          int16
	ships                           []ship
	projectiles                     []projectile
	fields                          []field
	structures                      []structure
	GameBoardMu                     sync.Mutex
	totalPieceCounter               uint16
	player1Credits                  int16
	player1Name                     string
	player2Credits                  int16
	player1Health                   int16
	player2Health                   int16
	player2Name                     string
	currentFrame                    uint16
	currentGameFrameTimeStamp       uint16
	lastSentGameFrameTimeStamp      uint16
	creditsForCapitalShipDamageBool bool
	numCreditsForCapitalShipDamage  int16
	creditsForShipDeathBool         bool
	refundPercentageForShipDeath    float64
}

func (gb *GameBoard) SetCurrentGameFrameTimeStamp(timeStamp uint16) {
	gb.currentGameFrameTimeStamp = timeStamp
}

func (gb *GameBoard) SetLastSentGameFrameTimeStamp(timeStamp uint16) {
	gb.lastSentGameFrameTimeStamp = timeStamp
}

func (gb *GameBoard) GetCurrentGameFrameTimeStamp() uint16 {
	return gb.currentGameFrameTimeStamp
}

func (gb *GameBoard) GetLastSentGameFrameTimeStamp() uint16 {
	return gb.lastSentGameFrameTimeStamp
}

func (gb *GameBoard) ForceGameOver() {
	gb.player1Health = 0
	gb.player2Health = 0
}

func (gb *GameBoard) ResetGameBoard() {
	gb.ships = make([]ship, 0)
	gb.projectiles = make([]projectile, 0)
	gb.fields = make([]field, 0)
	gb.structures = make([]structure, 0)
	gb.totalPieceCounter = 6
	gb.player1Credits = 0
	gb.player1Health = 1000
	gb.player2Health = 1000
	gb.player1Credits = 0
	gb.player2Credits = 0
	gb.setInitialShips()
}

// GetWidth get the width of the GameBoard
func (gb *GameBoard) GetWidth() int16 {
	return gb.width
}

func (gb *GameBoard) GetPlayer1Name() string {
	return gb.player1Name
}

func (gb *GameBoard) GetPlayer2Name() string {
	return gb.player2Name
}

func (gb *GameBoard) GetPlayer1Credits() int16 {
	return gb.player1Credits
}

func (gb *GameBoard) GetPlayer2Credits() int16 {
	return gb.player2Credits
}

func (gb *GameBoard) GetPlayer1Health() int16 {
	return gb.player1Health
}

func (gb *GameBoard) GetPlayer2Health() int16 {
	return gb.player2Health
}

// GetHeight get the height of the GameBoard
func (gb *GameBoard) GetHeight() int16 {
	return gb.height
}

func (gb *GameBoard) IsGameOver() bool {
	gameOver := false
	if (gb.player1Health <= 0) || (gb.player2Health <= 0) {
		gameOver = true
	}
	return gameOver
}

func (gb *GameBoard) GetWinner() uint8 {
	if gb.player1Health <= 0 {
		return 1
	}
	return 0
}

// CreateGameBoard create a new GameBoard
func CreateGameBoard(width int16, height int16) *GameBoard {
	gb := GameBoard{
		width:                           width,
		height:                          height,
		ships:                           make([]ship, 0),
		fields:                          make([]field, 0),
		projectiles:                     make([]projectile, 0),
		totalPieceCounter:               6,
		player1Credits:                  500,
		player2Credits:                  500,
		player1Health:                   1000,
		player2Health:                   1000,
		creditsForCapitalShipDamageBool: false,
		creditsForShipDeathBool:         false,
		refundPercentageForShipDeath:    0.10,
		numCreditsForCapitalShipDamage:  1,
	}
	gb.setInitialShips()
	return &gb
}

func getNextProjectilePositionXY(projectile *projectile, logger *log.Logger) (int16, int16) {
	projectilePositionX := getMostRecentInt16CurveKeyFrameValue(projectile.positionCurve.x)
	projectilePositionY := getMostRecentInt16CurveKeyFrameValue(projectile.positionCurve.y)
	hypotenuse := projectile.maxSpeed

	var nextProjectilePositionX int16
	var nextProjectilePositionY int16
	angle := projectile.angleShot
	nextProjectilePositionXFloat := float64(projectilePositionX) + (hypotenuse * math.Cos(degreesToRadians(angle)))
	nextProjectilePositionYFloat := float64(projectilePositionY) + (hypotenuse * math.Sin(degreesToRadians(angle)))
	nextProjectilePositionX = int16(math.Round(nextProjectilePositionXFloat))
	nextProjectilePositionY = int16(math.Round(nextProjectilePositionYFloat))
	//logger.Printf("Float X: %v Int16 X: %v Float y Int16 %v Y: %v", nextProjectilePositionXFloat, nextProjectilePositionX, nextProjectilePositionYFloat, nextProjectilePositionY)

	return nextProjectilePositionX, nextProjectilePositionY
}

// RunFrame runs a single step of the game (http://www.vergenet.net/~conrad/boids/pseudocode.html)
func (gb *GameBoard) RunFrame(frame uint16) {
	gb.currentFrame = frame
	logger := log.New(os.Stdout, "", 0)
	if gb.player1Health <= 0 || gb.player2Health <= 0 {
		// game over
		return
	}
	if gb.player1Credits < 1000 {
		gb.player1Credits = gb.player1Credits + 1
	}
	if gb.player2Credits < 1000 {
		gb.player2Credits = gb.player2Credits + 1
	}

	// 1. All ships shoot and then move
	// 2. Move all projectiles
	// 3. Detect any collisions between ships and projectiles and resolve them

	// 1. All ships shoot and then move
	for i := range gb.ships {
		ship := &gb.ships[i]
		if (ship.existenceCurve[len(ship.existenceCurve)-1]).value == 1 {
			// Ship no longer exists, dont do anything
			continue
		}
		// Shoot
		shoot(gb, ship, frame, logger)
		v := velocityVector{0, 0}
		if ship.piece.name == interceptorName {
			v = gb.getInterceptorVector(ship, frame, logger)
		}
		if ship.piece.name == capitalShipName {
			v = gb.getCapitalShipVector(ship, frame, logger)
		}
		if ship.piece.name == destroyerName {
			v = gb.getDestroyerVector(ship, frame, logger)
		}
		if ship.piece.name == flackShipName {
			v = gb.getFlackShipVector(ship, frame, logger)
		}
		// Move ship
		// Get X
		ship.velocity.x = limitVelocityFloat64(v.x, ship.maxSpeed)
		shipPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
		newShipPositionX := shipPositionX + ship.velocity.x

		// Get Y
		ship.velocity.y = limitVelocityFloat64(v.y, ship.maxSpeed)
		shipPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
		newShipPositionY := shipPositionY + ship.velocity.y

		// Limit X & Y
		// TODO: WILL NEED TO EVENTUALLY BE DONE TO MAKE IT SO THAT THE BOIDS DONT TURN SO SHARPLEY, BUT MORE GRADUALLY INSTEAD

		// Set X & Y
		if shipPositionY != newShipPositionY {
			newYFrame := int16CurveKeyFrame{newShipPositionY, frame}
			ship.positionCurve.y = append(ship.positionCurve.y, newYFrame)
		}
		if shipPositionX != newShipPositionX {
			newXFrame := int16CurveKeyFrame{newShipPositionX, frame}
			ship.positionCurve.x = append(ship.positionCurve.x, newXFrame)
		}
		// Move Hitbox
		if ship.hasPolygonHitbox {
			ship.polygonHitbox.Pos = c2d.NewVector(float64(newShipPositionX), float64(newShipPositionY))
		}

		if (newShipPositionY == shipPositionY) && (newShipPositionX == shipPositionX) {
			// No need to find rotation, its in the exact same position
			continue
		}

		// Get & Set Rotation
		if ship.shipClass == "INTERCEPTOR" {
			// We only want interceptors to rotate
			currentShipPosition := point64{float64(shipPositionX), float64(shipPositionY)}
			nextShipPosition := point64{float64(newShipPositionX), float64(newShipPositionY)}
			rotation := getRotation360(currentShipPosition, nextShipPosition)
			previousRotation := getMostRecentInt16CurveKeyFrameValue(ship.rotationCurve)
			if previousRotation != rotation {
				newRotationFrame := int16CurveKeyFrame{rotation, frame}
				ship.rotationCurve = append(ship.rotationCurve, newRotationFrame)
				rotationRads := degreesToRadians(float64(rotation))
				// Rotate Hitbox
				if ship.hasPolygonHitbox {
					ship.polygonHitbox.SetAngle(rotationRads)
				}
			}
		}
	}

	// Move all projectiles
	for i := range gb.projectiles {
		projectile := &gb.projectiles[i]
		if (len(projectile.existenceCurve) % 2) == 0 {
			// Projectile is not alive, do not move it
			continue
		}
		if projectile.existenceCurve[len(projectile.existenceCurve)-1].time == frame {
			// Projectile was created this frame, do not move it
			continue
		}
		projectile.timeExisted = projectile.timeExisted + 1
		if projectile.timeExisted > projectile.existenceTime {
			// Projectile's life is now over, do not move it
			projectile.existenceCurve = append(projectile.existenceCurve, uint8CurveKeyFrame{1, frame})

			// Make sure its death position is known
			lastPositionX := projectile.positionCurve.x[len(projectile.positionCurve.x)-1]
			if lastPositionX.time != frame {
				newXFrame := int16CurveKeyFrame{lastPositionX.value, frame}
				projectile.positionCurve.x = append(projectile.positionCurve.x, newXFrame)
			}
			lastPositionY := projectile.positionCurve.y[len(projectile.positionCurve.y)-1]
			if lastPositionY.time != frame {
				newYFrame := int16CurveKeyFrame{lastPositionY.value, frame}
				projectile.positionCurve.y = append(projectile.positionCurve.y, newYFrame)
			}
			continue
		}
		// Projectile exists, move it
		newProjectilePositionX, newProjectilePositionY := getNextProjectilePositionXY(projectile, logger)
		projectilePositionX := getMostRecentInt16CurveKeyFrameValue(projectile.positionCurve.x)
		projectilePositionY := getMostRecentInt16CurveKeyFrameValue(projectile.positionCurve.y)
		// Set X & Y
		if projectilePositionX != newProjectilePositionX {
			newXFrame := int16CurveKeyFrame{newProjectilePositionX, frame}
			projectile.positionCurve.x = append(projectile.positionCurve.x, newXFrame)
		}
		if projectilePositionY != newProjectilePositionY {
			newYFrame := int16CurveKeyFrame{newProjectilePositionY, frame}
			projectile.positionCurve.y = append(projectile.positionCurve.y, newYFrame)
		}
		projectile.polygonHitbox.Pos = c2d.Vector{X: float64(newProjectilePositionX), Y: float64(newProjectilePositionY)}
		handleProjectileCollisions(gb, projectile, frame)
	}
}

func giveCreditsForCapitalShipDamage(gb *GameBoard, player uint8, damage int16) {
	creditsToAdd := damage * gb.numCreditsForCapitalShipDamage
	if player == 0 {
		newCredits := creditsToAdd + gb.player1Credits
		if newCredits > 1000 {
			newCredits = 1000
		}
		gb.player1Credits = newCredits
	} else {
		newCredits := creditsToAdd + gb.player2Credits
		if newCredits > 1000 {
			newCredits = 1000
		}
		gb.player2Credits = newCredits
	}
}

func handleCapitalShipDamage(gb *GameBoard, ship *ship, damage int16) {
	if ship.player == 0 {
		gb.player1Health = ship.health
	} else {
		gb.player2Health = ship.health
	}
	if gb.creditsForCapitalShipDamageBool {
		giveCreditsForCapitalShipDamage(gb, ship.player, damage)
	}
}

func handleShipDeath(gb *GameBoard, ship *ship) {
	creditsToRefund := int16(gb.refundPercentageForShipDeath * float64(shipBlueprints[ship.name].creditsCost))
	if gb.creditsForShipDeathBool {
		if ship.player == 0 {
			newCredits := gb.player1Credits + creditsToRefund
			if newCredits > 1000 {
				newCredits = 1000
			}
			gb.player1Credits = newCredits
		} else {
			newCredits := gb.player2Credits + creditsToRefund
			if newCredits > 1000 {
				newCredits = 1000
			}
			gb.player2Credits = newCredits
		}
	}
}

func handleProjectileCollisions(gb *GameBoard, projectile *projectile, frame uint16) {
	for i := range gb.ships {
		ship := &gb.ships[i]
		exists := (ship.existenceCurve[len(ship.existenceCurve)-1].value != uint8(1))
		if ship.player != projectile.player && exists {
			if collision, _ := c2d.TestPolygonPolygon(ship.polygonHitbox, projectile.polygonHitbox); collision {
				damage := calculateDamage(ship, projectile)
				ship.health = ship.health - damage
				if ship.name == capitalShipName {
					handleCapitalShipDamage(gb, ship, damage)
				}
				if ship.health <= 0 {
					ship.existenceCurve = append(ship.existenceCurve, uint8CurveKeyFrame{1, frame})
					handleShipDeath(gb, ship)
				}
				projectile.existenceCurve = append(projectile.existenceCurve, uint8CurveKeyFrame{1, frame})

				// Make sure its death position is known
				lastPositionX := projectile.positionCurve.x[len(projectile.positionCurve.x)-1]
				if lastPositionX.time != frame {
					newXFrame := int16CurveKeyFrame{lastPositionX.value, frame}
					projectile.positionCurve.x = append(projectile.positionCurve.x, newXFrame)
				}
				lastPositionY := projectile.positionCurve.y[len(projectile.positionCurve.y)-1]
				if lastPositionY.time != frame {
					newYFrame := int16CurveKeyFrame{lastPositionY.value, frame}
					projectile.positionCurve.y = append(projectile.positionCurve.y, newYFrame)
				}
				return
			}
		}
	}
}

func (gb *GameBoard) getNextPieceID() uint16 {
	gb.totalPieceCounter = gb.totalPieceCounter + 1
	return gb.totalPieceCounter
}

// findAvailablePieceID find, if possible, a piece of the same type that no longer exists
func (gb *GameBoard) findAvailablePieceID(pieceType string, pieceName pieceName, player uint8) (uint16, bool) {
	pieceID := uint16(0)
	found := false
	if pieceType == "projectile" {
		for _, projectile := range gb.projectiles {
			if (projectile.name == pieceName) && (projectile.player == player) {
				if (len(projectile.existenceCurve) % 2) == 0 {
					lastExistenceIndex := len(projectile.existenceCurve) - 1
					if projectile.existenceCurve[lastExistenceIndex].time < gb.currentFrame {
						return projectile.pieceID, true
					}
				}
			}
		}
	}
	if pieceType == "ship" {
		for _, ship := range gb.ships {
			if (ship.name == pieceName) && (ship.player == player) {
				if (len(ship.existenceCurve) % 2) == 0 {
					lastExistenceIndex := len(ship.existenceCurve) - 1
					if ship.existenceCurve[lastExistenceIndex].time < gb.currentFrame {
						return ship.pieceID, true
					}
				}
			}
		}
	}
	return pieceID, found

}

func (gb *GameBoard) getPieceID(pieceType string, pieceName pieceName, player uint8) (uint16, bool) {
	pieceID, ok := gb.findAvailablePieceID(pieceType, pieceName, player)
	if !ok {
		pieceID = gb.getNextPieceID()
	}
	return pieceID, ok
}

func calculateDamage(ship *ship, projectile *projectile) int16 {
	if ship.armorClass == "LIGHT" {
		return int16(float64(projectile.damage) * projectile.damageMultiplierLight)
	}
	if ship.armorClass == "MEDIUM" {
		return int16(float64(projectile.damage) * projectile.damageMultiplierMedium)
	}
	if ship.armorClass == "HEAVY" {
		return int16(float64(projectile.damage) * projectile.damageMultiplierHeavy)
	}
	return projectile.damage
}

func interceptorShoot(gb *GameBoard, ship *ship, frame uint16, angle float64) {
	// Get our position
	shipPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	shipPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)

	// Get the pieceID for the projectile
	pieceID, isReusablePieceID := gb.getPieceID("projectile", laserName, ship.player)
	if isReusablePieceID {
		// Projectile is already in projectiles, its values just need to be changed
		for i := range gb.projectiles {
			projectile := &gb.projectiles[i]
			if projectile.pieceID == pieceID {
				projectile.timeExisted = 0
				projectile.angleShot = angle
				newYFrame := int16CurveKeyFrame{shipPositionY, frame}
				newXFrame := int16CurveKeyFrame{shipPositionX, frame}
				rotationFrame := int16CurveKeyFrame{int16(angle), frame}
				existenceFrame := uint8CurveKeyFrame{0, frame}
				projectile.positionCurve.x = append(projectile.positionCurve.x, newXFrame)
				projectile.positionCurve.y = append(projectile.positionCurve.y, newYFrame)
				projectile.rotationCurve = append(projectile.rotationCurve, rotationFrame)
				projectile.existenceCurve = append(projectile.existenceCurve, existenceFrame)
			}
		}

	} else {
		// Create and add the new projectile
		projectileSpawnInfo := newProjectileInfo{
			position: position{shipPositionX, shipPositionY},
			name:     laserName,
			pieceID:  pieceID,
			time:     frame,
			player:   ship.player,
			angle:    angle,
		}
		projectile := newProjectile(projectileSpawnInfo)
		gb.projectiles = append(gb.projectiles, *projectile)
	}
}

func destroyerShoot(gb *GameBoard, ship *ship, frame uint16, angle float64) {
	// Get our position
	shipPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	shipPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)

	// Get the pieceID for the projectile
	pieceID, isReusablePieceID := gb.getPieceID("projectile", bigLaserName, ship.player)
	if isReusablePieceID {
		// Projectile is already in projectiles, its values just need to be changed
		for i := range gb.projectiles {
			projectile := &gb.projectiles[i]
			if projectile.pieceID == pieceID {
				projectile.timeExisted = 0
				projectile.angleShot = angle
				newYFrame := int16CurveKeyFrame{shipPositionY, frame}
				newXFrame := int16CurveKeyFrame{shipPositionX, frame}
				rotationFrame := int16CurveKeyFrame{int16(angle), frame}
				existenceFrame := uint8CurveKeyFrame{0, frame}
				projectile.positionCurve.x = append(projectile.positionCurve.x, newXFrame)
				projectile.positionCurve.y = append(projectile.positionCurve.y, newYFrame)
				projectile.rotationCurve = append(projectile.rotationCurve, rotationFrame)
				projectile.existenceCurve = append(projectile.existenceCurve, existenceFrame)
			}
		}

	} else {
		// Create and add the new projectile
		projectileSpawnInfo := newProjectileInfo{
			position: position{shipPositionX, shipPositionY},
			name:     bigLaserName,
			pieceID:  pieceID,
			time:     frame,
			player:   ship.player,
			angle:    angle,
		}
		projectile := newProjectile(projectileSpawnInfo)
		gb.projectiles = append(gb.projectiles, *projectile)
	}
}

func flackShipShoot(gb *GameBoard, ship *ship, frame uint16, angle float64) {
	// Get our position
	shipPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
	shipPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
	for i := 0; i < 3; i++ {
		angleAdjuster := int16(0)
		if i != 1 {
			if i < 1 {
				angleAdjuster = -1 * int16(i) * 10
			}
			if i > 1 {
				angleAdjuster = int16(i) * 10
			}
		}
		adjustedAngle := angle + float64(angleAdjuster)
		// Get the pieceID for the projectile
		pieceID, isReusablePieceID := gb.getPieceID("projectile", flackLaserName, ship.player)
		if isReusablePieceID {
			// Projectile is already in projectiles, its values just need to be changed
			for i := range gb.projectiles {
				projectile := &gb.projectiles[i]
				if projectile.pieceID == pieceID {
					projectile.timeExisted = 0
					projectile.angleShot = adjustedAngle
					newYFrame := int16CurveKeyFrame{shipPositionY, frame}
					newXFrame := int16CurveKeyFrame{shipPositionX, frame}
					rotationFrame := int16CurveKeyFrame{int16(adjustedAngle), frame}
					existenceFrame := uint8CurveKeyFrame{0, frame}
					projectile.positionCurve.x = append(projectile.positionCurve.x, newXFrame)
					projectile.positionCurve.y = append(projectile.positionCurve.y, newYFrame)
					projectile.rotationCurve = append(projectile.rotationCurve, rotationFrame)
					projectile.existenceCurve = append(projectile.existenceCurve, existenceFrame)
				}
			}

		} else {
			// Create and add the new projectile
			projectileSpawnInfo := newProjectileInfo{
				position: position{shipPositionX, shipPositionY},
				name:     flackLaserName,
				pieceID:  pieceID,
				time:     frame,
				player:   ship.player,
				angle:    adjustedAngle,
			}
			projectile := newProjectile(projectileSpawnInfo)
			gb.projectiles = append(gb.projectiles, *projectile)
		}
	}

}

func shoot(gb *GameBoard, ship *ship, frame uint16, logger *log.Logger) {
	if ship.shipClass == "CAPITAL_SHIP" {
		return
	}
	if ship.currentWeaponCoolDownFrame == 0 {
		// Find a ship to shoot
		nearestEnemyShip := findNearestEnemyShip(gb, ship)
		// Get our position
		shipPositionX := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.x)
		shipPositionY := getMostRecentInt16CurveKeyFrameValue(ship.positionCurve.y)
		shipPosition := point64{float64(shipPositionX), float64(shipPositionY)}
		// Get current enemy position
		enemyShipPositionX := getMostRecentInt16CurveKeyFrameValue(nearestEnemyShip.positionCurve.x)
		enemyShipPositionY := getMostRecentInt16CurveKeyFrameValue(nearestEnemyShip.positionCurve.y)
		enemyShipPosition := point64{float64(enemyShipPositionX), float64(enemyShipPositionY)}

		// Get future enemy position
		distanceBetweenShips := lengthOfLine(shipPosition, enemyShipPosition)
		numFramesToHit := uint16(distanceBetweenShips / projectileBlueprints[laserName].maxSpeed)
		maxLeadFrames := uint16(60)
		enemyShipFuturePositionX := getCurveInt16ValueAtTime(nearestEnemyShip.positionCurve.x, frame+numFramesToHit, nearestEnemyShip.existenceCurve)
		enemyShipFuturePositionY := getCurveInt16ValueAtTime(nearestEnemyShip.positionCurve.y, frame+numFramesToHit, nearestEnemyShip.existenceCurve)

		enemyShipFuturePosition := point64{float64(enemyShipFuturePositionX), float64(enemyShipFuturePositionY)}
		// Set shoot angle to use future enemy position
		var angle float64
		if numFramesToHit < maxLeadFrames {
			// We should be able to lead the shot
			angle = getRotation360Float(enemyShipFuturePosition, shipPosition)
		} else {
			// It is unlikely that we would be able to lead the shot
			angle = getRotation360Float(enemyShipPosition, shipPosition)
		}
		//logger.Printf("Nearest Enemy Ship X: %v Y: %v Enemy Future Pos X: %v Y: %v My Position X: %v Y: %v Shooting at Angle: %v", enemyShipPositionX, enemyShipPositionY, enemyShipFuturePositionX, enemyShipFuturePositionY, shipPositionX, shipPositionY, angle)
		//angle = -45
		//logger.Printf("Shooting at angle/heading: %v", angle)
		if ship.shipClass == "INTERCEPTOR" {
			interceptorShoot(gb, ship, frame, angle)
		}
		if ship.shipClass == "DESTROYER" {
			destroyerShoot(gb, ship, frame, angle)
		}
		if ship.shipClass == "FLACK_SHIP" {
			flackShipShoot(gb, ship, frame, angle)
		}

		// We just fired, put the ship weapon on cooldown
		ship.currentWeaponCoolDownFrame = ship.weaponCoolDownFrames
	} else {
		// Ship weapon is on cooldown
		ship.currentWeaponCoolDownFrame = ship.currentWeaponCoolDownFrame - 1
	}

}

func (gb *GameBoard) setInitialShips() {
	initialShips := gb.createInitialShips()
	gb.ships = initialShips
}

func (gb *GameBoard) getInterceptorVector(ship *ship, rame uint16, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	// Get Vectors from Rules
	var v1 velocityVector
	if gb.enemyShipsWithPieceNameExist(interceptorName, ship.player) {
		// Attack interceptors->destroyers->flackShips->capitalShip
		v1 = boidRule10(gb, ship, interceptorName, 100, logger)
	} else {
		if gb.enemyShipsWithPieceNameExist(destroyerName, ship.player) {
			v1 = boidRule10(gb, ship, destroyerName, 100, logger)
		} else {
			if gb.enemyShipsWithPieceNameExist(flackShipName, ship.player) {
				v1 = boidRule10(gb, ship, flackShipName, 100, logger)
			} else {
				v1 = boidRule10(gb, ship, capitalShipName, 100, logger)
			}
		}
	}
	v4 := boidRule4Ship(gb, ship, logger)
	v3 := boidRule9Ship(gb, ship, 100, logger)
	v5 := boidRule8Ship(gb, ship, 10, logger)

	// Sum Vectors
	v.x = ship.velocity.x + v1.x + v4.x + v3.x + v5.x
	v.y = ship.velocity.y + v1.y + v4.y + v3.y + v5.y
	return v
}

func (gb *GameBoard) enemyShipsWithPieceNameExist(name pieceName, player uint8) bool {
	exists := false
	for i := range gb.ships {
		if (gb.ships[i].name == name) && (player != gb.ships[i].player) && (gb.ships[i].existenceCurve[len(gb.ships[i].existenceCurve)-1].value != 1) {
			exists = true
			return exists
		}
	}
	return exists
}

func (gb *GameBoard) getDestroyerVector(ship *ship, rame uint16, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	// Get Vectors from Rules
	var v1 velocityVector
	if gb.enemyShipsWithPieceNameExist(destroyerName, ship.player) {
		// Attack destroyers->flackShips->capitalShip
		v1 = boidRule10(gb, ship, destroyerName, 100, logger)
	} else {
		if gb.enemyShipsWithPieceNameExist(flackShipName, ship.player) {
			v1 = boidRule10(gb, ship, flackShipName, 100, logger)
		} else {
			v1 = boidRule10(gb, ship, capitalShipName, 100, logger)
		}
	}
	v4 := boidRule4Ship(gb, ship, logger)
	v3 := boidRule9Ship(gb, ship, 100, logger)
	v5 := boidRule8Ship(gb, ship, 1, logger)

	// Sum Vectors
	v.x = ship.velocity.x + v1.x + v4.x + v3.x + v5.x
	v.y = ship.velocity.y + v1.y + v4.y + v3.y + v5.y
	return v
}

func (gb *GameBoard) getFlackShipVector(ship *ship, rame uint16, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	// Get Vectors from Rules
	var v1 velocityVector
	if gb.enemyShipsWithPieceNameExist(interceptorName, ship.player) {
		// Attack interceptors->flackShips->capitalShip
		v1 = boidRule10(gb, ship, interceptorName, 100, logger)
	} else {
		if gb.enemyShipsWithPieceNameExist(flackShipName, ship.player) {
			v1 = boidRule10(gb, ship, flackShipName, 100, logger)
		} else {
			v1 = boidRule10(gb, ship, capitalShipName, 100, logger)
		}
	}
	v4 := boidRule4Ship(gb, ship, logger)
	v3 := boidRule9Ship(gb, ship, 100, logger)
	v5 := boidRule8Ship(gb, ship, 1, logger)

	// Sum Vectors
	v.x = ship.velocity.x + v1.x + v4.x + v3.x + v5.x
	v.y = ship.velocity.y + v1.y + v4.y + v3.y + v5.y
	return v
}

func (gb *GameBoard) getCapitalShipVector(ship *ship, frame uint16, logger *log.Logger) velocityVector {
	v := velocityVector{0, 0}
	return v
}

func findNearestEnemyShip(gb *GameBoard, thisShip *ship) *ship {
	var nearestEnemyShip ship
	var distanceToNearestEnemyShip float64
	numEnemyShips := 0
	myX := getMostRecentInt16CurveKeyFrameValue(thisShip.positionCurve.x)
	myY := getMostRecentInt16CurveKeyFrameValue(thisShip.positionCurve.y)
	myPos := point64{float64(myX), float64(myY)}
	for _, spShip := range gb.ships {
		if (spShip.existenceCurve[len(spShip.existenceCurve)-1]).value == 1 {
			// Ship does not exist, ignore it
			continue
		}
		if spShip.player != thisShip.player {
			enemyX := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.x)
			enemyY := getMostRecentInt16CurveKeyFrameValue(spShip.positionCurve.y)
			posOfEnemyShip := point64{float64(enemyX), float64(enemyY)}
			distanceToEnemyShip := lengthOfLine(myPos, posOfEnemyShip)
			if numEnemyShips == 0 {
				distanceToNearestEnemyShip = distanceToEnemyShip
				nearestEnemyShip = spShip
			} else {
				if distanceToEnemyShip < distanceToNearestEnemyShip {
					distanceToNearestEnemyShip = distanceToEnemyShip
					nearestEnemyShip = spShip
				}
			}
			numEnemyShips++
		}
	}
	return &nearestEnemyShip
}

func limitVector(currentPosition point64, nextPosition point64, previousPosition point64, maxRotationSpeed int16, logger *log.Logger) (int16, int16) {
	newShipPositionX := int16(nextPosition.x)
	newShipPositionY := int16(nextPosition.y)
	// previous rotation
	rotation1 := calculateRotationDegrees(currentPosition, nextPosition)
	// next rotation
	rotation2 := calculateRotationDegrees(previousPosition, currentPosition)
	rotationDifference := rotation2 - rotation1
	absRotationDifference := getAbsInt16(rotationDifference)
	logger.Printf("Previous rotation: %v Next Rotation: %v Abs Diff: %v", rotation1, rotation2, absRotationDifference)
	if absRotationDifference > maxRotationSpeed {
		// length of new vector
		length := int16(math.Sqrt(math.Pow(float64(nextPosition.x-currentPosition.x), 2) + math.Pow(float64(nextPosition.y-currentPosition.y), 2)))
		if rotation2 < 0 {
			negMaxRotationSpeed := (maxRotationSpeed * -1) + rotation1
			opposite := calculateOppositeUsingSin(negMaxRotationSpeed, length)
			adjacent := calculateAdjacentUsingCos(negMaxRotationSpeed, length)
			newShipPositionX = adjacent
			newShipPositionY = opposite
			logger.Printf("Actual Next Rotation: %v", negMaxRotationSpeed)
		} else {

			maxRotationSpeed += rotation1
			opposite := calculateOppositeUsingSin(maxRotationSpeed, length)
			adjacent := calculateAdjacentUsingCos(maxRotationSpeed, length)
			newShipPositionX = adjacent
			newShipPositionY = opposite
			logger.Printf("Actual Next Rotation: %v", maxRotationSpeed)
		}
	}

	return newShipPositionX, newShipPositionY
}

func limitVelocity(v int16, vlim int16) int16 {
	limitedVelocity := v
	absV := getAbsInt16(v)
	if absV > vlim {
		limitedVelocity = (v / absV) * vlim
	}
	return limitedVelocity
}

func limitVelocityFloat64(v int16, vlim float64) int16 {
	var float64limitedVelocity float64
	float64limitedVelocity = float64(v)
	absV := getAbsInt16(v)
	if float64(absV) > vlim {
		float64limitedVelocity = (float64(v) / float64(absV)) * vlim
	}
	return int16(float64limitedVelocity)
}

func (gb *GameBoard) vertexIsOnBoard(vertex point, boardPosition position) bool {
	var isOnBoard = true
	x := vertex.x + int16(boardPosition.x)
	y := vertex.y + int16(boardPosition.y)
	boardWidth := int16(gb.width)
	boardHeight := int16(gb.height)
	if (x < 0) || x > boardWidth-1 {
		isOnBoard = false
	}
	if (y < 0) || y > boardHeight-1 {
		isOnBoard = false
	}
	return isOnBoard
}

func newProjectile(newProjectileInfo newProjectileInfo) *projectile {
	blueprint := projectileBlueprints[newProjectileInfo.name]
	rotationCurve := createInt16Curve(int16(newProjectileInfo.angle), newProjectileInfo.time)
	positionCurveX := createInt16Curve(newProjectileInfo.position.x, newProjectileInfo.time)
	positionCurveY := createInt16Curve(newProjectileInfo.position.y, newProjectileInfo.time)
	existenceCurve := createUint8Curve(0, newProjectileInfo.time)
	positionCurve := positionCurve{positionCurveX, positionCurveY}

	return &projectile{
		piece: piece{
			identifier:       blueprint.identifier,
			positionCurve:    positionCurve,
			rotationCurve:    rotationCurve,
			player:           newProjectileInfo.player,
			name:             newProjectileInfo.name,
			pieceID:          newProjectileInfo.pieceID,
			existenceCurve:   existenceCurve,
			circleHitbox:     blueprint.circleHitbox,
			polygonHitbox:    blueprint.polygonHitbox,
			hasPolygonHitbox: blueprint.hasPolygonHitbox,
		},
		damage:                 blueprint.damage,
		maxSpeed:               blueprint.maxSpeed,
		angleShot:              newProjectileInfo.angle,
		existenceTime:          blueprint.existenceTime,
		timeExisted:            0,
		damageMultiplierLight:  blueprint.damageMultiplierLight,
		damageMultiplierMedium: blueprint.damageMultiplierMedium,
		damageMultiplierHeavy:  blueprint.damageMultiplierHeavy,
	}
}

func newShip(newShipInfo newShipInfo) *ship {
	blueprint := shipBlueprints[newShipInfo.name]
	rotationCurve := createInt16Curve(newShipInfo.rotation, newShipInfo.time)
	positionCurveX := createInt16Curve(newShipInfo.position.x, newShipInfo.time)
	positionCurveY := createInt16Curve(newShipInfo.position.y, newShipInfo.time)
	existenceCurve := createUint8Curve(0, newShipInfo.time)
	positionCurve := positionCurve{positionCurveX, positionCurveY}
	return &ship{
		piece: piece{
			identifier:       blueprint.identifier,
			positionCurve:    positionCurve,
			rotationCurve:    rotationCurve,
			player:           newShipInfo.player,
			name:             newShipInfo.name,
			pieceID:          newShipInfo.pieceID,
			existenceCurve:   existenceCurve,
			circleHitbox:     blueprint.circleHitbox,
			polygonHitbox:    blueprint.polygonHitbox,
			hasPolygonHitbox: blueprint.hasPolygonHitbox,
		},
		velocity:                   velocityVector{0, 0},
		maxSpeed:                   blueprint.maxSpeed,
		maxRotationSpeed:           blueprint.maxRotationSpeed,
		shipClass:                  blueprint.shipClass,
		armorClass:                 blueprint.armorClass,
		personalSpace:              blueprint.personalSpace,
		mass:                       blueprint.mass,
		weaponCoolDownFrames:       blueprint.weaponCoolDownFrames,
		currentWeaponCoolDownFrame: blueprint.currentWeaponCoolDownFrame,
		health:                     blueprint.health,
	}
}

func pieceExists(piece piece) bool {
	return len(piece.existenceCurve) == 1
}

func (gb *GameBoard) handlePlayerNameMessage(splitMessage []string) {
	if splitMessage[1] == "0" {
		gb.player1Name = splitMessage[2]
	} else {
		gb.player2Name = splitMessage[2]
	}
}

func (gb *GameBoard) HaveBothPlayerNames() bool {
	//fmt.Print("Inside HaveBothPlayerNames")
	//fmt.Printf("Player 1 Name: %v, Player 2 Name: %v", gb.player1Name, gb.player2Name)
	return (gb.player1Name != "") && (gb.player2Name != "")
}

func (gb *GameBoard) HandleMessageFromPlayer(message string) {
	fmt.Print("\nRECIEVED MESSAGE FROM PLAYER:")
	fmt.Print(message)
	splitMessage := strings.Split(message, " ")
	if splitMessage[0] == "PLAYERNAME" {
		gb.handlePlayerNameMessage(splitMessage)
	} else {
		shipToSpawnString := splitMessage[0]
		shipToSpawn := stringPieceNameMap[shipToSpawnString]
		playerString := splitMessage[1]
		playerInt, err := strconv.Atoi(playerString)
		if err != nil {
			fmt.Println(err)
			playerInt = 0
		}
		// Determine if the player has enough credits to spawn the ship
		shipCost := shipBlueprints[shipToSpawn].creditsCost
		if playerInt == 0 {
			// Player1
			if gb.player1Credits >= shipCost {
				gb.player1Credits = gb.player1Credits - shipCost
			} else {
				return
			}
		} else {
			//Player 2
			if gb.player2Credits >= shipCost {
				gb.player2Credits = gb.player2Credits - shipCost
			} else {
				return
			}
		}
		gb.spawnShip(uint8(playerInt), shipToSpawn)
	}
}

func (gb *GameBoard) spawnShip(player uint8, pieceName pieceName) {
	if pieceName == destroyerName {
		pieceID1 := gb.getNextPieceID()
		pieceID2 := gb.getNextPieceID()
		player1DestroyerPositions := []position{{100, 750}, {300, 750}}
		player2DestroyerPositions := []position{{100, 50}, {300, 50}}
		var positions []position
		var rotation = int16(90)
		if player == 0 {
			positions = player1DestroyerPositions
		}
		if player == 1 {
			positions = player2DestroyerPositions
			rotation = 270
		}
		playerDestroyerInfo1 := newShipInfo{
			position: positions[0],
			name:     destroyerName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID1,
			time:     gb.currentFrame,
		}
		playerDestroyerInfo2 := newShipInfo{
			position: positions[1],
			name:     destroyerName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID2,
			time:     gb.currentFrame,
		}
		ship1 := newShip(playerDestroyerInfo1)
		ship2 := newShip(playerDestroyerInfo2)
		gb.ships = append(gb.ships, *ship1, *ship2)
	}
	if pieceName == flackShipName {
		pieceID1 := gb.getNextPieceID()
		pieceID2 := gb.getNextPieceID()
		player1DestroyerPositions := []position{{100, 750}, {300, 750}}
		player2DestroyerPositions := []position{{100, 50}, {300, 50}}
		var positions []position
		var rotation = int16(90)
		if player == 0 {
			positions = player1DestroyerPositions
		}
		if player == 1 {
			positions = player2DestroyerPositions
			rotation = 270
		}
		playerDestroyerInfo1 := newShipInfo{
			position: positions[0],
			name:     flackShipName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID1,
			time:     gb.currentFrame,
		}
		playerDestroyerInfo2 := newShipInfo{
			position: positions[1],
			name:     flackShipName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID2,
			time:     gb.currentFrame,
		}
		ship1 := newShip(playerDestroyerInfo1)
		ship2 := newShip(playerDestroyerInfo2)
		gb.ships = append(gb.ships, *ship1, *ship2)
	}
	if pieceName == interceptorName {
		pieceID1 := gb.getNextPieceID()
		pieceID2 := gb.getNextPieceID()
		pieceID3 := gb.getNextPieceID()
		pieceID4 := gb.getNextPieceID()
		player1InterceptorPositions := []position{{100, 750}, {150, 750}, {250, 750}, {300, 750}}
		player2InterceptorPositions := []position{{100, 50}, {150, 50}, {250, 50}, {300, 50}}
		var positions []position
		var rotation = int16(90)
		if player == 0 {
			positions = player1InterceptorPositions
		}
		if player == 1 {
			positions = player2InterceptorPositions
			rotation = 270
		}
		player1InterceptorInfo1 := newShipInfo{
			position: positions[0],
			name:     interceptorName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID1,
			time:     gb.currentFrame,
		}
		player1InterceptorInfo2 := newShipInfo{
			position: positions[1],
			name:     interceptorName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID2,
			time:     gb.currentFrame,
		}
		player1InterceptorInfo3 := newShipInfo{
			position: positions[2],
			name:     interceptorName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID3,
			time:     gb.currentFrame,
		}
		player1InterceptorInfo4 := newShipInfo{
			position: positions[3],
			name:     interceptorName,
			rotation: rotation,
			player:   player,
			pieceID:  pieceID4,
			time:     gb.currentFrame,
		}
		ship1 := newShip(player1InterceptorInfo1)
		ship2 := newShip(player1InterceptorInfo2)
		ship3 := newShip(player1InterceptorInfo3)
		ship4 := newShip(player1InterceptorInfo4)
		debugging := false
		if !debugging {
			gb.ships = append(gb.ships, *ship1, *ship2, *ship3, *ship4)
		} else {
			gb.ships = append(gb.ships, *ship1)
		}
	}
}

func (gb *GameBoard) createInitialShips() []ship {
	initialShips := make([]ship, 2)
	player1CapitalShipInfo := newShipInfo{
		position: position{224, 750},
		name:     capitalShipName,
		rotation: 90,
		player:   0,
		pieceID:  0,
		time:     gb.currentFrame,
	}
	player2CapitalShipInfo := newShipInfo{
		position: position{224, 50},
		name:     capitalShipName,
		rotation: 270,
		player:   1,
		pieceID:  3,
		time:     gb.currentFrame,
	}

	initialShips[0] = *newShip(player1CapitalShipInfo)
	initialShips[1] = *newShip(player2CapitalShipInfo)

	return initialShips
}
