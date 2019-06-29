// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React,  { Component } from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import CircularProgress from '@material-ui/core/CircularProgress'
import Grid from "@material-ui/core/Grid";

const styles = theme => ({
  container: {
    position: 'absolute',
    width: '100%',
    height: '100%',
  },
  flexGrow: {
    height: '100%',
  }
});

class FullSizeCircularLoader extends Component {

  render() {
    const { classes } = this.props;
    return (
      <div className={classes.container}>
        <Grid container className={classes.flexGrow} direct="column" alignItems="center" justify="center">
          <CircularProgress className={classes.progress} />
        </Grid>
      </div>
    )
  }
}

export default withStyles(styles)(FullSizeCircularLoader)
