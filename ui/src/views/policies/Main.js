import React,  { Component } from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import Topbar from '../../components/Topbar';
import SectionHeader from '../../components/typo/SectionHeader';
import Table from '../../components/tables/EnhancedTable';

const backgroundShape = require('../../images/shape.svg');

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

function createData(id, name) {
  return { id, name, editUrl: `/policies/${id}/edit` };
}

const tableHeaders = [
  { id: 'id', numeric: true, disablePadding: false, label: 'Id' },
  { id: 'name', numeric: false, disablePadding: false, label: 'Name' },
  { id: 'action', numeric: false, disablePadding: false, label: 'Action' }
];

const tableData = {
  order: 'asc',
  orderBy: 'id',
  selected: [],
  data: [
    createData('1', 'Policy 1'),
    createData('2', 'Policy 2'),
    createData('3', 'Policy 3'),
    createData('4', 'Policy 4'),
    createData('5', 'Policy 5'),
    createData('6', 'Policy 6')
  ],
  page: 0,
  rowsPerPage: 5,
};

class PoliciesView extends Component {

  render() {
    const { classes } = this.props;
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
                <Table rows={tableHeaders} tableState={tableData} />
              </Grid>
            </Grid>
          </Grid>
        </div>
      </React.Fragment>
    )
  }
}

export default withStyles(styles)(PoliciesView);
