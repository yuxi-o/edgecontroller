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

  render() {
    const {
      classes,
      primaryButtonName,
      secondaryButtonName,
      primaryLink,
      secondaryLink
    } = this.props;

    console.log(this.props);

    return (
      <div className={classes.spaceTop}>
        <Button
          className={classes.primary}
          component={primaryLink}
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
