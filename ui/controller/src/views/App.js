import React, { Component } from 'react';
import ApiClient from '../api/ApiClient';
import { SchemaForm, utils } from 'react-schema-form';
import AppSchema from '../components/schema/App';
import Topbar from '../components/Topbar';
import withStyles from '@material-ui/core/styles/withStyles';
import { withSnackbar } from 'notistack';
import {
  Grid,
  Button
} from '@material-ui/core';

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
    paddingLeft: '20%',
    paddingRight: '20%'
  },
  gridSaveButton: {
    textAlign: 'right',
  }
});

class AppView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      error: null,
      showErrors: true,
      app: {},
    };
  }

  // GET /apps/:app_id
  getApp = () => {
    const { match } = this.props;

    const appID = match.params.id;

    ApiClient.get(`/apps/${appID}`)
      .then((resp) => {
        this.setState({
          loaded: true,
          app: resp.data || {},
        });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  // PATCH /apps/:app_id
  updateApp = () => {
    const { match } = this.props;
    const { app } = this.state;

    const appID = match.params.id;

    ApiClient.patch(`/apps/${appID}`, app)
      .then((resp) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`Successfully updated application.`, { variant: 'success' });
      })
      .catch((err) => {
        this.setState({
          loaded: true,
        });

        this.props.enqueueSnackbar(`${err.toString()}. Please try again later.`, { variant: 'error' });
      });
  }

  onModelChange = (key, val) => {
    const { app } = this.state;

    const newApp = app;

    utils.selectOrSet(key, newApp, val);

    this.setState({ app: newApp });
  }

  componentDidMount() {
    this.getApp();
  }

  render() {
    const { location: { pathname: currentPath }, classes } = this.props;

    const {
      loaded,
      showErrors,
      app,
    } = this.state;

    if (!loaded) {
      return <React.Fragment>Loading ...</React.Fragment>
    }

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
                schema={AppSchema.schema}
                form={AppSchema.form}
                model={app}
                onModelChange={this.onModelChange}
                showErrors={showErrors}
              />
            </Grid>
            <Grid item xs={12} className={classes.gridSaveButton}>
              <Button
                onClick={this.updateApp}
                variant="outlined"
                color="primary"
              >
                Save
            </Button>
            </Grid>
          </Grid>

        </div>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withSnackbar(AppView));
