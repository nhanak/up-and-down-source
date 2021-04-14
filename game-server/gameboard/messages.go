package gameboard

import "encoding/binary"

// ImageURLData url+piece+player data of a image
type ImageURLData struct {
	Player     uint16
	Identifier uint8
	URL        string
}

// pieceIDSentTracker: tracks whether a piece has been sent over the network or not before
var pieceIDSentTrackerPlayerOne = make(map[uint16]bool, 0)
var pieceIDSentTrackerPlayerTwo = make(map[uint16]bool, 0)

func (gb *GameBoard) ResetPieceIDSentTracker() {
	pieceIDSentTrackerPlayerOne = make(map[uint16]bool, 0)
	pieceIDSentTrackerPlayerTwo = make(map[uint16]bool, 0)
}

// GetImageURLData gets the image urls for all pieces
func (gb *GameBoard) GetImageURLData() []ImageURLData {
	images := make([]ImageURLData, 0)
	images = append(images, ImageURLData{Player: 0, Identifier: identifierMap[capitalShipName], URL: "/public/images/ships/capital_ship_player_one_v2.png"})
	images = append(images, ImageURLData{Player: 1, Identifier: identifierMap[capitalShipName], URL: "/public/images/ships/capital_ship_player_two_v2.png"})
	images = append(images, ImageURLData{Player: 0, Identifier: identifierMap[interceptorName], URL: "/public/images/ships/interceptor_player_one_v2.png"})
	images = append(images, ImageURLData{Player: 1, Identifier: identifierMap[interceptorName], URL: "/public/images/ships/interceptor_player_two_v2.png"})
	images = append(images, ImageURLData{Player: 0, Identifier: identifierMap[destroyerName], URL: "/public/images/ships/destroyer_player_one_v3.png"})
	images = append(images, ImageURLData{Player: 1, Identifier: identifierMap[destroyerName], URL: "/public/images/ships/destroyer_player_two_v3.png"})
	images = append(images, ImageURLData{Player: 0, Identifier: identifierMap[flackShipName], URL: "/public/images/ships/flack_ship_player_one_v2.png"})
	images = append(images, ImageURLData{Player: 1, Identifier: identifierMap[flackShipName], URL: "/public/images/ships/flack_ship_player_two_v2.png"})
	images = append(images, ImageURLData{Player: 0, Identifier: identifierMap[laserName], URL: "/public/images/projectiles/laser_player_one.png"})
	images = append(images, ImageURLData{Player: 1, Identifier: identifierMap[laserName], URL: "/public/images/projectiles/laser_player_two.png"})
	images = append(images, ImageURLData{Player: 0, Identifier: identifierMap[bigLaserName], URL: "/public/images/projectiles/big_laser_player_one.png"})
	images = append(images, ImageURLData{Player: 1, Identifier: identifierMap[bigLaserName], URL: "/public/images/projectiles/big_laser_player_two.png"})
	images = append(images, ImageURLData{Player: 0, Identifier: identifierMap[flackLaserName], URL: "/public/images/projectiles/flack_laser_player_one.png"})
	images = append(images, ImageURLData{Player: 1, Identifier: identifierMap[flackLaserName], URL: "/public/images/projectiles/flack_laser_player_two.png"})
	return images
}

func (gb *GameBoard) GetGameOverMessage(winner uint8) []byte {
	// msg: first byte is message type (2=GameOverMessage)
	msg := make([]byte, 2)
	msg[0] = 2
	msg[1] = winner
	return msg
}

func (gb *GameBoard) GetPlayersHealthAndCreditsMessage() []byte {
	// msg: first byte is message type (1=players health and resources)
	msg := make([]byte, 9)
	msg[0] = 1
	binary.BigEndian.PutUint16(msg[1:], uint16(gb.player1Health))
	binary.BigEndian.PutUint16(msg[3:], uint16(gb.player1Credits))
	binary.BigEndian.PutUint16(msg[5:], uint16(gb.player2Health))
	binary.BigEndian.PutUint16(msg[7:], uint16(gb.player2Credits))
	return msg
}

// GetPiecesMessage gets the binary message of all pieces on board to be sent by the server to the client
func (gb *GameBoard) GetPiecesMessage(lastSentGameFrame uint16, player uint8) []byte {
	// msg: first byte is message type (0=piecesCurves)
	// next two bytes are how many bytes total are in the message
	msg := make([]byte, 3)
	for _, ship := range gb.ships {
		pieceMsg, emptyMessage := createPieceMessage(&ship.piece, lastSentGameFrame, player)
		if emptyMessage {
			continue
		}
		msg = append(msg, pieceMsg...)
	}
	for _, projectile := range gb.projectiles {
		pieceMsg, emptyMessage := createPieceMessage(&projectile.piece, lastSentGameFrame, player)
		if emptyMessage {
			continue
		}
		msg = append(msg, pieceMsg...)
	}
	messageLength := uint16(len(msg))
	binary.BigEndian.PutUint16(msg[1:], messageLength)
	return msg
}
func getPieceIDSentTracker(player uint8) map[uint16]bool {
	if player == 0 {
		return pieceIDSentTrackerPlayerOne
	}
	return pieceIDSentTrackerPlayerTwo
}
func createPieceMessage(piece *piece, lastSentGameFrame uint16, player uint8) ([]byte, bool) {
	emptyMsg := false
	// 255 = message end
	messageEnd := make([]byte, 1)
	messageEnd[0] = 255
	pieceID := piece.pieceID
	pieceMsg := make([]byte, 2)
	binary.BigEndian.PutUint16(pieceMsg[0:], pieceID) // first 2 bytes of pieceMsg are the pieceId
	pieceIDSentTracker := getPieceIDSentTracker(player)
	if _, ok := pieceIDSentTracker[piece.pieceID]; !ok {
		// piece has never been sent, so send all its initial values
		initialPieceMsg := getInitialPieceMsg(*piece)
		pieceMsg = append(pieceMsg, initialPieceMsg...)
		pieceIDSentTracker[pieceID] = true

		// now send the latest values if the piece has any
		latestInitialPieceMessage := getLatestInitialPieceMsg(*piece)
		pieceMsg = append(pieceMsg, latestInitialPieceMessage...)

		// add the end of the message
		pieceMsg = append(pieceMsg, messageEnd...)

	} else {
		// piece has been sent before, send latest of what hasnt been sent
		followUpPieceMsg := getFollowUpPieceMsg(*piece, lastSentGameFrame)
		if len(followUpPieceMsg) == 0 {
			emptyMsg = true
		}
		pieceMsg = append(pieceMsg, followUpPieceMsg...)
		pieceMsg = append(pieceMsg, messageEnd...)
	}
	return pieceMsg, emptyMsg
}

func getInitialPieceMsg(piece piece) []byte {
	msg := make([]byte, 18)
	msg[0] = 0                                                                   // message type, 0 = initialPieceMsg
	msg[1] = piece.player                                                        //player
	msg[2] = piece.identifier                                                    //identifier
	msg[3] = piece.existenceCurve[0].value                                       //first existence value
	binary.BigEndian.PutUint16(msg[4:], piece.existenceCurve[0].time)            //first existence time
	binary.BigEndian.PutUint16(msg[6:], uint16(piece.positionCurve.x[0].value))  //first x pos value
	binary.BigEndian.PutUint16(msg[8:], piece.positionCurve.x[0].time)           //first x pos time
	binary.BigEndian.PutUint16(msg[10:], uint16(piece.positionCurve.y[0].value)) //first y pos value
	binary.BigEndian.PutUint16(msg[12:], piece.positionCurve.y[0].time)          //first y pos time
	binary.BigEndian.PutUint16(msg[14:], uint16(piece.rotationCurve[0].value))   //first rotation value
	binary.BigEndian.PutUint16(msg[16:], piece.rotationCurve[0].time)            //first rotation pos time
	return msg
}

func getFollowUpPieceMsg(piece piece, lastSentGameFrame uint16) []byte {
	// message header
	//0 = initial pieces
	// 1 = final piece info
	// 2 = all positional (x,y,rotation)
	// 3 = x and y
	// 4 = x
	// 5 = y
	// 6 = rotation
	// 7 = existence

	msg := make([]byte, 0)

	// check if there were existence changes
	existenceChange := false
	for i := range piece.existenceCurve {
		if piece.existenceCurve[i].time > lastSentGameFrame {
			existenceChange = true
		}
	}
	if !existenceChange {
		// there was no change in existence, just send the latest
		newPositionXAvailable := isNewerUint16CurveKeyFrameAvailable(piece.positionCurve.x, lastSentGameFrame)
		if newPositionXAvailable {
			mostRecentKeyFrame := getMostRecentUint16CurveKeyFrame(piece.positionCurve.x)
			msg = append(msg, createInt16CurveKeyFrameMessage(4, mostRecentKeyFrame)...)
		}
		newPositionYAvailable := isNewerUint16CurveKeyFrameAvailable(piece.positionCurve.y, lastSentGameFrame)
		if newPositionYAvailable {
			mostRecentKeyFrame := getMostRecentUint16CurveKeyFrame(piece.positionCurve.y)
			msg = append(msg, createInt16CurveKeyFrameMessage(5, mostRecentKeyFrame)...)
		}
		newRotationAvailable := isNewerUint16CurveKeyFrameAvailable(piece.rotationCurve, lastSentGameFrame)
		if newRotationAvailable {
			mostRecentKeyFrame := getMostRecentUint16CurveKeyFrame(piece.rotationCurve)
			msg = append(msg, createInt16CurveKeyFrameMessage(6, mostRecentKeyFrame)...)
		}
	} else {
		// there was a change in existence
		// all existence frames after lastSentGameFrame have to be sent, not just the latest
		for i := range piece.existenceCurve {
			if piece.existenceCurve[i].time > lastSentGameFrame {
				existenceKeyFrame := piece.existenceCurve[i]
				msg = append(msg, createUint8CurveKeyFrameMessage(7, existenceKeyFrame)...)
				// send the value of x at time of existence key frame
				for j := range piece.positionCurve.x {
					if piece.positionCurve.x[j].time == existenceKeyFrame.time {
						positionXKeyFrame := piece.positionCurve.x[j]
						msg = append(msg, createInt16CurveKeyFrameMessage(4, positionXKeyFrame)...)
					}
				}
				// send the value of y at time of existence key frame
				for k := range piece.positionCurve.y {
					if piece.positionCurve.y[k].time == existenceKeyFrame.time {
						positionYKeyFrame := piece.positionCurve.y[k]
						msg = append(msg, createInt16CurveKeyFrameMessage(5, positionYKeyFrame)...)
					}
				}
				// send the value of rotation at time of existence key frame
				for l := range piece.rotationCurve {
					if piece.rotationCurve[l].time == existenceKeyFrame.time {
						rotationKeyFrame := piece.rotationCurve[l]
						msg = append(msg, createInt16CurveKeyFrameMessage(6, rotationKeyFrame)...)
					}
				}
				// if this is the last existence key frame, and the key frame value was that it existed
				// we have to send the last position value
				if existenceKeyFrame.value == 0 {
					lastExistenceKeyFrame := piece.existenceCurve[len(piece.existenceCurve)-1]
					if lastExistenceKeyFrame == existenceKeyFrame {
						lastPositionXKeyFrame := piece.positionCurve.x[len(piece.positionCurve.x)-1]
						lastPositionYKeyFrame := piece.positionCurve.y[len(piece.positionCurve.y)-1]
						lastRotationKeyFrame := piece.rotationCurve[len(piece.rotationCurve)-1]
						if existenceKeyFrame.time < lastPositionXKeyFrame.time {
							msg = append(msg, createInt16CurveKeyFrameMessage(4, lastPositionXKeyFrame)...)
						}
						if existenceKeyFrame.time < lastPositionYKeyFrame.time {
							msg = append(msg, createInt16CurveKeyFrameMessage(5, lastPositionYKeyFrame)...)
						}
						if existenceKeyFrame.time < lastRotationKeyFrame.time {
							msg = append(msg, createInt16CurveKeyFrameMessage(6, lastRotationKeyFrame)...)
						}
					}
				}
			}
		}

	}
	return msg
}

func getLatestInitialPieceMsg(piece piece) []byte {
	msg := make([]byte, 0)
	if isInt16CurveLongerThanLength(piece.positionCurve.x, 1) {
		mostRecentKeyFrame := getMostRecentUint16CurveKeyFrame(piece.positionCurve.x)
		msg = append(msg, createInt16CurveKeyFrameMessage(4, mostRecentKeyFrame)...)
	}
	if isInt16CurveLongerThanLength(piece.positionCurve.y, 1) {
		mostRecentKeyFrame := getMostRecentUint16CurveKeyFrame(piece.positionCurve.y)
		msg = append(msg, createInt16CurveKeyFrameMessage(5, mostRecentKeyFrame)...)
	}
	if isInt16CurveLongerThanLength(piece.rotationCurve, 1) {
		mostRecentKeyFrame := getMostRecentUint16CurveKeyFrame(piece.rotationCurve)
		msg = append(msg, createInt16CurveKeyFrameMessage(6, mostRecentKeyFrame)...)
	}
	if isUint8CurveLongerThanLength(piece.existenceCurve, 1) {
		mostRecentKeyFrame := getMostRecentUint8CurveKeyFrame(piece.existenceCurve)
		msg = append(msg, createUint8CurveKeyFrameMessage(7, mostRecentKeyFrame)...)
	}
	return msg
}

func createInt16CurveKeyFrameMessage(messageHeader uint8, keyframe int16CurveKeyFrame) []byte {
	msg := make([]byte, 5)
	msg[0] = messageHeader
	binary.BigEndian.PutUint16(msg[1:], uint16(keyframe.value))
	binary.BigEndian.PutUint16(msg[3:], keyframe.time)
	return msg
}

func createUint8CurveKeyFrameMessage(messageHeader uint8, keyframe uint8CurveKeyFrame) []byte {
	msg := make([]byte, 4)
	msg[0] = messageHeader
	msg[1] = keyframe.value
	binary.BigEndian.PutUint16(msg[2:], keyframe.time)
	return msg
}
