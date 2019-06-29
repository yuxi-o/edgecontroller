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
import withStyles from '@material-ui/core/styles/withStyles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import { withRouter } from "react-router-dom";
import NodeCard from '../components/cards/CardItem';
import Topbar from '../components/Topbar';
import ApiClient from '../api/ApiClient';
import CircularLoader from '../components/progressbars/FullSizeCircularLoader';
import Button from '@material-ui/core/Button';
import AddIcon from '@material-ui/icons/Add';
import Typography from "@material-ui/core/Typography";
import AddNodeFormDialog from '../components/forms/AddNodeFormDialog';
import { withSnackbar } from 'notistack';

const styles = theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.grey['A500'],
    overflow: 'hidden',
    backgroundSize: 'cover',
    backgroundPosition: '0 400px',
    marginTop: 20,
    padding: 20,
    paddingBottom: 200
  },
  grid: {
    width: 1000
  },
  rightIcon: {
    marginLeft: '5px',
    marginRight: '0',
    fontSize: '20px',
  },
  sectionContainer: {
    marginTop: theme.spacing.unit * 4,
    marginBottom: theme.spacing.unit * 4
  },
  title: {
    fontWeight: 'bold',
  },
  subtitle: {
    display: 'inline-block',
  },
  addButton: {
    float: 'right',
  },
});

class NodesView extends Component {

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      showAddNodeForm: false,
    };

    this.handleClickOpen = this.handleClickOpen.bind(this);
    this.handleParentClose = this.handleParentClose.bind(this);
  }

  getNodes = () => {
    return ApiClient.get('/nodes');
  };

  handleClickOpen = () => {
    this.setState({ showAddNodeForm: !this.state.showAddNodeForm });
  };

  handleParentClose = () => {
    this.setState({ showAddNodeForm: false });
  };

  handleParentRefresh = () => {
    this.setState({ loaded: false });
    this.fetchNodes();
  };

  componentDidMount = () => {
    return this.fetchNodes();
  };

  fetchNodes = () => {
    return this.getNodes().then((resp) => {
      // Do Something
      if (resp.data) {
        this.setState({ loaded: true, nodes: resp.data.nodes })
      }
    }).catch((err) => {
      this.props.enqueueSnackbar(`Error loading edge nodes. Please try again later.`, {
        variant: 'error',
      });
      this.setState({ loaded: true });
    });
  };

  render() {
    const { location: { pathname: currentPath }, classes } = this.props;

    const renderNodes = () => {
      const { nodes } = this.state || {};
      if (nodes) {
        const nodeDialog = "You are about to delete a Edge Node. In order to re-enroll the deleted node, you may have to re-image it.";

        return Object.keys(nodes).map(key => {
          return (
            <NodeCard
              resourcePath="/nodes"
              key={nodes[key].id}
              CardItem={nodes[key]}
              dialogText={nodeDialog}
              excludeKeys={[]}
            />
          )
        })
      }
    };

    const renderAddNodeForm = () => {
      return (
        <AddNodeFormDialog
          open={this.state.showAddNodeForm}
          handleParentClose={this.handleParentClose}
          handleParentRefresh={this.handleParentRefresh}
        />
      );
    };

    const nodesGrid = () => (
      <Grid spacing={24} alignItems="center" justify="center" container className={classes.grid}>
        <Grid item xs={12}>
          <Grid container direction="row"
            justify="space-between"
            alignItems="flex-start"
            className={classes.sectionContainer}
          >
            <Grid item>
              <Typography variant="subtitle1" className={classes.title}>
                Edge Nodes
                </Typography>
              <Typography variant="body1" gutterBottom className={classes.subtitle}>
                List of Edge Nodes
              </Typography>
            </Grid>
            <Grid item xs={3}>
              <Button variant="contained" color="primary" className={classes.addButton} onClick={this.handleClickOpen}>
                Add Edge Node
                  <AddIcon className={classes.rightIcon} />
              </Button>
            </Grid>
          </Grid>
          {renderNodes()}
        </Grid>
      </Grid>
    );

    const circularLoader = () => (
      <CircularLoader />
    );

    return (
      <React.Fragment>
        <CssBaseline />
        <Topbar currentPath={currentPath} />
        <div className={classes.root}>
          <Grid container justify="center">
            {this.state.loaded ? nodesGrid() : circularLoader()}
            {renderAddNodeForm()}
          </Grid>
        </div>
      </React.Fragment>
    )
  }
}

export default withRouter(withSnackbar(withStyles(styles)(NodesView)));
