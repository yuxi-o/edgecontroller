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

import React from 'react';
import Typography from '@material-ui/core/Typography';
import withStyles from '@material-ui/core/styles/withStyles';

const styles = theme => ({
  loadingMessage: {
    position: 'absolute',
    top: '40%',
    left: '40%'
  }
});

function Loading(props) {
  const { classes, loading } = props;
  return (
    <div style={loading ? { display: 'block' } : { display: 'none' }} className={classes.loadingMessage}>
      <span role='img' aria-label='emoji' style={{ fontSize: 58, textAlign: 'center', display: 'inline-block', width: '100%' }}>👋</span>
      <Typography variant="h6">
        Waiting for input
      </Typography>
    </div>
  );
}

export default withStyles(styles)(Loading);
