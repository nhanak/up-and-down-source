import React from 'react';
import CreateAccountForm from '../components/forms/CreateAccountForm';
import "./Home.css"

class SignUp extends React.Component {
    render() {
        return(
            <div className="login-page">
                <div className="splash-div">
                    <h1 className="home-splash-text">Create Account</h1>
                    <CreateAccountForm/>
                </div>
            </div>);
    }
}

export default SignUp;