import React, { Component } from 'react';
import { TextField, Button } from 'material-ui';
import { Link } from 'react-router-dom';

export default class PasswordRequest extends Component {
	render() {
  	return (
    	<form style = {{ display: "flex", flexDirection: "column", alignItems: "center" }}>
				<TextField label="Email or phone" margin="normal" fullWidth required/>
				<Button 
					variant="raised" 
					color="primary"
					component={Link}
					to="/passwordreset"
					style={{ padding:"10px 30px", marginTop:"15px" }}
				>
					submit
				</Button>
			</form>
    )
  }
};