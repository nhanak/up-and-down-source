import React from 'react';
import "./CreateAccountFormInput.css";
import GreyBoxInput from "./GreyBoxInput";

function CreateAccountFormInput(props){
    const errors = props.errors.map(error=>(
        <p className="form-field-error">{error}</p>
    ));

    let errorsDiv =<div/>
    if (props.errors.length>0){
    errorsDiv =(
        <div className="form-field-error-div">
            <p className="form-field-error-title">Error: </p>
            {errors}
        </div>);
    }

    return(
        <div className="create-account-form-input-row">
            <p className="create-account-form-label">{props.label}</p>
            <GreyBoxInput size="medium" value={props.value} placeholder={props.placeholder} onChange={props.onChange}/>
            {errorsDiv}
        </div>);
}

export default CreateAccountFormInput;
