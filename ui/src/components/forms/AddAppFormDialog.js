import React,  { Component } from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import ApiClient from '../../api/ApiClient';
import CircularLoader from '../progressbars/FullSizeCircularLoader';
import Button from '@material-ui/core/Button/index';
import TextField from '@material-ui/core/TextField/index';
import Dialog from '@material-ui/core/Dialog/index';
import DialogActions from '@material-ui/core/DialogActions/index';
import DialogContent from '@material-ui/core/DialogContent/index';
import DialogContentText from '@material-ui/core/DialogContentText/index';
import DialogTitle from '@material-ui/core/DialogTitle/index';
import FormHelperText from "@material-ui/core/FormHelperText";
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';

const styles = theme => ({
  circularLoaderContainer: {
    position: 'absolute',
    width: '100%',
    height: '100%',
  },
  protocolSelect: {
    minWidth: '125px',
  },
});

class AddAppFormDialog extends Component {

  constructor(props) {
    super(props);

    const {open, handleParentClose, handleParentRefresh} = this.props;

    this.state = {
      open: open,
      loading: false,
      helperText: null,
      ports: [{port: 0, protocol: ''}],
    };

    this.handleDialogClose = this.handleDialogClose.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handlePortInput = this.handlePortInput.bind(this);
    this.handleProtocol = this.handleProtocol.bind(this);
    this.handleParentRefresh = handleParentRefresh.bind(this);
    this.handleParentClose = handleParentClose.bind(this);
  }

  static getDerivedStateFromProps(nextProps, prevState){
    if(nextProps.open !== prevState.open){
      return { open: nextProps.open};
    }

    return null;
  }

  handlePortInput = (e) => {
    let ports = [...this.state.ports];
    ports[e.target.dataset.id]['port'] = Math.trunc(e.target.value);
    this.setState({ports});
  };

  handleProtocol = (e) => {
    const [name, id] = e.target.name.split('-');
    let ports = [...this.state.ports];
    ports[id][name] = e.target.value;
    this.setState({ports});
  };

  handleInputChange = (e) => {
    if(e.target.name === 'cores' || e.target.name === 'memory') {
      this.setState({[e.target.name]: Math.trunc(e.target.value)});
      return ;
    }
    this.setState({[e.target.name]: e.target.value});
  };

  addPorts = (e) => {
    this.setState((prevState) => ({
      ports: [...prevState.ports, {port: 0, protocol: ''}]
    }));
  };

  handleDialogClose = () => {
    this.handleParentClose(!this.state.open);
  };

  clearPortList = () => {
    this.setState({ports: [{port: 0, protocol: ''}]});
  };

  handleSubmit = (e) => {
    e.preventDefault();

    const getAppFormValues = () => {
      const {name, version, type, vendor, description, cores, memory, ports, source} = this.state;

      return {name, version, type, vendor, description, cores, memory, ports, source};
    };

    if (this.state.loading === true) {
      return ;
    }

    this.setState({loading: true});

    return ApiClient.post('/apps', getAppFormValues())
      .then((resp) => {
        this.setState({loading: false});
        this.handleDialogClose();
        this.handleParentRefresh();
      })
      .catch((err) => {
        console.log(err);

        if("response" in err && "data" in err.response) {
          return this.setState({loading: false, submitError: true, helperText: err.response.data});
        }

        return this.setState({loading: false, submitError: true, helperText: 'Unknown error try again later'});
      });
  };

  render() {
    const {classes} = this.props;
    const circularLoader = () => (
      <div className={classes.circularLoaderContainer}>
        <CircularLoader />
      </div>
    );

    const dialogActions = () => (
      <DialogActions>
        <Button onClick={this.handleDialogClose} color="primary">
          Cancel
        </Button>
        <Button onClick={this.handleSubmit} type="submit" variant="contained" color="primary">
         Upload Application
        </Button>
      </DialogActions>
    );

    const generateInputField = (id, name, label) => (
      <TextField
        autoFocus
        margin="dense"
        onChange={this.handleInputChange}
        id={id}
        name={name}
        label={label}
        type={(name === 'cores' || name === 'memory') ? "number" : "text"}
        fullWidth
      />
    );

    const generateSelectField = (id, name, values) => {
      return (
        <React.Fragment>
          <InputLabel htmlFor={name}>{name}</InputLabel>
          <Select
            value={this.state[id] || ''}
            onChange={this.handleInputChange}
            name={name}
            id={id}
          >
            {values.map((val, idx) => {
              return (
                <MenuItem key={idx} value={val}>{val}</MenuItem>
              );
            })}
          </Select>
        </React.Fragment>
      )
    };

    const renderPortsInputField = () => {
      const button = (<button onClick={this.addPorts}>Add Additional Port</button>);

      return (
        <React.Fragment>
          {
            this.state.ports.map((val, idx) => {
              let portId = `port-${idx}`, protocolId = `protocol-${idx}`;
              return (
                <div className={classes.ports} key={idx}>
                  <TextField
                    autoFocus
                    margin="dense"
                    id={portId}
                    name={portId}
                    label="Port"
                    type="number"
                    inputProps={{
                      'data-id': idx,
                    }}
                    onChange={this.handlePortInput}
                  />
                  <FormControl className={classes.formControl}>
                    <InputLabel htmlFor={protocolId}>Protocol</InputLabel>
                    <Select
                      value={this.state.ports[idx].protocol}
                      onChange={this.handleProtocol}
                      name={protocolId}
                      data-id={idx}
                      id={protocolId}
                      className={classes.protocolSelect}
                    >
                      <MenuItem value="tcp">TCP</MenuItem>
                      <MenuItem value="udp">UDP</MenuItem>
                      <MenuItem value="sctp">SCTP</MenuItem>
                    </Select>
                  </FormControl>
                  {button}
                </div>
              )
            })
          }
        </React.Fragment>
      );
    };

    return (
      <React.Fragment>
        <Dialog
          open={this.state.open}
          onClose={this.handleDialogClose}
          onEnter={this.clearPortList}
          aria-labelledby="add-node-dialog-title"
          aria-describedby="add-node-dialog-description"
        >
          <DialogTitle id="add-node-dialog-title">Add an Application</DialogTitle>
          <DialogContent>
            <DialogContentText id="add-node-dialog-description">
            </DialogContentText>
            <form onSubmit={(e) => {e.preventDefault()}}>
              {generateInputField("name", "name", "Name")}
              {generateSelectField("type", "type", ['container', 'vm'])}
              {generateInputField("version", "version", "Version")}
              {generateInputField("vendor", "vendor", "Vendor")}
              {generateInputField("description", "description", "Description")}
              {generateInputField("cores", "cores", "Cores")}
              {generateInputField("memory", "memory", "Memory")}
              {generateInputField("source", "source", "Source")}
              {renderPortsInputField()}

              { this.state.helperText !== "" ?
                <FormHelperText id="component-error-text">
                  {this.state.helperText}
                </FormHelperText> : null
              }
            </form>
          </DialogContent>
          {dialogActions()}
          {(this.state.loading)? circularLoader() : null}
        </Dialog>
      </React.Fragment>
    )
  }
}

export default withStyles(styles)(AddAppFormDialog);
