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
import { SchemaForm, utils } from 'react-schema-form';
import PolicySchema from '../../components/schema/TrafficPolicy';
import PolicySchemaKOVN from '../../components/schema/TrafficPolicyKOVN';
import Topbar from '../../components/Topbar';
import withStyles from '@material-ui/core/styles/withStyles';
import { withSnackbar } from 'notistack';
import { Grid, Button } from '@material-ui/core';
import OrchestrationContext, {
  orchestrationModes,
} from '../../context/orchestrationContext';

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
    paddingLeft: '20%',
    paddingRight: '20%',
  },
  gridSaveButton: {
    textAlign: 'right',
  },
});

class PolicyView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      error: null,
      showErrors: true,
      policy: {},
    };
  }
  //Current orchestration mode context
  static contextType = OrchestrationContext;

  // GET /policies/:policy_id
  getPolicy = () => {
    const { match } = this.props;

    const policyID = match.params.id;

    if (!policyID) {
      this.setState({
        loaded: true,
      });
      return;
    }

    ApiClient.get(`${this.context.apiClientPath}/${policyID}`)
      .then((resp) => {
        this.setState({
          loaded: true,
          policy: resp.data || {},
        });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}.`, { variant: 'error' });
      });
  };

  // PATCH /policies/:policy_id
  updatePolicy = () => {
    const { match } = this.props;
    const { policy } = this.state;

    const policyID = match.params.id;

    ApiClient.patch(`${this.context.apiClientPath}/${policyID}`, policy)
      .then((resp) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`Successfully updated policy.`, {
          variant: 'success',
        });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}.`, { variant: 'error' });
      });
  };

  // POST /policies
  createPolicy = () => {
    const { history } = this.props;
    const { policy } = this.state;

    ApiClient.post(`${this.context.apiClientPath}`, policy)
      .then((resp) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`Successfully created policy.`, {
          variant: 'success',
        });
        // Delay the redirect so the user has a moment to breath
        setTimeout(() => {
          history.push('/policies');
        }, 250);
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}.`, { variant: 'error' });
      });
  };

  // DELETE /policies/:policy_id
  deletePolicy = () => {
    const { history, match } = this.props;
    const policyID = match.params.id;

    ApiClient.delete(`${this.context.apiClientPath}/${policyID}`)
      .then((resp) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`Deleted policy ${policyID}.`, {
          variant: 'success',
        });
        history.push('/policies');
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}.`, { variant: 'error' });
      });
  };

  onModelChange = (key, val) => {
    const { policy } = this.state;

    const newPolicy = policy;

    utils.selectOrSet(key, newPolicy, val);

    this.setState({ policy: newPolicy });
  };

  componentDidMount() {
    this.getPolicy();
  }

  render() {
    const {
      match,
      location: { pathname: currentPath },
      classes,
    } = this.props;

    const { loaded, showErrors, policy } = this.state;

    if (!loaded) {
      return <React.Fragment>Loading ...</React.Fragment>;
    }

    const currentPolicySchema =
      this.context.mode === orchestrationModes.kubernetes_ovn.name
        ? PolicySchemaKOVN
        : PolicySchema;

    return (
      <React.Fragment>
        <Topbar currentPath={currentPath} />
        <div className={classes.root}>
          <Grid
            container
            justify="center"
            alignItems="flex-end"
            spacing={24}
            className={classes.grid}
          >
            <Grid item xs={12}>
              <SchemaForm
                schema={currentPolicySchema.schema}
                form={currentPolicySchema.form}
                model={policy}
                onModelChange={this.onModelChange}
                showErrors={showErrors}
              />
            </Grid>
            {match.params.id ? (
              <Grid item xs={12} className={classes.gridSaveButton}>
                <Button onClick={this.deletePolicy}>Delete</Button>
              </Grid>
            ) : null}

            <Grid item xs={12} className={classes.gridSaveButton}>
              {match.params.id ? (
                <Button
                  onClick={this.updatePolicy}
                  variant="outlined"
                  color="primary"
                >
                  Save
                </Button>
              ) : (
                <Button
                  onClick={this.createPolicy}
                  variant="outlined"
                  color="primary"
                >
                  Create
                </Button>
              )}
            </Grid>
          </Grid>
        </div>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withSnackbar(PolicyView));
