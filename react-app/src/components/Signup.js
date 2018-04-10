import React, { Component } from 'react';
import { TextField, Button, FormControl, Input, InputLabel } from 'material-ui';
import MaskedInput from 'react-text-mask';
import { Redirect } from 'react-router-dom';

function InputMask() {
  return (
    <MaskedInput
      mask={['+',/\d/,/\d/,/\d/,' ','(',/\d/,/\d/,')',' ',/\d/,/\d/,/\d/,'-',/\d/,/\d/,'-',/\d/,/\d/]}
      placeholder="+123 (45) 678-90-12"
      guide={false}
      style={{ border:"none", width:"100%", outline:"none", fontFamily:"Roboto", fontSize:"16px", padding:"5px 0" }}
    />
  )
}

export default class Signup extends Component {
  constructor(props) {
    super(props);
    this.state = {
      text: "",
      isLoaded: false,
      emailValue: false,
      passwordValue: false,
      isDisabled: () => { 
        if (this.state.emailValue === false) {
          return true
        } else if (this.state.passwordValue === false) {
          return true
        } else {
          return false
        }
      }
    };
    this.submitForm = this.submitForm.bind(this);
    this.changeEmail = this.changeEmail.bind(this);
    this.changePassword = this.changePassword.bind(this);
  };

  submitForm() {
    let options = { method: "post" }
    fetch('/ui/signup?response_type=code&client_id={client_id}&state={state}&redirect_uri={redirect_uri}', options)
      .then((response) => {
          if ( response.ok ) {
            this.setState({ isLoaded: true });
          } else {
            this.setState({ text: "Something went wrong :(" });
          }
        }
      )
  }

  changeEmail = (e) => {
    if ( e.target.value.length > 0 ) {
      this.setState({ emailValue: true })
    } else {
      this.setState({ emailValue: false })
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
        <TextField label="Email address" margin="normal" onChange={this.changeEmail} fullWidth />
        <TextField label="Password" margin="normal" type="password" onChange={this.changePassword} fullWidth />
        <FormControl margin="normal" fullWidth >
          <InputLabel shrink={true}> Phone number </InputLabel>
          <Input inputComponent={InputMask} />
        </FormControl>
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
        { this.state.isLoaded ? <Redirect to="/verify" push/> : null } 
      </form>
    )
  }
};
