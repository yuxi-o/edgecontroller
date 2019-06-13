import React from 'react';
import ApiClient from '../api/ApiClient';
import MockAdapter from 'axios-mock-adapter';
import { MemoryRouter } from 'react-router-dom';
import { mount } from 'enzyme';
import AppsView from './AppsListing';
import CardItem from '../components/cards/CardItem';
import {SnackbarProvider} from "notistack";

describe('AppsListing', () => {
  it('Renders a single app', (done) => {
    const axiosMock = new MockAdapter(ApiClient.axiosInstance);
    axiosMock.onGet('/apps').reply(200, [{
      'id': 'app-id-1',
      'type': 'vm',
      'name': 'vm app',
      'version': '1.2.3',
      'vendor': 'smartedge',
      'description': 'VM Application',
      'cores': 4,
      'memory': 4096,
      'source': 'http://super-cool-cdn.com'
    }]);

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', {value: spyScrollTo});
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/apps']}>
          <AppsView/>
        </MemoryRouter>
      </SnackbarProvider>
    );
    const wrapper2 = wrapper.find('AppsView');

    setTimeout(() => {
      expect(wrapper2.state()).toHaveProperty('loaded', true);
      const apps = wrapper2.state().apps;
      expect(apps).toHaveLength(1);
      expect(wrapper2.update().find(CardItem)).toHaveLength(1);

      done();
    });
  });

  it('Renders multiple apps', (done) => {
    const axiosMock = new MockAdapter(ApiClient.axiosInstance);
    axiosMock.onGet('/apps').reply(200, [{
      'id': 'app-id-1',
      'type': 'vm',
      'name': 'vm app',
      'version': '1.2.3',
      'vendor': 'smartedge',
      'description': 'VM Application',
      'cores': 4,
      'memory': 4096,
      'source': 'http://super-cool-cdn.com'
    }, {
      'id': 'app-id-2',
      'type': 'container',
      'name': 'container app',
      'version': '1.2.3',
      'vendor': 'smartedge',
      'description': 'Container Application',
      'cores': 4,
      'memory': 4096,
      'source': 'http://super-cool-cdn.com'
    }]);

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', {value: spyScrollTo});
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/apps']}>
          <AppsView/>
        </MemoryRouter>
      </SnackbarProvider>
    );
    const wrapper2 = wrapper.find('AppsView');

    setTimeout(() => {
      expect(wrapper2.state()).toHaveProperty('loaded', true);
      const apps = wrapper2.state().apps;
      expect(apps).toHaveLength(2);
      expect(wrapper2.update().find(CardItem)).toHaveLength(2);

      done();
    });
  });
});
