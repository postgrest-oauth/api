import React, { Component } from 'react';
import { TextField, Button, Typography } from 'material-ui';

export default class Verify extends Component {
  render() {
    return(
      <form style = {{ display: "flex", flexDirection: "column", alignItems: "center" }}>
        <Typography>Please input verification code from email</Typography>
        <TextField label="Verification code" margin="normal" fullWidth required/>
        <Button variant="raised" color="primary" type="submit" style={{ padding:"10px 30px", marginTop:"15px" }}>submit</Button>
      </form>
    )
  }
};