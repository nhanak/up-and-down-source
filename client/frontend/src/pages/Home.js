import React from 'react';
import WhiteBoxInput from '../components/inputs/WhiteBoxInput';
import PlayButton from '../components/buttons/PlayButton';
import GameScreen from '../components/gamescreen/GameScreen';
import "./Home.css";
import "./Page.css";
import "../components/countdown/CountdownTimer";
import axios from 'axios';
import ConnectionManager from '../helpers/ConnectionManager';
var Filter = require('bad-words');
var badwordList = require('badwords-list')

class Home extends React.Component {
    constructor(props){
        super(props);
        let filter = new Filter();
        let bwlArray = badwordList.array;
        for (let i=0;i<bwlArray.length;i++){
            filter.addWords(bwlArray[i].toString());
        }
       
        this.state = {
            viewingGameScreen:false,
            lookingForMatch:false,
            connectedToMatch:false,
            gameserverIP:null,
            gameserverPort:null,
            haveGameServer:false,
            nickname: "",
            nicknameErrorInfo:"",
            errorInfo:"",
            badWordFilter:filter,
            connectionManager:null,
            tipOfTheDay:"",
        }
        this.local = true;
    }

    handleGameServerConnection = () => {
        this.setState({
            connectedToMatch:true,
        })
    }

    handleGameServerDisconnection = () => {
        this.setState({
            gameserverIP:null,
            gameserverPort:null,
            haveGameServer:false,
            connectedToMatch:false,
            lookingForMatch:false
        })
    }

    componentDidMount(){
        let tipsOfTheDay = [
            "Welcome to the super duper secret pre-pre-pre-alpha",
            "You did watch LOGH right anon?",
            "Starcraft for dummies",
            "Because you need to de-stress from League",
            "Because you need to de-stress from DOTA",
            "Because you need to de-stress from CS:GO",
            "Now with 1000% more Jajanken",
            "Now with 50% more salt",
            "If rock paper scissors was a RTS"
        ];
        const randomTip = tipsOfTheDay[Math.floor(Math.random() * tipsOfTheDay.length)];
        this.setState({tipOfTheDay:randomTip});
    }

    onClickPlayGameButton = () => {
        const {nickname} = this.state;
        //console.log("Play button clicked!");
        if ((this.state.lookingForMatch === false) && (this.state.connectedToMatch !== true) && (this.validNickname(nickname))){
            this.setState({errorInfo:"", searchInfo:"Looking for a match...", lookingForMatch:true});
            if (this.local){
                var ip = "127.0.0.1";
                var port = "7654";
                //console.log("IP: "+ip+ " Port: "+port);
                let finalNickname = nickname;
                if (finalNickname === ""){
                    finalNickname = this.generateNickname();
                }
                let connectionManager = new ConnectionManager(finalNickname, this.setViewingGameScreen, this.handleGameServerConnection, this.handleGameServerDisconnection);
                connectionManager.connect(ip, port, this.local);
                this.setState({nickname:finalNickname,gameserverIP:ip, gameserverPort:port, haveGameServer:true, connectionManager:connectionManager});
                
            }
            else{
                // console.log("Making Axios Request!");
                axios.get(window.location.href+'/matchmaker')
                .then(res => {
                // console.log(res.data);
                    var ip = res.data.IP;
                    var port = res.data.Port;
                    //console.log("IP: "+ip+ " Port: "+port);
                    
                    let finalNickname = nickname;
                    if (finalNickname === ""){
                        finalNickname = this.generateNickname();
                    }
                    let connectionManager = new ConnectionManager(finalNickname, this.setViewingGameScreen, this.handleGameServerConnection, this.handleGameServerDisconnection);
                    connectionManager.connect(ip, port, this.local);
                    this.setState({nickname:finalNickname, gameserverIP:ip, gameserverPort:port, haveGameServer:true, connectionManager:connectionManager});
                    
                }).catch(err =>{
                    //console.log("There was an internal server error");
                    //console.log(err);
                    if (err.response.status === 500){
                        this.setState({errorInfo:err.response.data.Error});
                    }
                    else{
                        this.setState({errorInfo:"There was an internal server error"});
                    }
                    this.setState({lookingForMatch:false})
                });
            }
        }
    }

    setViewingGameScreen = (isViewing) => {
        this.setState({viewingGameScreen:isViewing});
    }

    generateNickname = () => {
        const randomString = this.generateRandomString(5);
        const nickname = "anon-"+randomString;
        //console.log("Generated nickname is: "+nickname);
        return nickname
    }

    handleNicknameChange = (event) =>{
        //console.log("handle nickname change called")
        const newNickname = event.target.value;
        if (this.validNickname(newNickname)){
            this.setState({nickname:newNickname})
        }
    }

    validNickname = (nickname) => {
        let valid = true;
        let trimmedNickname = nickname.split(' ').join('');
        if (trimmedNickname.length<nickname.length){
            //console.log("invalid space char");
            this.setState({nicknameErrorInfo:"Nicknames cannot have space or tab characters"})
            valid = false;
            return valid;
        }

        //console.log(`passed nickname: ${nickname} length: ${nickname.length}`)

        if (trimmedNickname === ""){
            this.setState({nicknameErrorInfo:""});
            return valid;
        }
        if (trimmedNickname.length>20){
           // console.log("invalid too long");
            this.setState({nicknameErrorInfo:"Nicknames cannot be longer that 20 characters"})
            valid = false;
            return valid;
        }
        if (trimmedNickname.length <= 2){
            //console.log("invalid too short");
            this.setState({nicknameErrorInfo:"Nicknames must be 3 characters or longer"})
            valid = false;
            return valid;
        }

        const {badWordFilter} = this.state;
        let profane = badWordFilter.isProfane(trimmedNickname);

        if (profane){
            //console.log("invalid profane");
            this.setState({nicknameErrorInfo:"Nickname is profane"})
            valid = false;
            return valid;
        }
        //console.log(`Nickname is valid: ${valid}`);
        this.setState({nicknameErrorInfo:""});
        return valid
    }

    generateRandomString = (length) =>{
        var x = [...Array(length)].map(i=>(~~(Math.random()*36)).toString(36)).join('')
        return x
    }

    render(){
        const {
            lookingForMatch, 
            errorInfo,
            searchInfo,
            viewingGameScreen,
            connectionManager,
            nickname,
            tipOfTheDay,
            nicknameErrorInfo
         } = this.state;

        return (
            <div className="home-page">
                {!viewingGameScreen && (
                    <div className="flex-center-page-wrapper">
                        <div className="standard-page">
                            <h1 className="home-splash-text">up & down</h1>
                            <p className="tip-text">{tipOfTheDay}</p>
                            <div className="nickname-div">
                                <WhiteBoxInput placeholder="Nickname" size="large" onChange={this.handleNicknameChange}/>
                                <p className="nickname-error-text">{nicknameErrorInfo}</p>
                            </div>

                            <div className="play-button-div">
                                <PlayButton onClick={this.onClickPlayGameButton}/>
                            </div>
                            <div className="search-status-div">
                                {lookingForMatch && (
                                    <p className="tip-text">{searchInfo}</p>
                                )}
                                {!lookingForMatch && (
                                    <p className="error-text">{errorInfo}</p>
                                )}
                            </div>
                        </div>
                    </div>
                )}
                {viewingGameScreen && (
                    <div className="game-page">
                        <GameScreen nickname={nickname} lookingForMatch = {lookingForMatch} disconnect = {connectionManager.disconnect} getLatency = {connectionManager.getLatency} send={connectionManager.send} searchInfo={searchInfo} errorInfo={errorInfo} board={connectionManager.getBoard()} cnvs={connectionManager.getCnvs()} onClickPlayAgainButton={this.onClickPlayGameButton} setViewingGameScreen={this.setViewingGameScreen} hideFooter={this.props.hideFooter}/>
                    </div>
                )}
            </div>
        );
    }
}

export default Home;