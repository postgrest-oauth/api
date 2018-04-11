import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';
import { Redirect } from 'react-router-dom';

export default class PasswordRequest extends Component {
	constructor(props) {
    super(props);
    this.state = {
      text: "",
      isLoaded: false,
      codeValue: false,
      passwordValue: false,
      isDisabled: () => { 
        if (this.state.codeValue === false) {
          return true
        } else if (this.state.passwordValue === false) {
          return true
        } else {
          return false
        }
      }
    };
    this.submitForm = this.submitForm.bind(this);
    this.changeCode = this.changeCode.bind(this);
    this.changePassword = this.changePassword.bind(this);
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
  
  changeCode = (e) => {
    if ( e.target.value.length > 0 ) {
      this.setState({ codeValue: true })
    } else {
      this.setState({ codeValue: false })
    }
  };

  changePassword = (e) => {
    if ( e.target.value.length > 0 ) {
      this.setState({ passwordValue: true })
    } else {
      this.setState({ passwordValue: false })
    }
  };
	
	render() {
  	return (
    	<form className="form">
				<TextField label="Verification code" margin="normal" onChange={this.changeCode} fullWidth />
				<TextField label="New password" margin="normal" type="password" onChange={this.changePassword} fullWidth />
				<span style={{ color: "red" }}>{this.state.text}</span>
				<Button 
					variant="raised" 
					color="primary" 
					style={{ padding:"10px 30px", marginTop:"15px" }}
          onClick={this.submitForm}
          disabled={this.state.isDisabled()}
				>
					submit
				</Button>
				{ this.state.isLoaded ? <Redirect to="/signin" push/> : null }
			</form>
    )
  }
};