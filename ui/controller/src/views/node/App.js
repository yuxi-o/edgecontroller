import React, { Component } from 'react';
import ApiClient from "../../api/ApiClient";
import { withSnackbar } from 'notistack';

class AppView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      error: null,
      app: {},
      policy: {},
    };
  }

  // GET /nodes/:node_id/apps/:app_id
  getNodeApp = () => {
    const { nodeID, appID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/apps/${appID}`)
      .then((resp) => {
        this.setState({
          loaded: true,
          app: resp.data.apps || [],
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
  updateNodeApp = () => {
    const { nodeID, appID } = this.props;
    const { app } = this.state;

    ApiClient.patch(`/nodes/${nodeID}/apps/${appID}`, app)
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

  // GET /nodes/:node_id/apps/:app_id/policy
  getNodeAppPolicy = () => {
    const { nodeID, appID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/apps/${appID}/policy`)
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

  // PATCH /nodes/:node_id/apps/:app_id/policy
  updateNodeAppPolicy = () => {
    const { nodeID, appID } = this.props;
    const { policy } = this.state;

    ApiClient.patch(`/nodes/${nodeID}/apps/${appID}/policy`, policy)
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

  // DELETE /nodes/:node_id/apps/:app_id/policy
  deleteNodeAppPolicy = () => {
    const { nodeID, appID } = this.props;

    ApiClient.delete(`/nodes/${nodeID}/apps/${appID}/policy`)
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

  componentDidMount() {
    this.getNodeApp();
  }

  render() {
    const {
      loaded,
      app,
    } = this.state;

    if (!loaded) {
      return <React.Fragment>Loading ...</React.Fragment>
    }

    return (
      <React.Fragment>
        {JSON.stringify(app)}
      </React.Fragment>
    );
  }
};

export default withSnackbar(AppView);
