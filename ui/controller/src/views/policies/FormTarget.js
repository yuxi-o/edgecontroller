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
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
    width: 200,
  },
  fieldSet: {
    width: "100%",
  },
  dense: {
    marginTop: 19,
  },
  menu: {
    width: 200,
  },
  spacer: {
    width: "100%",
    padding: 20,
  }
});


class SourceForm extends React.Component {

  handleChange = name => event => {
    this.setState({ [name]: event.target.value });
  };

  render() {
    const { classes } = this.props;

    return (
      <form className={classes.container} noValidate autoComplete="off">
        <Typography variant="h6" className={classes.fieldSet} gutterBottom>
          Mac Filter
        </Typography>
                
        <TextField
          required
          id="standard-name"
          label="Description"
          className={classes.textField}
          margin="normal"
        />

        <div className={classes.spacer}></div>

        <TextField
          required
          id="standard-type"
          label="Target Action"
          className={classes.textField}
          margin="normal"
        />

        <div className={classes.spacer}></div>

        <TextField
          required
          id="standard-host"
          label="MAC Modifier"
          className={classes.textField}
          margin="normal"
        />

        <div className={classes.spacer}></div>

        <TextField
          required
          id="standard-error"
          label="IP Modifier"
          className={classes.textField}
          margin="normal"
        />

        <div className={classes.spacer}></div>

        <Button variant="contained" color="secondary" className={classes.button}>
          Submit
        </Button>
      </form>
    );
  }
}

SourceForm.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(SourceForm);
