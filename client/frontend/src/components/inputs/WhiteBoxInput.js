import React from 'react';
import './WhiteBoxInput.css'

function WhiteBoxInput(props) {

    function handleChange(event){
        props.onChange(event);
    }

    const size = props.size;
    if (size === "large"){
        return <input className="white-box-input white-box-input-large" value={props.value} placeholder={props.placeholder} onChange={handleChange}/>;
    }
    if (size === "medium"){
        return <input className="white-box-input white-box-input-medium" value={props.value} placeholder={props.placeholder} onChange={handleChange}/>;
    }
    else{
        return <input className="white-box-input" value={props.value} placeholder={props.placeholder} onChange={handleChange}/>;
    }
}

export default WhiteBoxInput