import React from 'react';
import SignInForm from '../components/forms/SignInForm';
import "./Home.css"
import "./Login.css"

class Login extends React.Component {
    render() {
        return(
            <div className="login-page">
                <div className="splash-div">
                    <h1 className="home-splash-text">Sign in</h1>
                    <SignInForm/>
                </div>
            </div>);
    }
}

export default Login;