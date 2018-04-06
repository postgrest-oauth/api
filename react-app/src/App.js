import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import Login from './components/Login';
import { Card, Grid } from 'material-ui';
import PasswordRequest from './components/PasswordRequest';
import PasswordReset from './components/PasswordReset';
import Verify from './components/Verify';

class App extends Component {
  render() {
    return (
      <Grid container style={{paddingTop:"120px", justifyContent:"center"}}>
        <Router>
          <Card raised = {true} style = {{ width: "420px", padding: "20px 50px" }}>
            <Route exact path="/" component={Login}/>
            <Route path="/passwordrequest" component={PasswordRequest}/>
            <Route path="/passwordreset" component={PasswordReset}/>
            <Route path="/verify" component={Verify}/>
          </Card>
        </Router>
      </Grid>
    );
  }
}

export default App;
