import React, { Component } from 'react';
import ApiClient from "../../api/ApiClient";
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import DialogContentText from '@material-ui/core/DialogContentText';
import { SchemaForm, utils } from 'react-schema-form';
import InterfaceSchema from '../../components/schema/NodeInterface';
import Select from '@material-ui/core/Select';
import { withSnackbar } from 'notistack';
import {
  Grid,
  Button
} from '@material-ui/core';

class InterfacesView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      error: null,
      showErrors: true,
      interfaces: [],
      open: false,
      interfaceID: '',
      nodeInterface: {},
      openPolicyDialog: false,
      policies: [],
      selectedTrafficPolicyID: '',
      selectedInterfaceID: ''
    };
  }

  getNodeInterfaces = () => {
    const { nodeID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/interfaces`)
      .then((resp) => {
        this.setState({
          loaded: true,
          interfaces: resp.data.interfaces || [],
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // PATCH /nodes/:node_id/interfaces/:interface_id
  updateNodeInterface = () => {
    const { nodeID } = this.props;
    const { nodeInterface } = this.state;

    ApiClient.patch(`/nodes/${nodeID}/interfaces/${nodeInterface.id}`, nodeInterface)
      .then((resp) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`Successfully updated node interface ${nodeInterface.id}.`, { variant: 'success' });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // GET /nodes/:node_id/interfaces/:interface_id
  getNodeInterface = (interfaceID) => {
    const { nodeID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/interfaces/${interfaceID}`)
      .then((resp) => {
        this.setState({
          loaded: true,
          nodeInterface: resp.data || {},
        });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // GET /nodes/:node_id/interfaces/:interface_id/policy
  getNodeInterfacePolicy = () => {
    const { nodeID, interfaceID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/interfaces/${interfaceID}/policy`)
      .then((resp) => {
        this.setState({
          loaded: true,
          policy: resp.data || {},
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // DELETE /nodes/:node_id/interfaces/:interface_id/policy
  deleteNodeInterfacePolicy = () => {
    const { nodeID, interfaceID } = this.props;

    ApiClient.delete(`/nodes/${nodeID}/interfaces/${interfaceID}/policy`)
      .then((resp) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`Successfully deleted policy on interface ${interfaceID}.`, { variant: 'success' });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  onModelChange = (key, val) => {
    const { nodeInterface } = this.state;

    const newInterface = nodeInterface;

    utils.selectOrSet(key, newInterface, val);

    this.setState({ policy: newInterface });
  }

  handleOpen = (interfaceID) => {
    this.getNodeInterface(interfaceID);

    this.setState({
      open: true,
      interfaceID: interfaceID,
    });
  };

  handleClose = () => {
    this.setState({ open: false });
  };

  handleUpdateInterface = () => {
    this.updateNodeInterface();
  };

  renderEditInterfaceDialog = () => {
    const {
      showErrors,
      nodeInterface,
    } = this.state

    return (
      <Dialog
        open={this.state.open}
        onClose={this.handleClose}
      >
        <DialogTitle id="form-dialog-title">Edit Interface</DialogTitle>
        <DialogContent>
          <SchemaForm
            schema={InterfaceSchema.schema}
            form={InterfaceSchema.form}
            model={nodeInterface}
            onModelChange={this.onModelChange}
            showErrors={showErrors}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={this.handleClose} color="primary">
            Cancel
            </Button>
          <Button onClick={this.handleUpdateInterface} color="primary">
            Save
            </Button>
        </DialogActions>
      </Dialog>
    );
  }

  renderTable = () => {
    const { interfaces } = this.state;

    return (
      <React.Fragment>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell align="right">Description</TableCell>
              <TableCell align="right">Driver</TableCell>
              <TableCell align="right">Type</TableCell>
              <TableCell align="right">Traffic Policy</TableCell>
              <TableCell align="right">Action</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {interfaces.map(row => (
              <TableRow key={row.id}>
                <TableCell component="th" scope="row">
                  {row.id}
                </TableCell>
                <TableCell align="right">{row.description}</TableCell>
                <TableCell align="right">{row.driver}</TableCell>
                <TableCell align="right">{row.type}</TableCell>
                <TableCell align="right">
                  <Button onClick={() => this.handleOpenPolicy(row.id)}>
                    Add
                  </Button>
                </TableCell>
                <TableCell align="right">
                  <Button onClick={() => this.handleOpen(row.id)}>
                    Edit
                </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        {this.renderEditInterfaceDialog()}
      </React.Fragment>
    );
  }

  // GET /policies
  getTrafficPolicies = () => {
    ApiClient.get(`/policies`)
      .then((resp) => {
        this.setState({
          loaded: true,
          policies: resp.data || [],
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
          error: err,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  handleAssignInterfaceTrafficPolicy = () => {
    const { nodeID } = this.props;
    const { selectedInterfaceID, selectedTrafficPolicyID } = this.state;

    const data = {
      id: selectedTrafficPolicyID,
    }

    ApiClient.patch(`/nodes/${nodeID}/apps/${selectedInterfaceID}`, data)
      .then((resp) => {
        this.setState({
          loaded: true,
          policies: resp.data || [],
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
          error: err,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });

    this.handleClosePolicy();
  }

  handleOpenPolicy = (interfaceID) => {
    this.setState({
      selectedInterfaceID: interfaceID,
      openPolicyDialog: true,
    });
  }

  handleClosePolicy = () => {
    this.setState({ openPolicyDialog: false });
  }

  handleChangeTrafficPolicy = event => {
    this.setState({ selectedTrafficPolicyID: event.target.value });
  }

  renderTrafficPolicyDialog = () => {
    const { policies } = this.state;

    return (
      <Dialog
        open={this.state.openPolicyDialog}
        onClose={this.handleClosePolicy}
      >
        <DialogTitle id="form-dialog-title">Assign Traffic Policy to Interface</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Choose a traffic policy to assign:
            </DialogContentText>
          <form>
            <FormControl>
              <Select
                native
                onChange={this.handleChangeTrafficPolicy}
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
          <Button onClick={() => this.handleAssignInterfaceTrafficPolicy()} color="primary">
            Assign
                </Button>
        </DialogActions>
      </Dialog>
    );
  }

  componentDidMount() {
    this.getNodeInterfaces();
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
          <Grid item xs={12}>
            {this.renderTable()}
          </Grid>
        </Grid>

        {this.renderTrafficPolicyDialog()}
      </React.Fragment>
    );
  }
};

export default withSnackbar(InterfacesView);
