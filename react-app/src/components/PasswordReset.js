import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';
import { Redirect } from 'react-router-dom';

export default class PasswordRequest extends Component {
	constructor(props) {
    super(props);
    this.state = {
      text: "",
      isLoaded: false
    };
    this.submitForm = this.submitForm.bind(this);
  };

  submitForm() {
    let options = { method: "post" }
    fetch('/ui/password/request', options)
      .then((response) => {
          if ( response.ok ) {
            this.setState({ isLoaded: true });
          } else {
            this.setState({ text: "Something went wrong :(" });
          }
        }
      )
	}
	
	render() {
  	return (
    	<form className="form">
				<TextField label="Verification code" margin="normal" type="password" fullWidth required/>
				<TextField label="New password" margin="normal" type="password" fullWidth required/>
				<span style={{ color: "red" }}>{this.state.text}</span>
				<Button 
					variant="raised" 
					color="primary" 
					style={{ padding:"10px 30px", marginTop:"15px" }}
					onClick={this.submitForm}
				>
					submit
				</Button>
				{ this.state.isLoaded ? <Redirect to="/signin" push/> : null }
			</form>
    )
  }
};