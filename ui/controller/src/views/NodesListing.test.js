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
import ApiClient from '../api/ApiClient';
import MockAdapter from 'axios-mock-adapter';
import { MemoryRouter } from 'react-router-dom';
import { mount } from 'enzyme';
import NodesView from './NodesListing';
import CardItem from '../components/cards/CardItem';
import { SnackbarProvider } from "notistack";

describe('NodesListing', () => {
  it('Renders a single node', (done) => {
    const axiosMock = new MockAdapter(ApiClient.axiosInstance);
    axiosMock.onGet('/nodes').reply(200, {
      nodes: [{
        "id": "CB0D7DA8-0B97-4668-9024-81415063A5C9",
        "name": "sample-name-us-west-2",
        "location": "Irvine, CA",
        "serial": "CB0D7DA8-0B97-4668-9024-81415063A5C9"
      }]
    });

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', { value: spyScrollTo });
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/nodes']}>
          <NodesView />
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
    axiosMock.onGet('/nodes').reply(200, {
      nodes: [{
        "id": "CB0D7DA8-0B97-4668-9024-81415063A5C9",
        "name": "sample-name-us-west-2",
        "location": "Irvine, CA",
        "serial": "CB0D7DA8-0B97-4668-9024-81415063A5C9"
      }, {
        "id": "7115cb33-0262-4b31-b9b9-4be96e4b0059",
        "name": "sample-name-us-west-3",
        "location": "Irvine, CA",
        "serial": "59ab1412-de83-4207-b002-c531fdb949ae"
      }]
    });

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', { value: spyScrollTo });
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/nodes']}>
          <NodesView />
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
