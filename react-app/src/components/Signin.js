import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';
import { Link } from 'react-router-dom';

export default class Signin extends Component {
  constructor(props) {
    super(props);
    this.state = {
      text: "",
      usernameValue: false,
      passwordValue: false,
      isDisabled: () => { 
        if (this.state.usernameValue === false) {
          return true
        } else if (this.state.passwordValue === false) {
          return true
        } else {
          return false
        }
      }
    };
    this.submitForm = this.submitForm.bind(this);
    this.changeUsername = this.changeUsername.bind(this);
    this.changePassword = this.changePassword.bind(this);
  };

  submitForm() {
    let options = { method: "post" }
    fetch('/ui/signin?response_type=code&client_id={client_id}&state={state}&redirect_uri={redirect_uri}', options)
      .then((response) => {
          if ( response.ok ) {
            window.location.replace('/ui//authorize?response_type=code&client_id={client_id}&state={state}&redirect_uri={redirect_uri}')
          } else {
            this.setState({ text: "Something went wrong :(" });
          }
        }
      )
  };

  changeUsername = (e) => {
    if ( e.target.value.length > 0 ) {
      this.setState({ usernameValue: true })
    } else {
      this.setState({ usernameValue: false })
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
        <TextField label="Username" margin="normal" onChange={this.changeUsername} fullWidth />
        <TextField label="Password" margin="normal" onChange={this.changePassword} type="password" fullWidth />
        <span style={{ color: "red" }}>{this.state.text}</span>
        <Button 
          variant="raised" 
          color="primary" 
          style={{ padding:"10px 30px", margin:"15px 0 10px 0" }}
          onClick={this.submitForm}
          disabled={this.state.isDisabled()}
        >
          submit
        </Button>
        <Link to="/password/request" className="forget-password-link"> Forgot your password? </Link>
      </form>
    )
  }
};
