import React from 'react'
import "./WhiteHyperlink.css"

function WhiteHyperlink(props){
    return(<a onClick={props.onClick} className="white-hyperlink">{props.label}</a>);
}

export default WhiteHyperlink;