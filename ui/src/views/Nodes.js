import React,  { Component } from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import { withRouter } from "react-router-dom";
import NodeCard from '../components/cards/NodeCard';
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

class NodesView extends Component {

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
                <SectionHeader title="Nodes" subtitle="View of nodes" />
                <NodeCard nodeID={1} />
                <NodeCard nodeID={2} />
                <NodeCard nodeID={3} />
                <NodeCard nodeID={4} />
              </Grid>
            </Grid>
          </Grid>
        </div>
      </React.Fragment>
    )
  }
}

export default withRouter(withStyles(styles)(NodesView));
