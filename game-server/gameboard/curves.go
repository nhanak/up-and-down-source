package gameboard

type positionCurve struct {
	x []int16CurveKeyFrame
	y []int16CurveKeyFrame
}

type int16CurveKeyFrame struct {
	value int16
	time  uint16
}

type uint8CurveKeyFrame struct {
	value uint8
	time  uint16
}

func createInt16Curve(val int16, time uint16) []int16CurveKeyFrame {
	keyframe := int16CurveKeyFrame{value: val, time: time}
	curve := make([]int16CurveKeyFrame, 1)
	curve[0] = keyframe
	return curve
}

func createUint8Curve(val uint8, time uint16) []uint8CurveKeyFrame {
	keyframe := uint8CurveKeyFrame{value: val, time: time}
	curve := make([]uint8CurveKeyFrame, 1)
	curve[0] = keyframe
	return curve
}

func getMostRecentInt16CurveKeyFrameValue(curve []int16CurveKeyFrame) int16 {
	return curve[len(curve)-1].value
}

func getMostRecentUint16CurveKeyFrameTime(curve []int16CurveKeyFrame) uint16 {
	return curve[len(curve)-1].time
}

func getMostRecentUint16CurveKeyFrame(curve []int16CurveKeyFrame) int16CurveKeyFrame {
	return curve[len(curve)-1]
}

func getMostRecentUint8CurveKeyFrame(curve []uint8CurveKeyFrame) uint8CurveKeyFrame {
	return curve[len(curve)-1]
}

func isNewerUint16CurveKeyFrameAvailable(curve []int16CurveKeyFrame, t uint16) bool {
	return (curve[len(curve)-1]).time > t
}

func isNewerUint8CurveKeyFrameAvailable(curve []uint8CurveKeyFrame, t uint16) bool {
	return (curve[len(curve)-1]).time > t
}

func isInt16CurveLongerThanLength(curve []int16CurveKeyFrame, length int) bool {
	return len(curve) > length
}

func isUint8CurveLongerThanLength(curve []uint8CurveKeyFrame, length int) bool {
	return len(curve) > length
}

func findNearestInt16StartCurveKeyFrame(curve []int16CurveKeyFrame, time uint16) int16CurveKeyFrame {
	smallestTimeDifference := time - curve[0].time
	smallestIndex := 0
	for i := 1; i < len(curve)-1; i++ {
		if curve[i].time < time {
			// make sure that this curve occurs before
			if (time - curve[i].time) < smallestTimeDifference {
				smallestTimeDifference = time - curve[i].time
				smallestIndex = i
			} else {
				return curve[smallestIndex]
			}
		} else {
			return curve[smallestIndex]
		}
	}
	return curve[smallestIndex]
}

func findNearestInt16EndCurveKeyFrame(curve []int16CurveKeyFrame, time uint16) int16CurveKeyFrame {
	smallestTimeDifference := curve[len(curve)-1].time - time
	smallestIndex := len(curve) - 1
	for i := len(curve) - 2; i > 0; i-- {
		// start from the end
		if curve[i].time > time {
			// make sure that this curve occurs after
			if (curve[i].time - time) < smallestTimeDifference {
				smallestTimeDifference = curve[i].time - time
				smallestIndex = i
			} else {
				return curve[smallestIndex]
			}
		} else {
			return curve[smallestIndex]
		}
	}
	return curve[smallestIndex]
}

func findNearestUint8StartCurveKeyFrame(curve []uint8CurveKeyFrame, time uint16) uint8CurveKeyFrame {
	smallestTimeDifference := time - curve[0].time
	smallestIndex := 0
	for i := 1; i < len(curve)-1; i++ {
		if curve[i].time < time {
			// make sure that this curve occurs before
			if (time - curve[i].time) < smallestTimeDifference {
				smallestTimeDifference = time - curve[i].time
				smallestIndex = i
			} else {
				return curve[smallestIndex]
			}
		} else {
			return curve[smallestIndex]
		}
	}
	return curve[smallestIndex]
}

func findNearestUint8EndCurveKeyFrame(curve []uint8CurveKeyFrame, time uint16) uint8CurveKeyFrame {
	smallestTimeDifference := curve[len(curve)-1].time - time
	smallestIndex := len(curve) - 1
	for i := len(curve) - 2; i > 0; i-- {
		// start from the end
		if curve[i].time > time {
			// make sure that this curve occurs after
			if (curve[i].time - time) < smallestTimeDifference {
				smallestTimeDifference = curve[i].time - time
				smallestIndex = i
			} else {
				return curve[smallestIndex]
			}
		} else {
			return curve[smallestIndex]
		}
	}
	return curve[smallestIndex]
}

func determineIfInt16CurveValueExists(existenceCurve []uint8CurveKeyFrame, frame uint16) bool {
	exists := false
	startCurveKeyFrame := findNearestUint8StartCurveKeyFrame(existenceCurve, frame)
	endCurveKeyFrame := findNearestUint8EndCurveKeyFrame(existenceCurve, frame)
	if startCurveKeyFrame.value == 0 {
		// start key frame says piece existed
		if endCurveKeyFrame.value == 1 {
			// end frame is a non existance frame
			if endCurveKeyFrame.time < frame {
				// the current time is before the end frame
				exists = true
			}
		}
	}
	if endCurveKeyFrame.value == 0 {
		// end curve key frame says piece existed
		if endCurveKeyFrame.time <= frame {
			// we are equal or past that time
			exists = true
		}
	}
	return exists
}

func getCurveInt16ValueAtTime(curve []int16CurveKeyFrame, time uint16, existenceCurve []uint8CurveKeyFrame) int16 {
	// check simple case where we only have one value
	if len(curve) == 1 {
		//console.log("Returning only value for curve");
		return curve[0].value
	}
	// check simple case where we have exactly this value
	for i := range curve {
		if curve[i].time == time {
			return curve[i].value
		}
	}
	// check if we can interpolate the value
	if len(curve) > 1 {
		// we have at least two points
		if curve[len(curve)-1].time > time {
			// we have a point further in the future than where we are currently, interpolation is possible

			// find closest startCurveKeyFrame and endCurveKeyFrame
			startCurveKeyFrame := findNearestInt16StartCurveKeyFrame(curve, time)
			endCurveKeyFrame := findNearestInt16EndCurveKeyFrame(curve, time)
			if !determineIfInt16CurveValueExists(existenceCurve, startCurveKeyFrame.time) {
				// piece did not exist at start time, we can't interpolate
				return endCurveKeyFrame.value
			}
			if (endCurveKeyFrame.time - startCurveKeyFrame.time) > 6 {
				// if we have obviously been reusing the same value for startCurveKeyFrame
				startCurveKeyFrame.time = endCurveKeyFrame.time - 6
			}

			slope := (endCurveKeyFrame.value - startCurveKeyFrame.value) / (int16(endCurveKeyFrame.time) - int16(startCurveKeyFrame.time))
			b := startCurveKeyFrame.value - (slope * int16(startCurveKeyFrame.time))
			interpolatedValue := (slope * int16(time)) + b
			return interpolatedValue
		}
	}
	// prediction of value
	if curve[len(curve)-1].time < time {
		// we only have points further in the past than the requested time
		if len(curve) > 1 {
			startCurveKeyFrame := curve[len(curve)-2]
			endCurveKeyFrame := curve[len(curve)-1]
			if !determineIfInt16CurveValueExists(existenceCurve, startCurveKeyFrame.time) {
				// piece did not exist at start time, we can't predict
				return endCurveKeyFrame.value
			}
			if (endCurveKeyFrame.time - startCurveKeyFrame.time) > 6 {
				// if we have obviously been reusing the same value for startCurveKeyFrame
				startCurveKeyFrame.time = endCurveKeyFrame.time - 6
			}
			slope := (endCurveKeyFrame.value - startCurveKeyFrame.value) / (int16(endCurveKeyFrame.time) - int16(startCurveKeyFrame.time))
			b := startCurveKeyFrame.value - (slope * int16(startCurveKeyFrame.time))
			predictedValue := (slope * int16(time)) + b
			return predictedValue
		}
	}
	return curve[0].value
}
