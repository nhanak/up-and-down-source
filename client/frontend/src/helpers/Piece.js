
class Piece{
    constructor(pieceID, player, identifier, existenceCurve, positionCurveX, positionCurveY, rotationCurve){
        this.pieceID = pieceID;
        this.player = player;
        this.identifier = identifier;
        this.existenceCurve = existenceCurve;
        this.positionCurveX = positionCurveX;
        this.positionCurveY = positionCurveY;
        this.rotationCurve = rotationCurve;
    }
}

export default Piece