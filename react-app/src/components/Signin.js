import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';
import { Link } from 'react-router-dom';

export default class Signin extends Component {

  render() {
    return (
      <form style = {{ display: "flex", flexDirection: "column", alignItems: "center" }}>
        <TextField label="Username" margin="normal" fullWidth required/>
        <TextField label="Password" margin="normal" type="password" fullWidth required/>
        <Button variant="raised" color="primary" type="submit" style={{ padding:"10px 30px", margin:"15px 0 10px 0" }}>submit</Button>
        <Link to="/passwordrequest" className="forget-password-link"> Forgot your password? </Link>
      </form>
    )
  }
};
