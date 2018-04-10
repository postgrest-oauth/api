import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';
import { Link } from 'react-router-dom';

export default class Signin extends Component {
  constructor(props) {
    super(props);
    this.state = {
      text: ""
    };
    this.submitForm = this.submitForm.bind(this);
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
  }

  render() {
    return (
      <form className="form">
        <TextField label="Username" margin="normal" fullWidth required/>
        <TextField label="Password" margin="normal" type="password" fullWidth required/>
        <span style={{ color: "red" }}>{this.state.text}</span>
        <Button 
          variant="raised" 
          color="primary" 
          type="submit" 
          style={{ padding:"10px 30px", margin:"15px 0 10px 0" }}
          onClick={this.submitForm}
        >
          submit
        </Button>
        <Link to="/password/request" className="forget-password-link"> Forgot your password? </Link>
      </form>
    )
  }
};
