import React from 'react';
import './GreyUnderlineInput.css'

function GreyUnderlineInput(props) {

    function handleChange(event){
        props.onChange(event.target.value);
    }

    const size = props.size;
    if (size === "large"){
        return <input className="grey-underline-input grey-underline-input-large" value={props.value} placeholder={props.placeholder} onChange={handleChange}/>;
    }
    if (size === "medium"){
        return <input className="grey-underline-input grey-underline-input-medium" value={props.value} placeholder={props.placeholder} onChange={handleChange}/>;
    }
    else{
        return <input className="grey-underline-input" value={props.value} placeholder={props.placeholder} onChange={handleChange}/>;
    }
}

export default GreyUnderlineInput