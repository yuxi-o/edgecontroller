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
import AppsView from './AppsListing';
import CardItem from '../components/cards/CardItem';
import { SnackbarProvider } from "notistack";

describe('AppsListing', () => {
  it('Renders a single app', (done) => {
    const axiosMock = new MockAdapter(ApiClient.axiosInstance);
    axiosMock.onGet('/apps').reply(200, {
      apps:
        [{
          "id": "CB0D7DA8-0B97-4668-9024-81415063A5C9",
          "type": "container",
          "name": "Sample App",
          "version": "1.2.3",
          "vendor": "Sample Vendor",
          "description": "Sample description goes here."
        }]
    });

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', { value: spyScrollTo });
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/apps']}>
          <AppsView />
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
    axiosMock.onGet('/apps').reply(200, {
      apps: [{
        "id": "CB0D7DA8-0B97-4668-9024-81415063A5C9",
        "type": "container",
        "name": "Sample App 1",
        "version": "1.2.3",
        "vendor": "Sample Vendor",
        "description": "Sample description goes here."
      }, {
        "id": "75a34d77-3de8-4d94-95ce-66d91e5c2815",
        "type": "container",
        "name": "Sample App 2",
        "version": "1.2.3",
        "vendor": "Sample Vendor",
        "description": "Sample description goes here."
      }]
    });

    const spyScrollTo = jest.fn();
    Object.defineProperty(global.window, 'scrollTo', { value: spyScrollTo });
    spyScrollTo.mockClear();
    sessionStorage.__STORE__['JWT'] = "FAKEJWTTOKEN";

    const wrapper = mount(
      <SnackbarProvider maxSnack={3}>
        <MemoryRouter initialEntries={['/apps']}>
          <AppsView />
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
