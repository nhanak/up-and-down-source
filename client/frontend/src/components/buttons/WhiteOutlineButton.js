import React from 'react';
import './WhiteOutlineButton.css'
import './MediumButton.css'

class WhiteOutlineButton extends React.Component {
    render() {
        const hoverGrey = this.props.hoverGrey
        if (!hoverGrey){
            return <button className="white-outline-button medium-button" type={this.props.type} onClick={this.props.onClick}>{this.props.label}</button>;
        }
        else{
            return <button className="white-outline-button hover-grey medium-button" type={this.props.type} onClick={this.props.onClick}>{this.props.label}</button>;
        }
    }
}

export default WhiteOutlineButton