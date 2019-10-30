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
import ApiClient from '../../api/ApiClient';

import withStyles from '@material-ui/core/styles/withStyles';
import Grid from '@material-ui/core/Grid';
import Topbar from '../../components/Topbar';
import Table from '../../components/tables/EnhancedTable';
import AddIcon from '@material-ui/icons/Add';
import { withSnackbar } from 'notistack';
import CircularLoader from '../../components/progressbars/FullSizeCircularLoader';
import CssBaseline from '@material-ui/core/CssBaseline';
import OrchestrationContext from '../../context/orchestrationContext';

import { Typography, Button } from '@material-ui/core';

const styles = (theme) => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.grey['A500'],
    overflow: 'hidden',
    backgroundSize: 'cover',
    backgroundPosition: '0 400px',
    marginTop: 20,
    padding: 20,
    paddingBottom: 200,
  },
  grid: {
    width: 1000,
  },
  addButton: {
    float: 'right',
  },
});

class PoliciesView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      policies: [],
    };
  }

  static contextType = OrchestrationContext;

  componentDidMount = () => {
    return this.getPolicies();
  };

  // GET /policies
  getPolicies = () => {
    ApiClient.get(`${this.context.apiClientPath}/policies`)
      .then((resp) => {
        this.setState({
          loaded: true,
          policies: resp.data.policies || [],
        });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(
          `${err.toString()}. Please try again later.`,
          { variant: 'error' }
        );
      });
  };

  handleOnClickAddPolicy = () => {
    const { history } = this.props;

    // Redirect user to add policies view
    history.push('/policies/add');
  };

  renderTable = () => {
    const { policies } = this.state;
    const data = [];

    const tableHeaders = [
      { id: 'id', numeric: true, disablePadding: false, label: 'ID' },
      { id: 'name', numeric: false, disablePadding: false, label: 'Name' },
      { id: 'action', numeric: false, disablePadding: false, label: 'Action' },
    ];

    // Append the edit view url to each policy
    policies.forEach((item) =>
      data.push({ ...item, editUrl: `/policies/${item.id}/edit` })
    );

    const tableData = {
      order: 'asc',
      orderBy: 'id',
      selected: [],
      data: data,
      page: 0,
      rowsPerPage: 10,
    };

    return <Table rows={tableHeaders} tableState={tableData} />;
  };

  render() {
    const { classes } = this.props;
    const currentPath = this.props.location.pathname;

    const circularLoader = () => <CircularLoader />;

    const policiesGrid = () => (
      <Grid
        spacing={24}
        alignItems="center"
        justify="center"
        container
        className={classes.grid}
      >
        <Grid item xs={12}>
          <Grid
            container
            direction="row"
            justify="space-between"
            alignItems="flex-start"
            className={classes.sectionContainer}
          >
            <Grid item>
              <Typography variant="subtitle1" className={classes.title}>
                Traffic Policies
              </Typography>
              <Typography
                variant="body1"
                gutterBottom
                className={classes.subtitle}
              >
                List of Traffic Policies
              </Typography>
            </Grid>
            <Grid item xs={3}>
              <Button
                variant="contained"
                color="primary"
                className={classes.addButton}
                onClick={this.handleOnClickAddPolicy}
              >
                Add Policy
                <AddIcon className={classes.rightIcon} />
              </Button>
            </Grid>
          </Grid>
          {this.renderTable()}
        </Grid>
      </Grid>
    );

    return (
      <React.Fragment>
        <CssBaseline />
        <Topbar currentPath={currentPath} />
        <div className={classes.root}>
          <Grid container justify="center">
            {this.state.loaded ? policiesGrid() : circularLoader()}
          </Grid>
        </div>
      </React.Fragment>
    );
  }
}

export default withSnackbar(withStyles(styles)(PoliciesView));
