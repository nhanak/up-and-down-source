import React from 'react';
import WhiteOutlineButton from '../buttons/WhiteOutlineButton';
import GreyUnderlineInput from '../inputs/GreyUnderlineInput';
import { Link } from 'react-router-dom';
import "./SignInForm.css";

class CreateAccountForm extends React.Component {
    render() {
        return(
            <div className="sign-in-form">
                <div>
                    <GreyUnderlineInput placeholder="username" size="large"/>
                </div>
                <div className="password-div">
                    <GreyUnderlineInput placeholder="password" size="large"/>
                </div>
                <div className="play-button-div">
                    <WhiteOutlineButton label="Sign in" hoverWhite={true}/>
                </div>
            </div>);
    }
}

export default CreateAccountForm;