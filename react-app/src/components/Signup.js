import React, { Component } from 'react';
import { TextField, Button, FormControl, Input, InputLabel } from 'material-ui';
import MaskedInput from 'react-text-mask';
import { Link } from 'react-router-dom'; 

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

  render() {
    return (
      <form style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
        <TextField label="Email address" margin="normal" type="email" fullWidth required/>
        <TextField label="Password" margin="normal" type="password" fullWidth required/>
        <FormControl margin="normal" fullWidth required>
          <InputLabel shrink={true}> Phone number </InputLabel>
          <Input inputComponent={InputMask} />
        </FormControl>
        <Button 
          variant="raised" 
          color="primary" 
          style={{ padding:"10px 30px", marginTop:"15px" }}
          component={Link}
          to="/verify"
        >
          submit
        </Button>
      </form>
    )
  }
};
