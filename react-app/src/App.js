import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Login from './components/Login';
import { Card, Grid } from 'material-ui';
import PasswordRequest from './components/PasswordRequest';
import PasswordReset from './components/PasswordReset';
import Verify from './components/Verify';

export default class App extends Component {
  render() {
    return (
      <Grid container style={{paddingTop:"120px", justifyContent:"center"}}>
        <Router basename="/ui">
          <Card raised = {true} style = {{ width: "420px", padding: "20px 50px" }}>
            <Switch>
              <Route exact path="/signin" component={Login}/>
              <Route path="/passwordrequest" component={PasswordRequest}/>
              <Route path="/passwordreset" component={PasswordReset}/>
              <Route path="/verify" component={Verify}/>
            </Switch>
          </Card>
        </Router>
      </Grid>
    );
  }
}

