import React,  { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import Topbar from '../../components/Topbar';
import SectionHeader from '../../components/typo/SectionHeader';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Typography from '@material-ui/core/Typography';

import SourceForm from './FormSource'
import DestinationForm from './FormDestination'
import TargetForm from './FormTarget'

const backgroundShape = require('../../images/shape.svg');

function TabContainer(props) {
  return (
    <Typography component="div" style={{ padding: 8 * 3 }}>
      {props.children}
    </Typography>
  );
}

TabContainer.propTypes = {
  children: PropTypes.node.isRequired,
};

const styles = theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.grey['A500'],
    overflow: 'hidden',
    background: `url(${backgroundShape}) no-repeat`,
    backgroundSize: 'cover',
    backgroundPosition: '0 400px',
    marginTop: 20,
    padding: 20,
    paddingBottom: 200
  },
  grid: {
    width: 1000
  }
})

class PolicyEdit extends Component {
  
  state = {
    value: 0,
  };

  handleChange = (event, value) => {
    this.setState({ value });
  };

  render() {
    const { classes } = this.props;
    const { value } = this.state;
    const currentPath = this.props.location.pathname

    return (
      <React.Fragment>
        <CssBaseline />
        <Topbar currentPath={currentPath} />
        <div className={classes.root}>
          <Grid container justify="center"> 
            <Grid spacing={24} alignItems="center" justify="center" container className={classes.grid}>
              <Grid item xs={12}>
                <SectionHeader title="Traffic Policies" subtitle="Select to View or Edit" />
                <AppBar position="static">
                  <Tabs value={value} onChange={this.handleChange}>
                    <Tab label="Source" />
                    <Tab label="Destination" />
                    <Tab label="Target" />
                  </Tabs>
                </AppBar>
                {value === 0 && <TabContainer><SourceForm title="Traffic Policies" subtitle="Select to View or Edit" /></TabContainer>}
                {value === 1 && <TabContainer><DestinationForm title="Traffic Policies" subtitle="Select to View or Edit" /></TabContainer>}
                {value === 2 && <TabContainer><TargetForm title="Traffic Policies" subtitle="Select to View or Edit" /></TabContainer>} 
              </Grid>
            </Grid>
          </Grid>
        </div>
      </React.Fragment>
    )
  }
}

export default withStyles(styles)(PolicyEdit);
