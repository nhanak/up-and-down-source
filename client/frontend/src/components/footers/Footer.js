
import WhiteHyperlink from '../hyperlinks/WhiteHyperlink';
import React from "react";
import { Link } from 'react-router-dom'
import "./Footer.css";

function Footer(props){
    return(
    <div className="footer">
        <div className="app-container">
            <div className="footer-row">
                <div className="made-with-love-container size-single-flex">
                    <p className="made-with-love">Made with</p>
                    <img className="heart-img" src="/heart_inverted.png"/>
                    <p className="made-with-love">in Canada</p>
                </div>
                <div className="footer-privacy-container size-single-flex">
            
                <div className="footer-item-left">
                    <Link to="/privacy">
                        <WhiteHyperlink label="Privacy"></WhiteHyperlink>
                    </Link>
                </div>
                <p className="made-with-love">-</p>
                <div className="footer-item-right">
                    <Link  to="/contact">
                        <WhiteHyperlink label="Contact"></WhiteHyperlink>
                    </Link>
                </div>
                </div>
                <div className="socials-container size-single-flex">
                    <a href="https://twitter.com/intent/tweet?url=www.up-and-down.io&text=Come play up %26 down!" target="_blank">
                        <img className="social-image" src="/twitter-52.png"/>
                    </a>
                    <a href="https://www.facebook.com/sharer/sharer.php?u=up-and-down.io" target="_blank">
                        <img className="social-image" src="/facebook-52.png"/>
                    </a>
                </div>
            </div>
        </div>
    </div>);
}

export default Footer;






