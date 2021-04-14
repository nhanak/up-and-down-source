import React, { useContext } from 'react';
import Home from './pages/Home';
import GameScreenPageOffline from './pages/GameScreenPageOffline';
import HowToPlay from './pages/How-To-Play';
import WhySignup from './pages/Why-Signup';
import Privacy from './pages/Privacy';
import Contact from './pages/Contact';
import {Route, Switch } from 'react-router-dom';
import MainNav from './components/navs/MainNav';
import MainMenu from './components/menus/MainMenu';
import Footer from './components/footers/Footer'

import './App.css';

class App extends React.Component{
    constructor(props){
        super(props);
        this.state = {hideFooter:false, isMenuOpen:false,}   
    }

    hideFooter = (boolVal) => {
        this.setState({hideFooter:boolVal});
    }

    toggleMenuOpen = () =>{
        const {isMenuOpen} = this.state;
        this.setState({isMenuOpen:!isMenuOpen});
    }

    render(){
        const {hideFooter, isMenuOpen} = this.state;
        return (
            <div className="app">
                <div className="app-content">
                    <MainNav toggleMenuOpen={this.toggleMenuOpen} isMenuOpen={isMenuOpen}/>  
                    <div className="app-container">
                        <Switch>
                            <Route 
                                path="/" 
                                render={(props) => <Home {...props} hideFooter={this.hideFooter}/>} 
                                exact />
                        </Switch>
                        <Switch>
                                <Route path="/game-screen-page-offline" component={GameScreenPageOffline} exact />
                        </Switch>
                        <Switch>
                                <Route path="/how-to-play" component={HowToPlay} exact />
                        </Switch>
                        <Switch>
                                <Route path="/why-signup" component={WhySignup} exact />
                        </Switch>
                        <Switch>
                                <Route path="/privacy" component={Privacy} exact />
                        </Switch>
                        <Switch>
                                <Route path="/contact" component={Contact} exact />
                        </Switch>
                    </div>
                </div>
                {!hideFooter && (
                    <Footer/>
                )}
                <MainMenu toggleMenuOpen={this.toggleMenuOpen} isMenuOpen={isMenuOpen}/>
            </div>
        );
    }
}
export default App;
