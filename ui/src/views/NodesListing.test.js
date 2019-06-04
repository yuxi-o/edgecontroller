import React from 'react';
import ApiClient from '../api/ApiClient';
import MockAdapter from 'axios-mock-adapter';
import { MemoryRouter } from 'react-router-dom';
import { mount } from 'enzyme';
import NodesView from './NodesListing';
import CardItem from '../components/cards/CardItem';
import {SnackbarProvider} from "notistack";

describe('NodesListing', () => {
  it('Renders a single node', (done) => {
    const axiosMock = new MockAdapter(ApiClient.axiosInstance);
    axiosMock.onGet('/nodes').reply(200, [{
      'id': 'node-id',
      'name': 'Mock Node',
      'location': 'Irvine, CA',
      'serial': 'node-serial-uuid',
      'grpc_target': 'GRPC_TARGET',
    }]);

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', {value: spyScrollTo});
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/nodes']}>
          <NodesView/>
        </MemoryRouter>
      </SnackbarProvider>
    );
    const wrapper2 = wrapper.find('NodesView');

    setTimeout(() => {
      expect(wrapper2.state()).toHaveProperty('loaded', true);
      const nodes = wrapper2.state().nodes;
      expect(nodes).toHaveLength(1);
      expect(wrapper2.update().find(CardItem)).toHaveLength(1);

      done();
    });
  });

  it('Renders multiple nodes', (done) => {
    const axiosMock = new MockAdapter(ApiClient.axiosInstance);
    axiosMock.onGet('/nodes').reply(200, [{
      'id': 'node-id',
      'name': 'Mock Node',
      'location': 'Irvine, CA',
      'serial': 'node-serial-uuid',
      'grpc_target': 'GRPC_TARGET',
    }, {
      'id': 'node-id-2',
      'name': 'Mock Node 2',
      'location': 'Irvine, CA',
      'serial': 'node-serial-uuid-2',
      'grpc_target': 'GRPC_TARGET',
    }]);

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', {value: spyScrollTo});
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/nodes']}>
          <NodesView/>
        </MemoryRouter>
      </SnackbarProvider>
    );
    const wrapper2 = wrapper.find('NodesView');

    setTimeout(() => {
      expect(wrapper2.state()).toHaveProperty('loaded', true);
      const nodes = wrapper2.state().nodes;
      expect(nodes).toHaveLength(2);
      expect(wrapper2.update().find(CardItem)).toHaveLength(2);

      done();
    });
  });
});
