
import React from 'react';
import "./MainMenu.css";
import WhiteMenuHyperlink from '../hyperlinks/WhiteMenuHyperlink';
import { Link } from 'react-router-dom'

function MainMenu(props){

    let mainMenuPageContainerClasses = "main-menu-page-container ";
    let mainMenuContainerClasses = "main-menu-container ";

    if (!props.isMenuOpen){
        mainMenuPageContainerClasses += "fade-out";
        mainMenuContainerClasses += "offscreen";
    }
    else{
        mainMenuPageContainerClasses += "fade-in";
    }

    return(
        <div className={mainMenuPageContainerClasses}>
            <div className={mainMenuContainerClasses}>
                <div className="menu-link-wrapper">
                    <Link to="/">
                        <WhiteMenuHyperlink onClick={props.toggleMenuOpen} label="Home"/>
                    </Link>
                </div>
                <div className="menu-link-wrapper">
                    <Link to="/how-to-play">
                        <WhiteMenuHyperlink onClick={props.toggleMenuOpen} label="How to Play"/>
                    </Link>
                </div>
                <div className="menu-link-wrapper">
                    <Link to="/contact">
                        <WhiteMenuHyperlink onClick={props.toggleMenuOpen} label="Contact"/>
                    </Link>
                </div>
                <div className="menu-link-wrapper">
                    <Link to="/privacy">
                        <WhiteMenuHyperlink onClick={props.toggleMenuOpen} label="Privacy"/>
                    </Link>
                </div>
            </div>
        </div>
    );
}

export default MainMenu;