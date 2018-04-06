import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';

export default class PasswordRequest extends Component {
	render() {
  	return (
    	<form style = {{ display: "flex", flexDirection: "column", alignItems: "center" }}>
				<TextField label="Verification code" margin="normal" fullWidth required/>
				<TextField label="New password" margin="normal" fullWidth required/>
        <Button variant="raised" color="primary" type="submit" style={{ padding:"10px 30px", marginTop:"15px" }}>submit</Button>
			</form>
    )
  }
};