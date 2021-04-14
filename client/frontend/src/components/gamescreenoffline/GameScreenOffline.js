import React from 'react';
import "./GameScreenOffline.css"
import UnitPurchaseButton from "../buttons/UnitPurchaseButton"
import PlayAgainBox from "../playagainbox/PlayAgainBox"

class GameScreenOffline extends React.Component{

    constructor(props){
        super(props);
    }

    render(){
        return (
            <div class="game-screen-wrapper">
                <div class="enemy-status-bar status-bar">
                    <div class="outer-status-bar-column">
                        <p class="player-name">Justin</p>
                    </div>
                    <div class="outer-status-bar-column">
                        <div class="resources-column">
                        <div class="outer-status-bar-column">
                                <div class="image-and-text-flex">
                                    <img className="health-image" src="/heart_inverted_dark.png"></img>
                                    <p>100</p>
                                </div>
                            </div>
                        <div class="outer-status-bar-column">
                                <div class="image-and-text-flex">
                                    <img className="health-image"  src="/yen_dark.png"></img>
                                    <p>1000</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <canvas id="mainCanvas" class="main-canvas"></canvas>
                <div class="player-status-bar status-bar">
                    <div class="outer-status-bar-column">
                        <p class="player-name">Neil</p>
                    </div>
                    <div class="outer-status-bar-column">
                        <div class="resources-column">
                        <div class="outer-status-bar-column">
                                <div class="image-and-text-flex">
                                    <img className="health-image" src="/heart_inverted.png"></img>
                                    <p>100</p>
                                </div>
                            </div>
                            <div class="outer-status-bar-column">
                                <div class="image-and-text-flex">
                                    <img className="health-image"  src="/yen.png"></img>
                                    <p>1000</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="controls-div">
                        <UnitPurchaseButton price={50} shortcut={"A"} imgAlt="Interceptor Image" unitImgSrc="/Interceptor.png" yenImgSrc="/yen_original.png" />
                        <UnitPurchaseButton price={200} shortcut={"S"} imgAlt="Flack Image" unitImgSrc="/Flack_Ship.png" yenImgSrc="/yen_original.png" />
                        <UnitPurchaseButton price={200} shortcut={"D"} imgAlt="Destroyer Image" unitImgSrc="/Destroyer.png" yenImgSrc="/yen_original.png" />
                </div>
                <PlayAgainBox victoryMessage={"You won!"}/>
            </div>
        );
    }
}

export default GameScreenOffline;

