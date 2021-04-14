import React from "react"
import "./UnitPurchaseButton.css"
class UnitPurchaseButton extends React.Component{
    constructor(props){
        super(props)
    }
    render(){
        //let howMuchToFillPercentage = Math.round((this.props.playerCredits/this.props.price)*100)
        let howMuchToFillPercentage =(this.props.playerCredits/this.props.price)*100;
        if (howMuchToFillPercentage > 100) {
            howMuchToFillPercentage = 100;
        }
        let greenClass = ""
        if (howMuchToFillPercentage === 100){
            greenClass = " green";
        }
        let stringFillPercentage = howMuchToFillPercentage.toString(10)+"%";
        return(
            <div className="unit-purchase-button-wrapper">
                <div className={"unit-purchase-button-border-thickener"+greenClass} style={{"height" : stringFillPercentage}}/>
                <button class={"unit-purchase-button"+greenClass} onClick={this.props.onClick}>
                    <div className="shortcut-and-name-wrapper">
                        <div class="shortcut-div">{this.props.shortcut}</div>
                        <div class="unit-name">{this.props.unitName}</div>
                    </div>
                    <img className="unit-image" src={this.props.unitImgSrc}></img>
                   
                    <div class="box">
                        <div className="price-row">
                        <img className="price-image" alt={this.props.imgAlt} src={this.props.yenImgSrc}></img>

                        <div class="price">{this.props.price}</div>
                        </div>
                    </div>
                </button>

            </div>
        );
    }
}

export default UnitPurchaseButton;

/*<img className="logo-nav-image" alt="updown game logo" src="/updown_logo_combined.png"> */