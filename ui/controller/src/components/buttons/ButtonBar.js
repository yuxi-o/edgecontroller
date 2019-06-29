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
import Button from '@material-ui/core/Button';

const styles = theme => ({
  primary: {
    marginRight: theme.spacing.unit * 2
  },
  secondary: {
    background: theme.palette.secondary['100'],
    color: 'white'
  },
  spaceTop: {
    marginTop: 20
  }
});

class ButtonBar extends Component {
  constructor(props) {
    super(props);
    this.handlePrimaryAction = this.handlePrimaryAction.bind(this);
  }

  handlePrimaryAction(e) {
    this.props.primaryButtonAction(e);
  }

  render() {
    const {
      classes,
      primaryButtonName,
      secondaryButtonName,
      primaryLink,
      secondaryLink
    } = this.props;

    return (
      <div className={classes.spaceTop}>
        <Button
          className={classes.primary}
          component={primaryLink}
          onClick={this.handlePrimaryAction}
        >
          {primaryButtonName}
        </Button>
        <Button
          variant="contained"
          color="primary"
          className={classes.secondary}
          component={secondaryLink}
        >
          {secondaryButtonName}
        </Button>
      </div>
    )
  }
}

export default withStyles(styles)(ButtonBar);
