import React, { Component } from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import Topbar from '../../components/Topbar';
import Table from '../../components/tables/EnhancedTable';
import AddIcon from '@material-ui/icons/Add';
import {
  Typography,
  Button,
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
    width: 1000
  },
  addButton: {
    float: 'right',
  },
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
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
    };
  }

  handleOnClickAddPolicy = () => {
    const { history } = this.props;

    // Redirect user to add policies view
    history.push('/policies/add');
  };

  render() {
    const { classes } = this.props;
    const currentPath = this.props.location.pathname

    const policiesGrid = () => (
      <Grid spacing={24} alignItems="center" justify="center" container className={classes.grid}>
        <Grid item xs={12}>
          <Grid container direction="row"
            justify="space-between"
            alignItems="flex-start"
            className={classes.sectionContainer}
          >
            <Grid item>
              <Typography variant="subtitle1" className={classes.title}>
                Traffic Policies
                </Typography>
              <Typography variant="body1" gutterBottom className={classes.subtitle}>
                List of Traffic Policies
                </Typography>
            </Grid>
            <Grid item xs={3}>
              <Button variant="contained" color="primary" className={classes.addButton} onClick={this.handleOnClickAddPolicy}>
                Add Policy
                  <AddIcon className={classes.rightIcon} />
              </Button>
            </Grid>
          </Grid>
          <Table rows={tableHeaders} tableState={tableData} />
        </Grid>
      </Grid>
    );

    return (
      <React.Fragment>
        <CssBaseline />
        <Topbar currentPath={currentPath} />
        <div className={classes.root}>
          <Grid container justify="center">
            {policiesGrid()}
          </Grid>
        </div>
      </React.Fragment>
    )
  }
}

export default withStyles(styles)(PoliciesView);
