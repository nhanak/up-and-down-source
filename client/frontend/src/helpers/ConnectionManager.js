import Board from "./Board"
import Cnvs from "./Cnvs"
import Piece from "./Piece"
import CurveKeyFrame from "./CurveKeyFrame"

class ConnectionManager {
    constructor(nickname, setViewingGameScreen, handleGameServerConnection, handleGameServerDisconnection) {
        this.nickname = nickname;
        this.setViewingGameScreen = setViewingGameScreen;
        this.handleGameServerConnection = handleGameServerConnection;
        this.handleGameServerDisconnection = handleGameServerDisconnection;
        this.conn = null;
        this.board = new Board(0,0);
        this.cnvs = new Cnvs(0,0);
        this.recieved_first_message = false;
        this.lastMessageRecievedTimestamp = null;
        this.secondToLastMessageRecievedTimestamp = null;
        this.ip = null;
        this.port = null;
    }

    getLatency = () => {
        let latency = null;
        if ((this.lastMessageRecievedTimestamp!==null)&&(this.secondToLastMessageRecievedTimestamp!==null)){
            latency = this.lastMessageRecievedTimestamp - this.secondToLastMessageRecievedTimestamp;
        }
        return latency;
    }

    bigResetAndGameOverTrue = () => {
        let prevVictoryMessage = this.board.victoryMessage;
        this.handleGameServerDisconnection();
        this.resetConnectionManager();
        this.board.gameIsOver = true;
        this.board.victoryMessage = prevVictoryMessage;
    }

    send = (val) => {
        //console.log("Sending message: "+ val);
        if (!this.board.gameIsOver && (this.conn !== null)){
            this.conn.send(val)
        }
        if (this.conn === null){
            //console.info("Send connection was null :(");
            this.bigResetAndGameOverTrue();
        }
    }

    getBoard = () => {
        return this.board;
    }

    getCnvs = () => {
        return this.cnvs;
    }

    updateTimeStamps = () => {
        var d = new Date();
        var newLastMessageRecievedTimestamp = d.getTime();
        this.secondToLastMessageRecievedTimestamp = this.lastMessageRecievedTimestamp
        this.lastMessageRecievedTimestamp = newLastMessageRecievedTimestamp
    }

    resetConnectionManager = () => {
        this.conn = null;
        this.board = new Board(0,0);
        this.cnvs = new Cnvs(0,0);
        this.recieved_first_message = false;
        this.lastMessageRecievedTimestamp = null;
        this.secondToLastMessageRecievedTimestamp = null;
        this.ip = null;
        this.port = null;
    }

    disconnect = () => {
        this.resetConnectionManager();
        this.handleGameServerDisconnection();
    }

    connect = (ip, port, local) => {
        this.board.gameIsOver = false;
        this.ip = ip;
        this.port = port;
        this.loadImages();
        if (!local){
            this.conn = new WebSocket(`wss://up-and-down.io/join?ip=${ip}&port=${port}`)
        }
        else{
            this.conn = new WebSocket(`ws://${ip}:${port}/join`)
        }

        this.conn.onerror = (ev => {
            console.info("Websocket Connection Error");
            this.bigResetAndGameOverTrue();
        });

        this.conn.addEventListener("open", ev => {
            console.info("Websocket connected");
            this.handleGameServerConnection();
        })

        this.conn.addEventListener("close", ev => {
            console.info("Websocket closed");
            this.bigResetAndGameOverTrue();
        })

        this.conn.addEventListener("message", ev => {
            if (typeof ev.data !== "string") {
                this.updateTimeStamps();
                let fileReader = new FileReader();
                fileReader.readAsArrayBuffer(ev.data);
                fileReader.onload = (event) => {
                    let arrayBuffer = fileReader.result;
                    const view = new DataView(arrayBuffer);
                    let headerPosition = 0
                    while (headerPosition !== view.byteLength){
                        let messageHeader = view.getUint8(headerPosition);
                        //console.log(`Looking at message header: ${messageHeader} at headerPosition: ${headerPosition}`);
                            if (messageHeader === 0){
                                // Pieces Curves message
                                //console.log("Message was a curves message!");
                                const numBytesInMessage = view.getUint16(headerPosition+1);
                                const startByte = headerPosition+3
                                this.decodeCurvesMessage(view, startByte, numBytesInMessage);
                                headerPosition = headerPosition+numBytesInMessage
                                if (this.recieved_first_message === false){
                                    this.recieved_first_message = true;
                                }
                            }
                            if (messageHeader === 1){
                                // Player Health and Credits message
                                //console.log("Message was a Health and credits message!");
                                this.board.player1Health = view.getUint16(headerPosition+1);
                                this.board.player1Credits = view.getUint16(headerPosition+3);
                                this.board.player2Health = view.getUint16(headerPosition+5);
                                this.board.player2Credits = view.getUint16(headerPosition+7);
                                headerPosition = headerPosition + 9
                            }
                            if (messageHeader === 2){
                                // Game Over message
                                //console.log("Message was a game over message!");
                                let victor = view.getUint8(headerPosition+1);
                                headerPosition = headerPosition + 2;
                                this.board.victor = victor
                                let victoryMessage = ""
                                if (this.board.victor === this.board.player){
                                    victoryMessage = "You win!"
                                }
                                else{
                                    victoryMessage = "You lose"
                                }
                                this.handleGameServerDisconnection();
                                this.board.gameIsOver = true;
                                this.board.victoryMessage = victoryMessage;
                                //console.log(`Player ${victor+1} was victorious!`);
                            }
                        }
                };
                return
            }
            else{
                const data = JSON.parse(ev.data);
                //console.log("Data is: ");
                //console.log(data);
                if (data.hasOwnProperty("Height")){
                    // This is board data
                    //console.log("Player is: "+data.Player);
                    this.board.setDimensions(data.Width, data.Height);
                    this.board.setPlayer(data.Player);
                    this.cnvs.heightRatio = this.board.height / this.cnvs.height;
                    this.cnvs.widthRatio = this.board.width / this.cnvs.width;

                    // Let the game know players name
                    //console.log("PLAYERNAME "+this.board.player+" "+this.nickname);
                    this.send("PLAYERNAME "+this.board.player+" "+this.nickname);
                }
                if (data.hasOwnProperty("Player1Name")){
                    if (this.board.player === 0){
                        //console.log("Recieved enemy name: "+data.Player2Name)
                        this.board.setEnemyPlayerName(data.Player2Name)
                        this.setViewingGameScreen(true);
                    }
                    else{
                        //console.log("Recieved enemy name: "+data.Player1Name)
                        this.board.setEnemyPlayerName(data.Player1Name)
                        this.setViewingGameScreen(true);
                    }
                }
            }
        });
    }


    decodeCurvesMessage = (view, startByte, numBytesInMessage) =>{
        let currentByte = startByte
      
        //console.log(`Recieved bytes:${view.byteLength}`);
       
        //printUint8Bytes(view);
    
        while (currentByte!==numBytesInMessage){
            let pieceID = view.getUint16(currentByte)
            currentByte += 2
            let hitEnd = false;
            while (!hitEnd){
                let messageHeader = view.getUint8(currentByte);
                //console.log(`Message header code: ${messageHeader}`);
                if (messageHeader === 0){
                    
                    //console.log(`Message header concerned initial piece message`)
                    // initialPieceMsg
                    let piece = new Piece(
                        pieceID,view.getUint8(currentByte+1),
                        view.getUint8(currentByte+2),
                        [new CurveKeyFrame(view.getUint8(currentByte+3), view.getUint16(currentByte+4))],
                        [new CurveKeyFrame(view.getInt16(currentByte+6),view.getUint16(currentByte+8))],
                        [new CurveKeyFrame(view.getInt16(currentByte+10),view.getUint16(currentByte+12))],
                        [new CurveKeyFrame(view.getInt16(currentByte+14),view.getUint16(currentByte+16))]);
                    //console.log("New Piece Created: ");
                    //console.log(piece);
                    this.board.addPiece(piece);
                    currentByte = currentByte+18;
                }
                if (messageHeader === 4){
                    // position x curve key frame
                    //console.log(`Message header concerned x curve key frame`)
                    let keyframe = new CurveKeyFrame(view.getInt16(currentByte+1),view.getUint16(currentByte+3));
                    this.board.addPositionXCurveKeyFrame(pieceID, keyframe);
                    currentByte = currentByte + 5;
                }
                if (messageHeader === 5){
                    // position y curve key frame
                    //console.log(`Message header concerned y curve key frame`)
                    let keyframe = new CurveKeyFrame(view.getInt16(currentByte+1),view.getUint16(currentByte+3));
                    this.board.addPositionYCurveKeyFrame(pieceID, keyframe);
                    currentByte = currentByte + 5;
                    //console.log(`GOT NEW Y POSITION CURVE KEYFRAME FOR ${pieceID}`);
                    //console.log(keyframe);
                }
                if (messageHeader === 6){
                    // rotation curve key frame
                    //console.log(`Message header concerned y curve key frame`)
                    let keyframe = new CurveKeyFrame(view.getInt16(currentByte+1),view.getUint16(currentByte+3));
                    this.board.addRotationCurveKeyFrame(pieceID, keyframe);
                    //console.log(`GOT NEW ROTATION KEYFRAME FOR ${pieceID}`);
                    //console.log(keyframe);
                    currentByte = currentByte + 5;
                }
                if (messageHeader === 7){
                    // existence curve key frame
                    //console.log(`Message header concerned existence curve key frame`)
                    let keyframe = new CurveKeyFrame(view.getUint8(currentByte+1),view.getUint16(currentByte+2));
                    this.board.addExistenceCurveKeyFrame(pieceID, keyframe);
                    currentByte = currentByte + 4;
                }
                if (messageHeader === 255){
                    // hit the end of the message, go to the next piece
                    //console.log(`xx Hit end of message header xx`)
                    currentByte = currentByte + 1;
                    hitEnd = true
                }
            }
            /* DEBUG START*/
            let gotPiece = this.board.getPieceWithPieceID(pieceID);
            if (gotPiece!==null){
                if (gotPiece.identifier === 255){
                    //console.log(gotPiece);
                }
            }
            /* DEBUG END */ 
        }
        return currentByte;
    }

    handleRecievedImages = (responseText) => {
        //console.log("Recieved images")
        const imageData = JSON.parse(responseText);
        this.board.images = imageData;
        this.loadImages();
    }

    loadImages = () => {
        const capital_ship_player_one_image = new Image();
        capital_ship_player_one_image.src = "/images/ships/capital_ship_player_one_v2.png";

        const capital_ship_player_two_image = new Image();
        capital_ship_player_two_image.src = "/images/ships/capital_ship_player_two_v2.png";

        const interceptor_player_one_image = new Image();
        interceptor_player_one_image.src = "/images/ships/interceptor_player_one_v2.png";

        const interceptor_player_two_image = new Image();
        interceptor_player_two_image.src = "/images/ships/interceptor_player_two_v2.png";

        const destroyer_player_one_image = new Image();
        destroyer_player_one_image.src = "/images/ships/destroyer_player_one_v3.png";

        const destroyer_player_two_image = new Image();
        destroyer_player_two_image.src = "/images/ships/destroyer_player_two_v3.png";

        const flack_ship_player_one_image = new Image();
        flack_ship_player_one_image.src = "/images/ships/flack_ship_player_one_v2.png";

        const flack_ship_player_two_image = new Image();
        flack_ship_player_two_image.src = "/images/ships/flack_ship_player_two_v2.png";

        const laser_player_one_image = new Image();
        laser_player_one_image.src = "/images/projectiles/laser_player_one.png";

        const laser_player_two_image = new Image();
        laser_player_two_image.src = "/images/projectiles/laser_player_two.png";

        const big_laser_player_one_image = new Image();
        big_laser_player_one_image.src = "/images/projectiles/big_laser_player_one.png";

        const big_laser_player_two_image = new Image();
        big_laser_player_two_image.src = "/images/projectiles/big_laser_player_two.png";

        const flack_laser_player_one_image = new Image();
        flack_laser_player_one_image.src = "/images/projectiles/flack_laser_player_one.png";

        const flack_laser_player_two_image = new Image();
        flack_laser_player_two_image.src = "/images/projectiles/flack_laser_player_two.png";



        this.board.images = [
            { Player:0, Identifier: 0, image: capital_ship_player_one_image },
            { Player:1, Identifier: 0, image: capital_ship_player_two_image },
            { Player:0, Identifier: 1, image: interceptor_player_one_image }, 
            { Player:1, Identifier: 1, image: interceptor_player_two_image },
            { Player:0, Identifier: 2, image: destroyer_player_one_image }, 
            { Player:1, Identifier: 2, image: destroyer_player_two_image }, 
            { Player:0, Identifier: 3, image: flack_ship_player_one_image }, 
            { Player:1, Identifier: 3, image: flack_ship_player_two_image }, 
            { Player:0, Identifier: 255, image: laser_player_one_image }, 
            { Player:1, Identifier: 255, image: laser_player_two_image }, 
            { Player:0, Identifier: 254, image: big_laser_player_one_image }, 
            { Player:1, Identifier: 254, image: big_laser_player_two_image }, 
            { Player:0, Identifier: 253, image: flack_laser_player_one_image }, 
            { Player:1, Identifier: 253, image: flack_laser_player_two_image }, 

        ]



/*
        //console.log("Loading images...")
        for (let i=0;i<this.board.images.length;i++){
            this.board.images[i].image = new Image();
            this.board.images[i].image.src = ;
            //console.log("Loaded image @: ");
            //console.log(this.board.images[i].image.src);
        }
        //console.log("Loaded images")
        */
    }

    httpGetAsync = (theUrl, callback) => {
        //console.log("Getting url: "+theUrl)
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() { 
            if (xmlHttp.readyState === 4 && xmlHttp.status === 200)
                callback(xmlHttp.responseText);
        }
        xmlHttp.open("GET", theUrl, true); // true for asynchronous 
        xmlHttp.send(null);
    }
}

export default ConnectionManager