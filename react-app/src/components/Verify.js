import React, { Component } from 'react';
import { TextField, Button, Typography } from 'material-ui';

export default class Verify extends Component {
  constructor(props) {
    super(props);
    this.state = {
      text: "",
      textColor: "",
      codeValue: false,
      isDisabled: () => { 
        if (this.state.codeValue === false) {
          return true
        } else {
          return false
        }
      }
    };
    this.submitForm = this.submitForm.bind(this);
    this.changeCode = this.changeCode.bind(this);
  };

  submitForm() {
    let options = { method: "post" }
    fetch('/ui/verify', options)
      .then((response) => {
          if ( response.ok ) {
            this.setState({ text: "Success! :)", textColor: "green" });
          } else {
            this.setState({ text: "Something went wrong :(", textColor: "red" });
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

  render() {
    return(
      <form className="form">
        <Typography>Please input verification code from email</Typography>
        <TextField label="Verification code" margin="normal" onChange={this.changeCode} fullWidth />
        <span style={{ color: this.state.textColor }}>{this.state.text}</span>
        <Button 
          variant="raised" 
          color="primary" 
          style={{ padding:"10px 30px", marginTop:"15px" }}
          onClick={this.submitForm}
          disabled={this.state.isDisabled()}
        >
          submit
        </Button>
      </form>
    )
  }
};