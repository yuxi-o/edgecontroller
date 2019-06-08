import React, { Component } from 'react';
import ApiClient from "../../api/ApiClient";

export default class AppView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      errored: false,
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
          app: resp.data || [],
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
          errored: true,
          error: err,
        });
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
          errored: true,
          error: err,
        });
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
          errored: true,
          error: err,
        });
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
          errored: true,
          error: err,
        });
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
          errored: true,
          error: err,
        });
      });
  }

  componentDidMount() {
    this.getNodeApp();
  }

  render() {
    const {
      loaded,
      errored,
      error,
      app,
    } = this.state;

    if (!loaded) {
      return <React.Fragment>Loading ...</React.Fragment>
    }

    if (errored) {
      return <React.Fragment>{error.toString()}</React.Fragment>
    }

    return (
      <React.Fragment>
        {JSON.stringify(app)}
      </React.Fragment>
    );
  }
};
