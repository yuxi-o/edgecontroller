// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React, { Component } from 'react';
import { Route, Redirect, Switch } from 'react-router-dom';
import { BrowserRouter as Router } from "react-router-dom";
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Userplanes from './components/Userplanes';
import Userplane from './components/Userplane';
import Header from './components/Header';
import ErrorBoundary from './components/ErrorBoundary';

const useStyles = theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    minHeight: '100vh',
  },
  header: {
    padding: 20,
    flexGrow: 1,
  },
  main: {
    width: 'auto',
    marginTop: theme.spacing.unit * 14,
    marginBottom: theme.spacing.unit * 2,
    marginLeft: theme.spacing.unit * 2,
    marginRight: theme.spacing.unit * 2,
    [theme.breakpoints.up(600 + theme.spacing.unit * 4)]: {
      width: 800,
      marginLeft: 'auto',
      marginRight: 'auto',
    },
  },
});

class App extends Component {
  render() {
    const { classes } = this.props;

    return (
      <Router>
        <div>
          <header>
            <Header className={classes.header} />
          </header>
          <main className={classes.main}>
            <ErrorBoundary>
              <Switch>
                <Route
                  exact
                  path="/"
                  render={() => <Redirect to="/userplanes" />}
                />
                <Route
                  exact
                  path="/userplanes"
                  component={Userplanes}
                />
                <Route
                  exact
                  path="/userplanes/create"
                  render={(props) => <Userplane {...props} createMode={true} />}
                />
                <Route
                  exact
                  path="/userplanes/:id"
                  component={Userplane}
                />
                <Route
                  render={() => <span>404 Not Found</span>}
                />
              </Switch>
            </ErrorBoundary>
          </main>
        </div>
      </Router>
    )
  }
}

App.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(useStyles)(App);
