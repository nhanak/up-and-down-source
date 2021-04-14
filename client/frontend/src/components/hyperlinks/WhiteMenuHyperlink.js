import React from 'react'
import "./WhiteMenuHyperlink.css"

function WhiteMenuHyperlink(props){
    return(<a onClick={props.onClick} className="white-menu-hyperlink">{props.label}</a>);
}

export default WhiteMenuHyperlink;