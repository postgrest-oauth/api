import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';
import { Redirect } from 'react-router-dom';

export default class PasswordRequest extends Component {
  constructor(props) {
    super(props);
    this.state = {
      text: "",
      isLoaded: false,
      inputValue: false,
      isDisabled: () => { 
        if (this.state.inputValue === false) {
          return true
        } else {
          return false
        }
      }
    };
    this.submitForm = this.submitForm.bind(this);
    this.changeInput = this.changeInput.bind(this);
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

  changeInput = (e) => {
    if ( e.target.value.length > 0 ) {
      this.setState({ inputValue: true })
    } else {
      this.setState({ inputValue: false })
    }
  };
  
	render() {
  	return (
    	<form className="form">
				<TextField label="Email or phone" margin="normal" onChange={this.changeInput} fullWidth />
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
        { this.state.isLoaded ? <Redirect to="/password/reset" push/> : null }
			</form>
    )
  }
};