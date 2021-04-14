import React from 'react';
import "./Home.css";
import "./Page.css"

function WhySignup() {
    return (
        <div className="how-to-play">
            <div className="standard-page">
                <h1 className="home-splash-text">Why signup?</h1>
                <div className="standard-text-container">
                    <h2 className="page-subheader">Choose your username</h2>
                    <p className="standard-text">Anonymous players have a randomly generated names. By signing up you can pick a username that you prefer.</p>
                    <h2 className="page-subheader">Track your progress</h2>
                    <p className="standard-text">up & down is a competitive game. With an account you can track your match making rating (MMR), see your match history as well as see your position on leaderboards.</p>
                    <h2 className="page-subheader">Lobby based matchmaking</h2>
                    <p className="standard-text">Down the road, we hope to implement lobby based matchmaking, so that you can play against friends. This will only be available to people who have accounts.</p>
                </div>
            </div>
        </div>
    );
}

export default WhySignup;