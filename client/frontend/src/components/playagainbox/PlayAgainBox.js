import React from 'react';
import "./PlayAgainBox.css"
import "../buttons/WhiteOutlineButton"
import WhiteOutlineButton from '../buttons/WhiteOutlineButton';


function PlayAgainBox(props){
    let victoryMessage = props.victoryMessage;
    if ((victoryMessage === "")&&(!props.lookingForMatch)){
        victoryMessage = "Menu"
    }
    return(
        <div className="play-again-box-container">
            {!props.lookingForMatch && (
                <p className="victory-message">{victoryMessage}</p>
            )}
            {!props.lookingForMatch && (
                <div>
                    <div className="play-again-button-row">
                        <WhiteOutlineButton onClick={props.onClickPlayAgain} hoverGrey={true} label="[R] Play again"/>
                    </div>
                    <div className="play-again-bottom-row">
                        <WhiteOutlineButton onClick={props.onClickHome} hoverGrey={true} label="Home"/>
                    </div>
                </div>
             )}
            <div className="search-status-div">
                {props.lookingForMatch && (
                    <p className="search-text">{props.searchInfo}</p>
                )}
                {!props.lookingForMatch && (
                    <p className="error-text-play-again-box">{props.errorInfo}</p>
                )}
            </div>
        </div>
    );
}

export default PlayAgainBox