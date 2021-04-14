import React from 'react';
import './WhiteOutlineButton.css'
import './LargeButton.css'

class PlayButton extends React.Component {
    render() {
            return <button className="white-outline-button hover-grey large-button" type={this.props.type} onClick={this.props.onClick}>Play</button>;
        }
    
}

export default PlayButton