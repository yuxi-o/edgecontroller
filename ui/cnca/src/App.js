/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

import React, { Component } from 'react';
import { Route, Redirect, Switch } from 'react-router-dom';
import { BrowserRouter as Router } from "react-router-dom";
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Subscriptions from './components/Subscriptions';
import Subscription from './components/Subscription';
import SubscriptionModify from './components/SubscriptionModify';
import Services from './components/Services.js';
import Service from './components/Service.js';
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
                  render={() => <Redirect to="/services" />}
                />
                <Route
                  exact
                  path="/services"
                  component={Services}
                />
                <Route
                  exact
                  path="/services/create"
                  render={(props) => <Service {...props} createMode={true}/>}
                />
                <Route
                  exact
                  path="/services/:id"
                  component={Service}
                />
                <Route
                  exact
                  path="/subscriptions"
                  component={Subscriptions}
                />
                <Route
                  exact
                  path="/subscriptions/create"
                  render={(props) => <Subscription {...props} createMode={true} />}
                />
                <Route
                  exact
                  path="/subscriptions/edit/:id"
                  component={Subscription}
                />
                <Route
                  exact
                  path="/subscriptions/patch/:id"
                  component={SubscriptionModify}
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
