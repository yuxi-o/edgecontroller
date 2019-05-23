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

class TextFields extends React.Component {

  handleChange = name => event => {
    this.setState({ [name]: event.target.value });
  };

  render() {
    const { classes } = this.props;

    return (
      <form className={classes.container} noValidate autoComplete="off">
        <Typography variant="h6" className={classes.fieldSet} gutterBottom>
          DNS Record
        </Typography>
                
        <TextField
          required
          id="standard-name"
          label="Name"
          className={classes.textField}
          margin="normal"
        />

        <TextField
          required
          id="standard-type"
          label="Type"
          className={classes.textField}
          margin="normal"
        />

        <TextField
          required
          id="standard-host"
          label="Host"
          className={classes.textField}
          margin="normal"
        />

        <TextField
          required
          id="standard-error"
          label="TTL"
          className={classes.textField}
          margin="normal"
        />

        <div className={classes.spacer} />

        <Typography variant="h6" className={classes.fieldSet} gutterBottom>
          DNS Forwarders
        </Typography>
                
        <TextField
          id="standard-name"
          label="Dns Forwarders"
          className={classes.textField}
          onChange={this.handleChange('name')}
          margin="normal"
        />

        <div className={classes.spacer} />
        
        <Button variant="contained" color="secondary" className={classes.button}>
          Submit
        </Button>
      </form>
    );
  }
}

TextFields.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(TextFields);
