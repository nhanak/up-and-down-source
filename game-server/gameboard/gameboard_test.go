package gameboard

import (
	"testing"
)

func Test_pieceIsOnBoard(t *testing.T) {

	capitalShipTestInfo := newShipInfo{
		position: position{0, 0},
		name:     capitalShipName,
		player:   0,
		rotation: 0,
	}
	capitalShipTestInfo2 := newShipInfo{
		position: position{30, 30},
		name:     capitalShipName,
		player:   0,
		rotation: 0,
	}
	capitalShipTestInfo3 := newShipInfo{
		position: position{20, 10},
		name:     capitalShipName,
		player:   0,
		rotation: 0,
	}
	capitalShipTestInfo4 := newShipInfo{
		position: position{19, 10},
		name:     capitalShipName,
		player:   0,
		rotation: 0,
	}
	tables := []struct {
		boardWidth  uint16
		boardHeight uint16
		piece       *piece
		onBoard     bool
	}{
		{5, 5, &newShip(capitalShipTestInfo).piece, false},
		{100, 100, &newShip(capitalShipTestInfo2).piece, true},
		{100, 100, &newShip(capitalShipTestInfo3).piece, true},
		{100, 100, &newShip(capitalShipTestInfo4).piece, false},
	}
	for i, table := range tables {
		gb := gameBoard{width: table.boardWidth, height: table.boardHeight}
		pieceIsOnBoard := gb.isOnBoard(table.piece)
		if pieceIsOnBoard != table.onBoard {
			t.Errorf("piece.isOnBoard returned incorrect value for piece #%v, got: %v, want: %v on board dim(%v, %v)", i, pieceIsOnBoard, table.onBoard, table.boardWidth, table.boardHeight)
		}
	}
}

func Test_vertexIsOnBoard(t *testing.T) {
	tables := []struct {
		width          uint16
		height         uint16
		boardPositionX uint16
		boardPositionY uint16
		vertexX        int16
		vertexY        int16
		onBoard        bool
	}{
		{10, 10, 0, 0, -5, 0, false},
		{10, 10, 0, 0, 0, 0, true},
		{10, 10, 0, 0, 0, 9, true},
		{10, 10, 0, 0, 9, 9, true},
		{10, 10, 0, 5, 0, 0, true},
		{10, 10, 5, 0, 0, 0, true},
		{10, 10, 5, 0, 5, 0, false},
		{10, 10, 0, 5, 0, 5, false},
		{10, 10, 0, 5, 0, 4, true},
		{10, 10, 5, 0, 4, 0, true},
		{10, 10, 5, 0, -5, 0, true},
		{10, 10, 0, 5, 0, -5, true},
	}
	for _, table := range tables {
		gb := gameBoard{width: table.width, height: table.height}
		vertex := point{x: table.vertexX, y: table.vertexY}
		bp := position{x: table.boardPositionX, y: table.boardPositionY}
		vertexIsOnBoard := gb.vertexIsOnBoard(vertex, bp)
		if vertexIsOnBoard != table.onBoard {
			t.Errorf("vertexIsOnBoard returned incorrect value, got: %v, want: %v on board dim(%v, %v) for vertex(%v,%v) @ boardPosition(%v,%v)", vertexIsOnBoard, table.onBoard, table.width, table.height, table.vertexX, table.vertexY, table.boardPositionX, table.boardPositionY)
		}
	}
}
