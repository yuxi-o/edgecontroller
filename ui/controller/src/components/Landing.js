/* SPDX-License-Identifier: Apache-2.0
 * Copyright © 2019 Intel Corporation
 */

import React, { Component } from 'react';
import { withStyles } from '@material-ui/core/styles';
import './LandingHeader.css'
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
  Button,
  Typography,
  Paper,
  AppBar,
  Link,
} from '@material-ui/core';

const CONTROLLER_URL = process.env.REACT_APP_CONTROLLER_UI_URL
const CUPS_URL = process.env.REACT_APP_CUPS_UI_URL
const CNCA_URL = process.env.REACT_APP_CNCA_UI_URL

const styles = theme => ({
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

  paper: {
    marginTop: theme.spacing.unit * 3,
    marginBottom: theme.spacing.unit * 3,
    padding: theme.spacing.unit * 2,
    [theme.breakpoints.up(600 + theme.spacing.unit * 3 * 2)]: {
      marginTop: theme.spacing.unit * 6,
      marginBottom: theme.spacing.unit * 6,
      padding: theme.spacing.unit * 3,
    },
  },
});

class Landing extends Component {
  _isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      hasError: false,
    }
  }

  _cancelIfUnmounted = (action) => {
    if (this._isMounted) {
      action();
    }
  }

  componentWillUnmount() {
    // Signal to cancel any pending async requests to prevent setting state
    // on an unmounted component.
    this._isMounted = false;
  }

  async componentDidMount() {
    this._isMounted = true;
  }

  render() {
    const {
      classes,
    } = this.props;

    const {
      hasError,
      error,
    } = this.state;

    if (hasError) {
      throw error;
    }

    const LandingChoiceTableRow = ({ match, history, item }) => {
      return (
        <TableRow>
          <TableCell>
            <Button
              onClick={() => window.location.assign(`${CONTROLLER_URL}/login/`)}
              variant="outlined"
            >
              Infrastructure Manager
            </Button>
          </TableCell>
          <TableCell>
            <Button
              onClick={() => window.location.assign(`${CUPS_URL}/`)}
              variant="outlined"
            >
              LTE CUPS Core Network
            </Button>
          </TableCell>
          <TableCell>
            <Button
              onClick={() => window.location.assign(`${CNCA_URL}/`)}
              variant="outlined"
            >
              5G Next-Gen Core Network
            </Button>
          </TableCell>
        </TableRow>
      );
    }

    return (
      <div>
        <AppBar className={classes.header}>
          <Typography variant="h6" component="h2" id="title">
            <Link to="/">
              Infrastructure and Network Management
            </Link>
          </Typography>
        </AppBar>
              
        <Paper className={classes.paper}>
          <Table>
            <TableBody>
              {
                <LandingChoiceTableRow/>
              }
            </TableBody>
          </Table>

        </Paper>
      </div>
    );
  }
};

export default withStyles(styles)(Landing);

