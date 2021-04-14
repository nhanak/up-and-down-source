import React from 'react';
import "./Home.css";
import "./Page.css";
import GameScreenOffline from '../components/gamescreenoffline/GameScreenOffline';

class GameScreenPageOffline extends React.Component {
    constructor(props){
        super(props);
    }

    render(){
        return (
            <div className="home-page">
                <div className="game-page">
                    <GameScreenOffline/>
                </div>
            </div>
            );
    }
}

export default GameScreenPageOffline