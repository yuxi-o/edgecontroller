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

import React, { Component } from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import Typography from '@material-ui/core/Typography';
import { withRouter } from 'react-router-dom';

const styles = theme => ({
  sectionContainer: {
    marginTop: theme.spacing.unit * 4,
    marginBottom: theme.spacing.unit * 4
  },
  title: {
    fontWeight: 'bold'
  }
});

class SectionHeader extends Component {
  render() {
    const { classes, title, subtitle} = this.props;
    return (
      <div className={classes.sectionContainer}>
        <Typography variant="subtitle1" className={classes.title}>
          {title}
        </Typography>
        <Typography variant="body1" gutterBottom>
          {subtitle}
        </Typography>
      </div>
    )
  }
}

export default withRouter(withStyles(styles)(SectionHeader));
