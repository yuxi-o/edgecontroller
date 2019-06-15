import React, { Component } from 'react';
import ApiClient from "../../api/ApiClient";
import AppView from './App';
import { withSnackbar } from 'notistack';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import {
  Grid,
  Button,
} from '@material-ui/core';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import {
  Add
} from '@material-ui/icons';

class AppsView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      error: null,
      open: false,
      nodeApps: [],
      nodeAppsStatus: {},
      apps: [],
      selectedAppID: '',
      openPolicyDialog: false,
      policies: [],
      selectedTrafficPolicyID: '',
    };
  }

  // GET /nodes/:node_id/apps
  getNodeApps = () => {
    const { nodeID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/apps`)
      .then((resp) => {
        this.setState({
          loaded: true,
          nodeApps: resp.data.apps || []
        })

        this.getNodeAppsStatus();
      })
      .catch((err) => {
        this.setState({
          loaded: true,
          nodeApps: [],
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // GET /nodes/:node_id/apps/:app_id
  getNodeAppsStatus = async () => {
    const { nodeID } = this.props;
    const { nodeApps, nodeAppsStatus } = this.state;

    nodeApps.forEach(app => {
      ApiClient.get(`/nodes/${nodeID}/apps/${app.id}`)
        .then((resp) => {
          this.setState({
            nodeAppsStatus: {
              ...nodeAppsStatus,
              [app.id]: resp.data.status || 'unknown',
            },
          })
        })
        .catch((err) => {
          this.setState({
            error: err,
            nodeAppsStatus: {
              ...nodeAppsStatus,
              [app.id]: 'unknown',
            },
          });

          this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
        });
    })
  }

  // GET /apps
  getApps = () => {
    ApiClient.get(`/apps`)
      .then((resp) => {
        this.setState({
          loaded: true,
          apps: resp.data.apps || [],
        });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
          error: err,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // POST /nodes/:node_id/apps
  deployNodeApp = () => {
    const { nodeID } = this.props;
    const { selectedAppID } = this.state;

    const data = {
      app_id: selectedAppID,
    };

    ApiClient.post(`/nodes/${nodeID}/apps`, data)
      .then((resp) => {
        this.setState({
          loaded: true,
        })

        this.props.enqueueSnackbar(`Successfully deployed app ${selectedAppID}.`, { variant: 'success' });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });
        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // DELETE /nodes/:node_id/apps/:app_id
  deleteNodeApp = (appID) => {
    const { nodeID } = this.props;

    ApiClient.delete(`/nodes/${nodeID}/apps/${appID}`)
      .then((resp) => {
        this.setState({
          loaded: true,
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // PATCH /nodes/:node_id/apps/:app_id
  commandNodeApp = (appID, command) => {
    const { nodeID } = this.props;

    const data = {
      command: command,
    }

    ApiClient.patch(`/nodes/${nodeID}/apps/${appID}`, data)
      .then((resp) => {
        this.setState({
          loaded: true,
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // GET /policies
  getTrafficPolicies = () => {
    ApiClient.get(`/policies`)
      .then((resp) => {
        this.setState({
          loaded: true,
          policies: resp.data.policies || [],
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  handleAssignAppTrafficPolicy = () => {
    const { nodeID } = this.props;
    const { selectedAppID, selectedTrafficPolicyID } = this.state;

    const data = {
      id: selectedTrafficPolicyID,
    }

    ApiClient.patch(`/nodes/${nodeID}/apps/${selectedAppID}`, data)
      .then((resp) => {
        this.setState({
          loaded: true,
          policies: resp.data || [],
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });

    this.handleClosePolicy();
  }

  handleOpen = () => {
    this.setState({ open: true });
  };

  handleClose = () => {
    this.setState({ open: false });
  };

  handleOpenPolicy = (appID) => {
    this.setState({
      selectedAppID: appID,
      openPolicyDialog: true,
    });
  }

  handleClosePolicy = () => {
    this.setState({ openPolicyDialog: false });
  }

  handleDeployApp = () => {
    this.deployNodeApp();

    this.handleClose();
  }

  handleChange = event => {
    this.setState({ selectedAppID: event.target.value });
  };

  renderNodeApp = (appID) =>
    <AppView
      nodeID={this.props.nodeID}
      appID={appID}
    />

  renderDeployDialog = () => {
    const { apps } = this.state;

    return (
      <Dialog
        open={this.state.open}
        onClose={this.handleClose}
      >
        <DialogTitle id="form-dialog-title">Deploy Application to Node</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Choose an application to deploy:
        </DialogContentText>
          <form>
            <FormControl>
              <Select
                native
                value={this.state.selectedAppID}
                onChange={this.handleChange}
              >
                {
                  apps.map(item => <option key={item.id} value={item.id}>{item.id}</option>)
                }
              </Select>
            </FormControl>
          </form>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.handleClose} color="primary">
            Cancel
            </Button>
          <Button onClick={this.handleDeployApp} color="primary">
            Deploy
            </Button>
        </DialogActions>
      </Dialog>
    );
  }

  renderTrafficPolicyDialog = () => {
    const { policies, selectedAppID } = this.state;

    return (
      <Dialog
        open={this.state.openPolicyDialog}
        onClose={this.handleClosePolicy}
      >
        <DialogTitle id="form-dialog-title">Assign Traffic Policy to Node</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Choose a traffic policy to assign:
        </DialogContentText>
          <form>
            <FormControl>
              <Select
                native
                value={this.state.selectedTrafficPolicyID}
                onChange={this.handleChange}
              >
                {
                  policies.map(item => <option value={item.id}>{item.id}</option>)
                }
              </Select>
            </FormControl>
          </form>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.handleClosePolicy} color="primary">
            Cancel
            </Button>
          <Button onClick={() => this.handleAssignAppTrafficPolicy(selectedAppID)} color="primary">
            Assign
            </Button>
        </DialogActions>
      </Dialog>
    );
  }

  renderTable = () => {
    const {
      nodeApps,
      nodeAppsStatus,
    } = this.state;

    return (
      <React.Fragment>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell align="right">Status</TableCell>
              <TableCell align="right">Type</TableCell>
              <TableCell align="right">Name</TableCell>
              <TableCell align="right">Traffic Policy</TableCell>
              <TableCell align="right">Lifecycle</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {
              nodeApps.map(row => {
                return (
                  <TableRow key={row.id} >
                    <TableCell component="th" scope="row">
                      {row.id}
                    </TableCell>
                    <TableCell align="right">
                      {nodeAppsStatus[row.id]}
                    </TableCell>
                    <TableCell align="right">{row.type}</TableCell>
                    <TableCell align="right">{row.name}</TableCell>
                    <TableCell align="right">
                      <Button onClick={() => this.handleOpenPolicy(row.id)}>
                        Add
                      </Button>
                    </TableCell>
                    <TableCell align="right">
                      <Button onClick={() => this.commandNodeApp(row.id, 'start')}>
                        Start
                      </Button>
                      <Button onClick={() => this.commandNodeApp(row.id, 'stop')}>
                        Stop
                      </Button>
                      <Button onClick={() => this.commandNodeApp(row.id, 'restart')}>
                        Restart
                      </Button>
                    </TableCell>
                    <TableCell align="right">
                      <Button onClick={() => this.deleteNodeApp(row.id)}>
                        Delete
                      </Button>
                    </TableCell>
                  </TableRow>
                );
              })
            }
          </TableBody>
        </Table>
      </React.Fragment>
    );
  }

  componentDidMount() {
    this.getApps();
    this.getNodeApps();
    this.getTrafficPolicies();
  }

  render() {
    const {
      loaded
    } = this.state;

    if (!loaded) {
      return <React.Fragment>Loading ...</React.Fragment>
    }

    return (
      <React.Fragment>
        <Grid container>
          <Grid item xs={12} style={{ textAlign: 'right' }}>
            <Button
              onClick={this.handleOpen}
              variant="outlined"
              color="primary"
            >
              Deploy App
              <Add />
            </Button>
          </Grid>
          <Grid item xs={12}>
            {this.renderTable()}
          </Grid>
        </Grid>

        {this.renderDeployDialog()}
        {this.renderTrafficPolicyDialog()}
      </React.Fragment>
    );
  }
};

export default withSnackbar(AppsView);
