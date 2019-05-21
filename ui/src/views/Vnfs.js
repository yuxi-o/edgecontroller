import React,  { Component } from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import VnfCard from '../components/cards/VnfCard';
import Topbar from '../components/Topbar';
import SectionHeader from '../components/typo/SectionHeader';
const backgroundShape = require('../images/shape.svg');

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

class VnfsView extends Component {

  render() {
    const { location: {pathname: currentPath}, classes } = this.props;

    return (
      <React.Fragment>
        <CssBaseline />
        <Topbar currentPath={currentPath} />
        <div className={classes.root}>
          <Grid container justify="center">
            <Grid spacing={24} alignItems="center" justify="center" container className={classes.grid}>
              <Grid item xs={12}>
                <SectionHeader title="Vnfs" subtitle="View of VNFs" />
                <VnfCard vnfID={1} />
                <VnfCard vnfID={2} />
                <VnfCard vnfID={3} />
                <VnfCard vnfID={4} />
              </Grid>
            </Grid>
          </Grid>
        </div>
      </React.Fragment>
    )
  }
}

export default withStyles(styles)(VnfsView);
