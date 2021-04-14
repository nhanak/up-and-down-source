import React from 'react';
import "./CountdownTimer.css"

function CountdownTimer(props){
    let properHeight = props.containerHeight*0.5;
    let countdownContainerClasses = "countdown-container hideable "+props.countdownContainerAnimation;
    let versusContainerClasses = "versus-container hideable "+props.versusContainerAnimation
    let timedCountdownContainerClasses = "timed-countdown-container hideable "+props.timedCountdownContainerAnimation;
    let numberThreeClasses = "number three hideable "+ props.numberAnimation;
    let numberTwoClasses = "number two hideable "+ props.numberAnimation;
    let numberOneClasses = "number one hideable "+ props.numberAnimation;
    let numberGoClasses = "number go hideable "+ props.numberAnimation;
    return(
        <div id="countdownContainer" className={countdownContainerClasses} style={{width:props.containerWidth, height:props.containerHeight}}>
            <div id="versusContainer" className={versusContainerClasses}>
                <p className="match-start-header">{props.opponentName}</p>
                <p className="match-start-header">vs.</p>
                <p className="match-start-header">{props.playerName}</p>
            </div>
            <div id="timedCountdownContainer" className={timedCountdownContainerClasses}>
                <p className="match-starting-header">Match starts in...</p>
                <div className="circle-counting-down-container" style={{width:props.containerWidth, height:props.containerHeight}}>
                    <img style={{maxWidth:props.containerWidth, maxHeight:properHeight}} className="countdown-image" src="/rock_paper_scissors_v5.png"></img>
                    <p id="numberThree" className={numberThreeClasses}>3</p>
                    <p id="numberTwo" className={numberTwoClasses}>2</p>
                    <p id="numberOne" className={numberOneClasses}>1</p>
                    <p id="numberGo" className={numberGoClasses}>Go!</p>
                </div>
            </div>
        </div>
    )
}

export default CountdownTimer;

