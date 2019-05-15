import React, { Component } from 'react';
import { MuiThemeProvider } from '@material-ui/core/styles';
import './App.css';
import Routes from './routes'
import MuiTheme from './MuiTheme';

class App extends Component {
  render() {
    return (
      <div>
        <MuiThemeProvider theme={MuiTheme}>
          <Routes />
        </MuiThemeProvider>
      </div>
    );
  }
}

export default App;
