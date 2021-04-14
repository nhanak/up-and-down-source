import React from 'react';
import WhiteBoxInput from '../components/inputs/WhiteBoxInput';
import PlayButton from '../components/buttons/PlayButton';
import { useAuth0 } from '../contexts/auth0-context';
import "./Home.css";
import "./Page.css";
import axios from 'axios';

class Home extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            lookingForMatch:false,
            connectedToMatch:false,
        }
    }

    onClick = () =>{
        console.log("Play button clicked!");
        if ((this.state.lookingForMatch === false) && (this.state.connectedToMatch!==true)){
            this.setState({lookingForMatch:true})
            axios.get(window.location.href+'/matchmaker?nickname=satan')
            .then(res => {
                var connectionSplit = res.data.split('"');
                var ipAndPort = connectionSplit[1];
                var ipAndPortArray = ipAndPort.split(':');
                var ip = ipAndPortArray[0];
                var port = ipAndPortArray[1];
                console.log("IP: "+ip+ " Port: "+port);
                var conn = new WebSocket(`ws://${ip}:${port}/subscribe`)

                conn.onerror = function(ev){
                    console.log("Websocket Connection Error");
                    this.setState({
                        lookingForMatch:false,
                        connectedToMatch:false,
                    })

                }

                conn.addEventListener("open", ev => {
                    console.info("Websocket connected");
                    this.setState({
                        lookingForMatch:false,
                        connectedToMatch:true,
                    })
                })

                conn.addEventListener("close", ev => {
                    console.info("Websocket closed");
                    this.setState({
                        lookingForMatch:false,
                        connectedToMatch:false,
                    })
                })
            })
        }
    }

    render(){
        const { isLoading, user } = useAuth0();
        const {lookingForMatch, connectedToMatch} = this.state;
        return (
            <div className="home-page">
                {!user && (
                    <div className="standard-page">
                        <h1 className="home-splash-text">up & down</h1>
                        <p className="tip-text">Rock, paper scissors for the 21st century</p>
                        <div className="nickname-div">
                            <WhiteBoxInput placeholder="Nickname" size="large"/>
                        </div>
                    
                        <div className="play-button-div">
                            <PlayButton onClick={this.onClick}/>
                        </div>
                        <div className="search-status-div">
                            {lookingForMatch && (
                                <p>Looking for match...</p>
                            )}
                            {connectedToMatch && (
                                <p>Found and connected to match!</p>
                            )}
                        </div>
                    </div>
                )}
                    {user && (
                        <div>
                    <h1 className="home-splash-text">Welcome back,</h1>
                    <h1 className="home-splash-text">{user.nickname}</h1>
                    </div>
                    )}
            </div>
        );
    }
}

export default Home;