import React from 'react';
import "./GameScreen.css";
import PlayAgainBox from "../playagainbox/PlayAgainBox";
import UnitPurchaseButton from "../buttons/UnitPurchaseButton";
import CountdownTimer from "../countdown/CountdownTimer";

class GameScreen extends React.Component{
    constructor(props){
        super(props);
        this.state = {
            currentFrame: 0,
            framerateMS:16,
            gameLoopID:null,
            restartGameLoopID:null,
            countdownContainerAnimation: "instant-fade-in-long-fade-out",
            versusContainerAnimation:"fade-in-and-out",
            numberAnimation:"number-fade-in-and-out",
            timedCountdownContainerAnimation: "fade-in-and-out-longer",
            savedCanvasWidth:0,
            savedCanvasHeight:0,
            previousOpponentName:null,
            previousOpponentHealth:null,
            previousOpponentCredits:null,
            myPreviousHealth:null,
            myPreviousCredits:null,
            gameCurrentlyOver:false,
        };
        this.defaultCountdownContainerAnimation = "instant-fade-in-long-fade-out";
        this.defaultCountdownContainerAnimationTwo = "not-quite-instant-fade-in-long-fade-out";
        this.defaultVersusContainerAnimation = "fade-in-and-out";
        this.defaultNumberAnimation = "number-fade-in-and-out";
        this.defaultTimedCountdownContainerAnimation = "fade-in-and-out-longer";
    }

    onClickPlayAgain = () => {
        this.props.onClickPlayAgainButton();
    }

    onClickHome = () => {
        this.props.disconnect();
        this.props.setViewingGameScreen(false);
    }

    drawBoard = (board, time) => {
        const canvas = document.getElementById("mainCanvas");
        var ctx = canvas.getContext('2d');
        ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
        this.drawPieces(ctx, board.pieces, board.images, board.player,board.width, board.height, time);
    }

    drawPieces = (ctx, pieces, images, player, boardWidth, boardHeight, time) =>{
        for (let i=0;i<pieces.length;i++){
            const piece = pieces[i];
            let exists = this.determineIfPieceExists(piece, time, false);
            if (exists){
                let positionCurveX = piece.positionCurveX;
                let positionCurveY = piece.positionCurveY;
                const pieceImage = this.getPieceImage(piece, images, player);
                let xPositionAtTime = this.getCurveValueAtTime(positionCurveX, time, piece);
                let yPositionAtTime = this.getCurveValueAtTime(positionCurveY, time, piece);
                if ((xPositionAtTime === null)||(yPositionAtTime === null)){
                    // could not interpolate or client predict, dont draw this
                    continue;
                }
                //let xPositionAtTime = this.getCurveValueAtTimeNoInterpolationNoClientSidePrediction(positionCurveX, time);
                //let yPositionAtTime = this.getCurveValueAtTimeNoInterpolationNoClientSidePrediction(positionCurveY, time);
                let rotationAtTimeDegrees = this.getCurveValueAtTimeNoInterpolationNoClientSidePrediction(piece, piece.rotationCurve, time);
               
                // we rotate by 180 cause ship images are all facing right
                let actualRotationAtTimeDegrees = rotationAtTimeDegrees+180;
                if (player===1){
                    // Player 2 needs the board flipped
                    actualRotationAtTimeDegrees = rotationAtTimeDegrees;
                    xPositionAtTime = boardWidth - xPositionAtTime;
                    yPositionAtTime = boardHeight - yPositionAtTime;
                }
                const rotationAtTime = this.degreesToRadians(actualRotationAtTimeDegrees);
                ctx.save();
                ctx.translate(Math.round(xPositionAtTime/this.props.cnvs.widthRatio), Math.round(yPositionAtTime/this.props.cnvs.heightRatio));
                ctx.rotate(rotationAtTime);
                ctx.drawImage(pieceImage, -25, -25);
                ctx.restore();
            }
        }
    }

    degreesToRadians = (degrees) =>{
        return degrees * (Math.PI/180);
    }

    getCurveValueAtTimeNoInterpolationNoClientSidePrediction = (piece, curve, time) =>{
        let nearestTimeDifference = 999999;
        let nearestCurveValue = null;
        for (let i=0;i<curve.length;i++){
            if (Math.abs(time-curve[i].time)<nearestTimeDifference){
                nearestTimeDifference = Math.abs(time-curve[i].time)
                nearestCurveValue = curve[i].value
            }
        }
        return nearestCurveValue;
    }

    newFindNearestStartCurveKeyFrame = (curve, time) => {
        let nearestTimeDifference = 999999;
        let nearestCurveKeyFrame = null;
        for (let i=0;i<curve.length;i++){
            if (((time-curve[i].time)<nearestTimeDifference)&&((time-curve[i].time)>=0)){
                nearestTimeDifference = time-curve[i].time
                nearestCurveKeyFrame = curve[i]
            }
        }
        return nearestCurveKeyFrame;
    }

    newFindNearestEndCurveKeyFrame = (curve, time) => {
        let nearestTimeDifference = 999999;
        let nearestCurveKeyFrame = null;
        for (let i=curve.length-1;i>=0;i--){
            if ((Math.abs(time-curve[i].time)<nearestTimeDifference)&&((time-curve[i].time)<=0)){
                nearestTimeDifference = Math.abs(time-curve[i].time)
                nearestCurveKeyFrame = curve[i]
            }
        }
        return nearestCurveKeyFrame;
    }
    
    getCurveValueAtTime = (curve, time, piece) =>{
        let value = null;

        // check simple case where we only have one value
        if (curve.length===1){
            return curve[0].value
        }

        // check simple case where we have exactly this value
        for (let i=0;i<curve.length;i++){
            if(curve[i].time===time){
                return curve[i].value
            }
        }

        // check if we can interpolate the value
        if (curve[curve.length-1].time>time){
            // we have a point further in the future than where we are currently, interpolation is possible
            let startCurveKeyFrame = this.newFindNearestStartCurveKeyFrame(curve, time)
            let endCurveKeyFrame = this.newFindNearestEndCurveKeyFrame(curve, time);
            
            // interpolate
            let slope = (endCurveKeyFrame.value - startCurveKeyFrame.value)/(endCurveKeyFrame.time - startCurveKeyFrame.time);
            let b = startCurveKeyFrame.value - (slope*startCurveKeyFrame.time);
            let interpolatedValue = (slope*time)+b;
            return interpolatedValue;
        }
        
        // check if we can client side predict the value
       if (curve[curve.length-1].time<time){
            // we only have points further in the past than the requested time
            let startCurveKeyFrame = curve[curve.length-2];
            let endCurveKeyFrame = curve[curve.length-1];
            let lastExistenceCurveKeyFrame = piece.existenceCurve[piece.existenceCurve.length-1];
            if ((lastExistenceCurveKeyFrame.value===0)&&(lastExistenceCurveKeyFrame.time<=startCurveKeyFrame.time)){
                // last available existence frame says this piece existed
                let slope = (endCurveKeyFrame.value - startCurveKeyFrame.value)/(endCurveKeyFrame.time - startCurveKeyFrame.time);
                let b = startCurveKeyFrame.value - (slope*startCurveKeyFrame.time);
                let predictedValue = (slope*time)+b;
                return predictedValue;
            }
        }
        return value;
    }

    determineIfPieceExists = (piece, time, debug) =>{
        let exists = false;
        for (let i=piece.existenceCurve.length-1;i>=0;i--){
            let pieceOfCurve = piece.existenceCurve[i];
            if (pieceOfCurve.time<=time){
                if (pieceOfCurve.value === 0){
                    exists = true
                }
                if (debug){
                    console.log(`Piece exists: ${exists} @ time: ${time}`);
                    console.log(piece.existenceCurve);
                }
                return exists
            }
        }
        if (debug){
            console.log(`Piece has no existence curve at or before time: ${time}`);
        }
        return exists;
    }

    getPieceImage = (piece, images, player) =>{
        for (let i=0;i<images.length;i++){
            if (images[i].Identifier === piece.identifier){
                if (player === 1){
                    // Player 2, has to use the opposite image
                    if (images[i].Player !== piece.player){
                        return images[i].image;
                    }
                    /* start debug */
                    if (piece.identifier === 255){
                        //console.log(piece.existenceCurve)
                    }
                    /* end debug*/
                }
                else{
                    // Player 1 uses all images as normal
                    if (images[i].Player === piece.player){
                        return images[i].image;
                    }
                }
            }
        }
    }
    // restart animation: really dumb, but DOM wont update so we need a timeout so that react doesnt update on same batch
    //https://github.com/facebook/react/issues/7142
    restartCountdownTimerAnimations = () => {
        this.setState({
            countdownContainerAnimation:"", 
            versusContainerAnimation:"", 
            numberAnimation:"", 
            timedCountdownContainerAnimation:""}, () => {
            setTimeout(() => this.setState({
                countdownContainerAnimation:this.defaultCountdownContainerAnimationTwo, 
                versusContainerAnimation:this.defaultVersusContainerAnimation,
                numberAnimation:this.defaultNumberAnimation,
                timedCountdownContainerAnimation:this.defaultTimedCountdownContainerAnimation}), 50)
        });

    }

    mainGameLoop = () => {
            if (this.props.board.gameIsOver){
                let myPreviousHealth = 0;
                let myPreviousCredits = 0;
                let previousOpponentCredits = 0;
                let previousOpponentHealth = 0;
                let previousOpponentName = this.props.board.enemyPlayerName;
                if (this.props.board.player === 0){
                    // Player is Player 1
                    myPreviousHealth = Math.ceil(this.props.board.player1Health/10);
                    myPreviousCredits= this.props.board.player1Credits;
                    previousOpponentHealth = Math.ceil(this.props.board.player2Health/10);
                    previousOpponentCredits = this.props.board.player2Credits;
                }
                else{
                    // Player is Player 2
                    myPreviousHealth = Math.ceil(this.props.board.player2Health/10);
                    myPreviousCredits = this.props.board.player2Credits;
                    previousOpponentHealth = Math.ceil(this.props.board.player1Health/10);
                    previousOpponentCredits = this.props.board.player1Credits;
                }

                const {framerateMS} = this.state;
                var restartGameLoopID = setInterval(this.restartGameLoop, framerateMS);
                this.setState({
                    gameCurrentlyOver:true,
                    restartGameLoopID:restartGameLoopID, 
                    previousOpponentCredits:previousOpponentCredits, 
                    previousOpponentName:previousOpponentName, 
                    previousOpponentHealth:previousOpponentHealth, 
                    myPreviousHealth:myPreviousHealth, 
                    myPreviousCredits:myPreviousCredits});
                this.shutdownGameLoop();
                return;
            }
            this.handleLatency();
            const {currentFrame} = this.state;
            this.drawBoard(this.props.board, currentFrame);
            const newCurrentFrame = currentFrame + 1;
            this.setState({currentFrame:newCurrentFrame});
    }

    // handleLatency(): attempts to make the game appear smooth regardless of the latency by adjusting what the current frame
    // the game is displaying if there is too much drift
    handleLatency = () => {
        const {currentFrame} = this.state;
        // 6.25 frames occur in 100ms
        const maxFramesAllowedAhead = 12;
        const maxFramesAllowedBehind = 12;
        // Are we past the biggest time in a key frame by some arbitrary amount?
        // Go to the biggest key frame
        const largestFrame = this.props.board.getLargestTimeInAllCurves();
        if (largestFrame === 0){
            return;
        }
        const maxFramesAhead = largestFrame + maxFramesAllowedAhead;
        if (maxFramesAhead < currentFrame){
            // we are past what is acceptable, we have to adjust the current frame
            let newCurrentFrame = largestFrame - maxFramesAllowedBehind; 
            if (newCurrentFrame < 0){
                newCurrentFrame = 0;
            }
            this.setState({currentFrame:newCurrentFrame});
        }
        const maxFramesBehind = largestFrame - maxFramesAllowedBehind;
        if (currentFrame < maxFramesBehind){
            // we are behind what is acceptable, we have to adjust the current frame
            let newCurrentFrame = largestFrame - 12;
            if (newCurrentFrame < 0){
                newCurrentFrame = 0;
            }
            this.setState({currentFrame:newCurrentFrame})
        }

    }

    restartGameLoop = () =>{
        if (!this.props.board.gameIsOver && (this.props.board.enemyPlayerName !== null)){
            const {framerateMS, restartGameLoopID} = this.state;
            try{
                clearInterval(restartGameLoopID);
                this.setState({restartGameLoopID:null});
                this.updateDimensions();
                this.restartCountdownTimerAnimations();
                setTimeout(()=>{
                    var gameLoopID = setInterval(this.mainGameLoop, framerateMS);
                    this.setState({currentFrame: 0, gameLoopID:gameLoopID, gameCurrentlyOver:false,});
                    
                },6500);
            }
            catch(err){
                console.log(err);
            }
        }
    }

    shutdownGameLoop = () => {
        const {gameLoopID} = this.state;
        if (gameLoopID != null){
            try{
                clearInterval(gameLoopID)
                this.setState({gameLoopID:null})
            }
            catch(err){
                console.log(err);
            }
        }
    }

    componentWillUnmount(){
        this.shutdownGameLoop();
        this.props.hideFooter(false);
        window.removeEventListener('keydown', this.handleKeyDown);
        window.removeEventListener("resize", this.updateDimensions);
    }

    handleKeyDown = (event) => {
        const {gameCurrentlyOver} = this.state;
        if (!gameCurrentlyOver){
            if (event.code==="KeyQ"){
                this.handleFlackClick();
            }
            if (event.code==="KeyW"){
                this.handleDestroyerClick();
            }
            if (event.code==="KeyE"){
                this.handleInterceptorClick();
            }
        }
        if (gameCurrentlyOver && !(this.props.lookingForMatch)){
            if (event.code==="KeyR"){
                this.onClickPlayAgain();
            }
        }
    }

    componentDidMount(){
        //console.log("GameScreen mounted!")
        this.props.hideFooter(true);
        window.addEventListener('keydown', this.handleKeyDown);
        const canvas = document.getElementById("mainCanvas");
        var ctx = canvas.getContext('2d');
 
        // Get the new width
        let newWidth = 500;
        if (newWidth >= window.innerWidth-100){
            newWidth=window.innerWidth-100;
        }

        if (window.innerWidth<482){
            newWidth = window.innerWidth;
        }

        // Get the new height
        let newHeight = window.innerHeight*0.8;
        if (newHeight > 650){
            newHeight = 650
        }

        if (newHeight+90+53+50+50+20> window.innerHeight){
            newHeight = window.innerHeight - 90 - 53 - 50 - 50 - 20;
        }
        // because border is 2px
        newWidth = newWidth - 2;
        this.props.cnvs.setDimensions(newWidth, newHeight);
        ctx.canvas.width  = this.props.cnvs.width;
        ctx.canvas.height = this.props.cnvs.height;

        window.addEventListener("resize", this.updateDimensions);

        const {framerateMS} = this.state;
        var gameLoopID = setInterval(this.mainGameLoop, framerateMS)
        this.setState({gameLoopID:gameLoopID})
        this.updateDimensions();
    }

    updateDimensions = () => {
        const canvas = document.getElementById("mainCanvas")
        var ctx = canvas.getContext('2d');

        // Get the new width
        let newWidth = 500;
        if (newWidth >= window.innerWidth-100){
            newWidth=window.innerWidth-100;
        }

        if (window.innerWidth<482){
            newWidth = window.innerWidth;
        }

        // Get the new height
        let newHeight = window.innerHeight*0.8;
        if (newHeight > 650){
            newHeight = 650
        }

        if (newHeight+104+53+50+50+20> window.innerHeight){
            newHeight = window.innerHeight - 104 - 53 - 50 - 50 - 20;
        }
        newWidth = newWidth - 2;
        //cnvs.setDimensions(((window.innerHeight*0.8)/(16/9)), window.innerHeight*0.8);
        this.props.cnvs.setDimensions(newWidth, newHeight);
        ctx.canvas.width  = this.props.cnvs.width;
        ctx.canvas.height = this.props.cnvs.height;
        this.props.cnvs.heightRatio = this.props.board.height / this.props.cnvs.height;
        this.props.cnvs.widthRatio = this.props.board.width / this.props.cnvs.width;
    }

    handleInterceptorClick = () => {
        this.props.send("INTERCEPTOR "+this.props.board.player);
    }

    handleDestroyerClick = () => {
        this.props.send("DESTROYER "+this.props.board.player);
    }

    handleFlackClick = () => {
        this.props.send("FLACK_SHIP "+this.props.board.player);
    }

    render() {
        const {savedCanvasWidth, savedCanvasHeight} = this.state;
        let mainCanvasWidth = savedCanvasWidth;
        let mainCanvasHeight = savedCanvasHeight;
        if ((this.props.cnvs.width !== 0) && (this.props.cnvs.height !== 0)){
            if (this.props.cnvs.width != mainCanvasWidth){
                mainCanvasWidth = this.props.cnvs.width;
                this.setState({savedCanvasWidth:this.props.cnvs.width});
            }
            if (this.props.cnvs.height != mainCanvasHeight){
                mainCanvasHeight = this.props.cnvs.height;
                this.setState({savedCanvasHeight:this.props.cnvs.height});
            }
        }
        const latency = this.props.getLatency();
        const { board } = this.props;
        const { nickname } = this.props;
        let playerHealth, playerCredits, opponentHealth, opponentName, opponentCredits;
        if (this.props.board.player === 0){
            // Player is Player 1
            playerHealth = Math.ceil(this.props.board.player1Health/10);
            playerCredits= this.props.board.player1Credits;
            opponentHealth = Math.ceil(this.props.board.player2Health/10);
            opponentName = this.props.board.enemyPlayerName;
            opponentCredits = this.props.board.player2Credits;
        }
        else{
            // Player is Player 2
            playerHealth = Math.ceil(this.props.board.player2Health/10);
            playerCredits= this.props.board.player2Credits;
            opponentHealth = Math.ceil(this.props.board.player1Health/10);
            opponentName = this.props.board.enemyPlayerName;
            opponentCredits = this.props.board.player1Credits;
        }

        // Basically, show old board until we actually have a new opponent
        if (this.props.board.enemyPlayerName === null){
            const {previousOpponentName} = this.state;
            if (previousOpponentName !== null){
                const {previousOpponentHealth, previousOpponentCredits, myPreviousCredits, myPreviousHealth} = this.state;
                playerHealth = myPreviousHealth;
                playerCredits = myPreviousCredits;
                opponentCredits = previousOpponentCredits;
                opponentHealth = previousOpponentHealth;
                opponentName = previousOpponentName;
            }
        }
        return(
            <div className="game-screen-wrapper">
                <div className="enemy-status-bar status-bar">
                    <div className="outer-status-bar-column">
                        <p className="player-name">{opponentName}</p>
                    </div>
                    <div className="outer-status-bar-column">
                        <div className="resources-column">
                        <div className="outer-status-bar-column">
                                <div className="image-and-text-flex">
                                    <img className="health-image" src="/heart_inverted_dark.png"></img>
                                    <p>{opponentHealth}</p>
                                </div>
                            </div>
                        <div className="outer-status-bar-column">
                                <div className="image-and-text-flex">
                                    <img className="health-image"  src="/yen_dark.png"></img>
                                    <p>{opponentCredits}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="main-canvas-container" style={{width:mainCanvasWidth, height:mainCanvasHeight}}>
                    <CountdownTimer playerName={nickname} opponentName={opponentName} containerWidth={mainCanvasWidth} containerHeight={this.props.cnvs.height} countdownContainerAnimation={this.state.countdownContainerAnimation} versusContainerAnimation={this.state.versusContainerAnimation} numberAnimation={this.state.numberAnimation} timedCountdownContainerAnimation={this.state.timedCountdownContainerAnimation}/>
                    <canvas id="mainCanvas" className="main-canvas"/>
                    {(board.gameIsOver || (board.enemyPlayerName === null)) && (
                    <PlayAgainBox errorInfo={this.props.errorInfo} searchInfo = {this.props.searchInfo} lookingForMatch={this.props.lookingForMatch} onClickPlayAgain={this.onClickPlayAgain} onClickHome={this.onClickHome} victoryMessage={this.props.board.victoryMessage}/>
                )}
                </div>
                <div className="player-status-bar status-bar" style={{width:mainCanvasWidth}}>
                    <div className="outer-status-bar-column">
                        <p className="player-name">{nickname}</p>
                    </div>
                    <div className="outer-status-bar-column">
                        <div className="resources-column">
                        <div className="outer-status-bar-column">
                                <div className="image-and-text-flex">
                                    <img className="health-image" src="/heart_inverted.png"></img>
                                    <p>{playerHealth}</p>
                                </div>
                            </div>
                            <div className="outer-status-bar-column">
                                <div className="image-and-text-flex">
                                    <img className="health-image"  src="/yen.png"></img>
                                    <p>{playerCredits}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="controls-div">
                        <UnitPurchaseButton unitName={"rock"} playerCredits={playerCredits} price={200} shortcut={"Q"} imgAlt="Flack Image" unitImgSrc="/Flack_Ship.png" yenImgSrc="/yen_original.png" onClick={this.handleFlackClick} />
                        <UnitPurchaseButton unitName={"paper"} playerCredits={playerCredits} price={200} shortcut={"W"} imgAlt="Destroyer Image" unitImgSrc="/Destroyer_v2.png" yenImgSrc="/yen_original.png" onClick={this.handleDestroyerClick} />
                        <UnitPurchaseButton unitName={"scissors"} playerCredits={playerCredits} price={50} shortcut={"E"} imgAlt="Interceptor Image" unitImgSrc="/Interceptor.png" yenImgSrc="/yen_original.png" onClick={this.handleInterceptorClick} />
                </div>
            </div>
        );
    }
}

export default GameScreen;