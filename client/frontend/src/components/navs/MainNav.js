import React from 'react';
import WhiteOutlineButton from '../buttons/WhiteOutlineButton';
import WhiteHyperlink from '../hyperlinks/WhiteHyperlink';
import WhiteFillButton from '../buttons/WhiteFillButton';
import { useAuth0 } from '../../contexts/auth0-context';
import { Link } from 'react-router-dom'
import './MainNav.css';
import { HamburgerSpin } from 'react-animated-burgers'
//https://github.com/AuvikAlive/react-animated-burgers
class MainNav extends React.Component {
    //const { isLoading, user, loginWithRedirect, logout} = useAuth0();
    constructor(props){
        super(props)
        this.state={isMenuOpen:false,}
    }

    menuButtonClick = () =>{
        const {isMenuOpen} = this.state;
        this.setState({isMenuOpen:!isMenuOpen});
    }

    logoImageClick = () => {
        if (this.props.isMenuOpen){
            this.props.toggleMenuOpen();
        }
    }

    render(){
        return(
            <div className="main-nav ">
                <div className="app-container">
                    <Link to="/">
                        <img onClick={this.logoImageClick} className="logo-nav-image" alt="updown game logo" src="/updown_logo_combined.png"></img>
                    </Link>
                        <div className="main-nav-end-div">
                            <div className="end-col-menu">
                                <HamburgerSpin className="menu-button-additional-styles" isActive={this.props.isMenuOpen} toggleButton={this.props.toggleMenuOpen} barColor="white" buttonWidth={30}/>
                            </div>
                        </div>

                    </div>
            </div>);
    }
}

export default MainNav;


/*
<WhiteHyperlink label="Sign in" onClick={loginWithRedirect}/>
<WhiteOutlineButton hoverGrey={true} label="Sign up" onClick={loginWithRedirect}/>
<WhiteOutlineButton hoverGrey={true} label="Logout" onClick={()=>logout({returnTo:"http://localhost:8000/"})}/>


{!isLoading && !user && (
    <Link to="/why-signup">
        <WhiteHyperlink label="Why sign up?"/>
    </Link>
)}


                    {!isLoading && user && (
                        <div className="main-nav-end-div">
                            <div className="main-nav-col">
                            <WhiteOutlineButton hoverGrey={true} label="Logout"/>
                            </div>
                        </div>
                    )}




                                                <div className="main-nav-col end-col-text">
                            <WhiteHyperlink label="Sign in"/>
                            </div>
                            <div className="end-col-button">
                                <WhiteOutlineButton hoverGrey={true} label="Sign up"/>
                            </div>
*/