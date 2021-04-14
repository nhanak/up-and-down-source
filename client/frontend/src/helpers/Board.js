class Board{
    constructor(width, height){
        this.width = width;
        this.height = height;
        this.pieces = [];
        this.images = [];
        this.player1Credits = 0;
        this.player1Health = 1000;
        this.player2Credits = 0;
        this.player2Health = 1000;
        this.gameOver = false;
        this.victor = null;
        this.player= 0;
        this.enemyPlayerName = null;
        this.gameIsOver = false;
        this.victoryMessage = "";
    }

    setDimensions(width, height){
        this.width = width;
        this.height = height;
    }

    addPiece(piece){
        this.pieces.push(piece);
    }

    getPieceWithPieceID(pieceID){
        let piece = null;
        for(let i=0;i<this.pieces.length;i++){
            if (this.pieces[i].pieceID===pieceID){
                return this.pieces[i];
            }
        }
        return piece;
    }

    clearPieces(){
        this.pieces = [];
    }

    setPlayer(player){
        this.player = player
    }

    setEnemyPlayerName(name){
        this.enemyPlayerName = name;
    }

    addPositionXCurveKeyFrame(pieceID, keyframe){
        for(let i=0;i<this.pieces.length;i++){
            if (this.pieces[i].pieceID === pieceID){
                this.pieces[i].positionCurveX.push(keyframe);
            }
        }
    }

    addPositionYCurveKeyFrame(pieceID, keyframe){
        for(let i=0;i<this.pieces.length;i++){
            if (this.pieces[i].pieceID === pieceID){
                this.pieces[i].positionCurveY.push(keyframe);
            }
        }
    }

    addRotationCurveKeyFrame(pieceID, keyframe){
        for(let i=0;i<this.pieces.length;i++){
            if (this.pieces[i].pieceID === pieceID){
                this.pieces[i].rotationCurve.push(keyframe);
            }
        }
    }

    addExistenceCurveKeyFrame(pieceID, keyframe){
        for(let i=0;i<this.pieces.length;i++){
            if (this.pieces[i].pieceID === pieceID){
                this.pieces[i].existenceCurve.push(keyframe);
            }
        }
    }

    getLargestTimeInAllCurves(){
        let largestTime = 0;
        let largestTimesArray = [];
        for(let i=0;i<this.pieces.length;i++){
            let largestExistenceCurveTime = this.getLargestTimeInCurve(this.pieces[i].existenceCurve);
            let largestRotationCurveTime = this.getLargestTimeInCurve(this.pieces[i].rotationCurve);
            let largestPositionCurveXTime = this.getLargestTimeInCurve(this.pieces[i].positionCurveX);
            let largestPositionCurveYTime = this.getLargestTimeInCurve(this.pieces[i].positionCurveY);
            largestTimesArray.push(largestExistenceCurveTime);
            largestTimesArray.push(largestRotationCurveTime);
            largestTimesArray.push(largestPositionCurveXTime);
            largestTimesArray.push(largestPositionCurveYTime);
        }
        for (let i=0;i<largestTimesArray.length;i++){
            if (largestTimesArray[i]>largestTime){
                largestTime = largestTimesArray[i];
            }
        }
        return largestTime;
    }

    getLargestTimeInCurve(curve){
        return curve[curve.length-1].time
    }
}

export default Board