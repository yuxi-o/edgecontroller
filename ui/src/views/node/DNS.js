import React, { Component } from 'react';
import ApiClient from "../../api/ApiClient";
import { SchemaForm, utils } from 'react-schema-form';
import DNSSchema from '../../components/schema/NodeDnsConfigApply';
import {
  Grid,
  Button
} from '@material-ui/core';

export default class DNSView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      errored: false,
      error: null,
      dns: [],
    };
  }

  // GET /nodes/:node_id/dns
  getNodeDNSConfig = () => {
    const { nodeID } = this.props;

    ApiClient.get(`/nodes/${nodeID}/dns`)
      .then((resp) => {
        this.setState({
          loaded: true,
          dns: resp.data || {},
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

  // PATCH /nodes/:node_id/dns
  applyDNSConfig = () => {
    const { nodeID } = this.props;
    const { dns } = this.state;

    ApiClient.patch(`/nodes/${nodeID}/dns`, dns)
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

  // DELETE /nodes/:node_id/dns/:dns_id
  deleteDNSConfig = (dnsID) => {
    const { nodeID } = this.props;

    ApiClient.delete(`/nodes/${nodeID}/dns/${dnsID}`)
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

  onModelChange = (key, val) => {
    const { dns } = this.state;

    const newDNS = dns;

    utils.selectOrSet(key, newDNS, val);

    this.setState({ dns: newDNS });
  }

  componentDidMount() {
    this.getNodeDNSConfig();
  }

  render() {
    const {
      loaded,
      // errored,
      // error,
      showErrors,
      dns,
    } = this.state;

    if (!loaded) {
      return <React.Fragment>Loading ...</React.Fragment>
    }

    return (
      <React.Fragment>
        <Grid
          container
          justify="center"
          alignItems="flex-end"
          spacing={24}
        >
          <Grid item xs={12}>
            <SchemaForm
              schema={DNSSchema.schema}
              form={DNSSchema.form}
              model={dns}
              onModelChange={this.onModelChange}
              showErrors={showErrors}
            />
          </Grid>
          <Grid item xs={12}>
            <Button
              onClick={this.applyDNSConfig}
              variant="outlined"
              color="primary"
            >
              Save
            </Button>
          </Grid>
        </Grid>
      </React.Fragment>
    );
  }
};
