import React, { Component } from 'react';
import ApiClient from "../../api/ApiClient";

export default class InterfaceView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      error: null,
      interface: {},
    };
  }

  // GET /nodes/:node_id/interfaces/:interface_id
  getNodeInterface = () => {
    const { nodeID, interfaceID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/interfaces/${interfaceID}`)
      .then((resp) => {
        this.setState({
          loaded: true,
          interface: resp.data || {},
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
          error: err,
        });
      });
  }

  // PATCH /nodes/:node_id/interfaces
  updateNodeInterface = (interface) => {
    const { nodeID } = this.props;

    // TODO: Get from user form
    const data = {
      interfaces: [
        interface,
      ],
    };

    ApiClient.patch(`/nodes/${nodeID}/interfaces`, data)
      .then((resp) => {
        this.setState({
          loaded: true,
        })
      })
      .catch((err) => {
        this.setState({
          loaded: true,
          error: err,
        });
      });
  }

  componentDidMount() {
    this.getNodeInterface();
  }

  render() {
    const {
      loaded,
      error,
      interface,
    } = this.state;

    if (!loaded) {
      return <React.Fragment>Loading ...</React.Fragment>
    }

    return (
      <React.Fragment>
        {JSON.stringify(interface)}
      </React.Fragment>
    );
  }
};
