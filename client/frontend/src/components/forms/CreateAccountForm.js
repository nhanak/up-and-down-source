import React from 'react';
import WhiteOutlineButton from '../buttons/WhiteOutlineButton';
import CreateAccountFormInput from '../inputs/CreateAccountFormInput';
import { Link } from 'react-router-dom';
import "./CreateAccountForm.css";
import "./Form.css";
import {FormFieldValidity, FormSubmissionValidityObject} from "./FormSubmissionDataStructures";

class CreateAccountForm extends React.Component {
    constructor(props){
        super(props);   
        this.state={
            email: "",
            username: "",
            password: "",
            reenteredPassword:"",
            emailValid:false,
            usernameValid:false,
            passwordValid:false,
            reenteredPasswordValid:false,
            formValid:false,
            emailErrors:[],
            userNameErrors:[],
            passwordErrors:[],
            reenteredPasswordErrors:[],
            formErrors:[],
        }
    }

    handleSubmit = (event) => {
        event.preventDefault();
        this.validateForm();
        if (!this.showFormError){
            // Actually send the form
        }
    }

    validateForm = () => {
        const isValidEmail = this.validateEmail();
        const isValidUsername = this.validateUsername();
        const isValidPassword = this.validatePassword();
        const isValidReenteredPassword = this.validateReenteredPassword();

        const fieldValidityArray = [isValidEmail, isValidUsername, isValidPassword, isValidReenteredPassword];

        let isValidForm = true;
        for (let i=0;i<fieldValidityArray.length;i++){
            const isValid = fieldValidityArray[i];
            if (!isValid){
                isValidForm = false;
            }
        }

        this.setState({
            emailValid: isValidEmail,
            usernameValid:isValidUsername, 
            passwordValid: isValidPassword,
            reenteredPasswordValid: isValidReenteredPassword,
            formValid:isValidForm,
        })
    }

    validateEmail = () => {
        let isValid = false;
        let errors = [];
        const {email} = this.state;
        if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(email)){
            isValid = true;
        } 
        else{
            errors.push("Please enter a valid email")
        }
        if (errors.length === 0){
            // No errors, check database if the email is in use
        }
        this.setState({emailErrors: errors});
        return isValid;
    }

    validateUsername = () => {
        let isValid = true;
        const {username} = this.state;
        let errors = [];
        if (username.length<3){
            isValid = false;
            errors.push("Usernames must be longer than 2 characters");
        }
        if (username.length>20){
            isValid = false;
            errors.push("Usernames must be shorter than 20 characters");
        }
        if (!/^[A-Za-z]/.test(username)){
            isValid = false;
            errors.push("Usernames must start with a letter");
        }
        if (!/^[A-Za-z0-9]+$/.test(username)){
            isValid=false;
            errors.push("Usernames can only contain letters and numbers");
        }
        this.setState({userNameErrors:errors})
        return isValid;
    }

    validatePassword = () => {
        let isValid = true;
        const {password} = this.state;
        let errors = [];
        if (password.length<1){
            errors.push("Please enter a password")
            this.setState({passwordErrors: errors});
            return false;
        }
        if (password.length<8){
            errors.push("Passwords must be longer than 8 characters")
            isValid = false;
        }
        if (!/[a-z]/.test(password)){
            errors.push("Passwords must have at least one lower case character [a-z]");
            isValid = false;
        }
        if (!/[A-Z]/.test(password)){
            errors.push("Passwords must have at least one upper case character [A-Z]");
            isValid = false;
        }
        if (!/[!@#$%^&*]/.test(password)){
            errors.push("Passwords must have at least one special character [!@#$^&*]");
            isValid = false;
        }
        this.setState({passwordErrors:errors});
        return isValid;
    }

    validateReenteredPassword = () => {
        let isValid = true;
        const {password, reenteredPassword} = this.state;
        let errors = [];
        if (password !== reenteredPassword){
            isValid = false;
            errors.push("Passwords do not match")
        }
        this.setState({reenteredPasswordErrors:errors});
        return isValid;
    }

    handleEmailChange = (value) => {
        this.setState({email:value});
    }

    handleUsernameChange = (value) => {
        this.setState({username:value});
    }

    handlePasswordChange = (value) => {
        this.setState({password:value});
    }

    handleReenteredPasswordChange = (value) => {
        this.setState({reenteredPassword:value});
    }

    render() {
        return(
            <form className="create-account-form" autoComplete="off" onSubmit={this.handleSubmit}>
                <CreateAccountFormInput label="Email" value={this.state.email} onChange={this.handleEmailChange} placeholder="" errors={this.state.emailErrors}/>
                <CreateAccountFormInput label="Username" value={this.state.username} onChange={this.handleUsernameChange} placeholder="" errors={this.state.userNameErrors}/>
                <CreateAccountFormInput label="Password" value={this.state.password} onChange={this.handlePasswordChange} placeholder="" errors={this.state.passwordErrors}/>
                <CreateAccountFormInput label="Reenter Password" value={this.state.reenteredPassword} onChange={this.handleReenteredPasswordChange} placeholder="" errors={this.state.reenteredPasswordErrors}/>
                <div className="play-button-div">
                    <WhiteOutlineButton label="Sign up" hoverWhite={true} type="submit"/>
                </div>
            </form>);
    }
}

export default CreateAccountForm;
