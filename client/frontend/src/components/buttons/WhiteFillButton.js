import React from 'react';
import './WhiteFillButton.css'
import './MediumButton.css'

class WhiteFillButton extends React.Component {
    render() {
        return <button className="white-fill-button medium-button" onClick={this.props.onClick}>{this.props.label}</button>;
    }
}

export default WhiteFillButton