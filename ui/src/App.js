import React, { Component } from 'react';
import { BrowserRouter } from "react-router-dom";
import { MuiThemeProvider } from '@material-ui/core/styles';
import './App.css';
import Routes from './routes'
import MuiTheme from './MuiTheme';
import '@material-ui/icons';
import 'typeface-lato';
import 'typeface-roboto';
import {SnackbarProvider} from "notistack";

class App extends Component {
  render() {
    return (
      <SnackbarProvider maxSnack={3}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'right'
        }}
      >
        <div>
          <MuiThemeProvider theme={MuiTheme}>
            <BrowserRouter>
              <Routes />
            </BrowserRouter>
          </MuiThemeProvider>
        </div>
      </SnackbarProvider>
    );
  }
}

export default App;
