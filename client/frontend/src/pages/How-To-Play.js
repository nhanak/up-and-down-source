import React from 'react';
import "./Home.css";
import "./Page.css";
import "./How-To-Play.css";

function HowToPlay() {
    return (
        <div className="how-to-play">
            <div className="standard-page">
                <h1 className="page-header">How to Play</h1>
                <div className="standard-text-container">
                <p className="standard-text underline">tl;dr: Press the <span className="green-text">green</span> highlighted buttons at the bottom of the game screen to purchase units. Counter the units your opponent is making so that your units can attack your opponents base and make their health fall to zero.</p>
                    <p className="standard-text">up & down is a Real Time Strategy (RTS) game where the objective is to have your units destroy the enemy base. Unlike most RTS games, in up & down you do not directly control your units, instead <i>they</i> decide what they should shoot and how to maneuver.</p>
                    <p className="standard-text">The strategy part of up & down is what units you decide to purchase for ¥ in order to counter what your opponent is making.</p>
                    <p className="standard-text">The units follow a rock, paper, scissors type counter system (see below). Counter what your opponent is making and win!</p>
                    
                </div>
                <div className="counter-chart-div">
                <h2 className="page-subheader">Counter Chart</h2>
                <img width="250px"src="/rock_paper_scissors_v5.png"></img>
                </div>
            </div>
        </div>
    );
}

export default HowToPlay;