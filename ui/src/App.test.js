import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { mount, shallow } from 'enzyme';
import MuiThemeProvider from '@material-ui/core/styles/MuiThemeProvider';
import Routes from "./routes";
import MuiTheme from "./MuiTheme";
import Login from "./components/Login";

// Use React Dom to render App Component
it('renders without crashing (ReactDOM)', () => {
  const div = document.createElement('div');
  ReactDOM.render(<App />, div);
  ReactDOM.unmountComponentAtNode(div);
});

// Use Enzyme Shallow rendering for App Component
it('renders without crashing (Shallow)', () => {
  const sut = shallow(<App />);
  expect(sut.contains('div'));
});

// Use Enzyme Mount for App Component
it('renders without crashing (Mount)', () => {
  const sut = mount(<App />);
  const ThemeWrapper = (
    <MuiThemeProvider theme={MuiTheme}>
      <Routes />
    </MuiThemeProvider>
  );

  expect(sut.contains(ThemeWrapper)).toEqual(true);
});

// Validate Route exists

it('Prompts Login screen', () => {
  const sut = mount(<App />);

  expect(sut.contains(Login));

  // Expect total of 3 input field
  expect(sut.find('input')).toHaveLength(3);
});
