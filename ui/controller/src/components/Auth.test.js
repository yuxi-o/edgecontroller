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
import Auth from './Auth';
import MockAdapter from 'axios-mock-adapter';

describe('Auth Test', () => {
  it('Successful Login will store the Token...', async(done) => {
    const axiosMock = new MockAdapter(ApiClient.axiosInstance);

    axiosMock.onPost('/auth').reply(200, {
     'token': 'SUCCESS_TOKEN',
    });

    await Auth.login('test@email.org', '1234');
    expect(sessionStorage.setItem).toHaveBeenLastCalledWith('JWT', 'SUCCESS_TOKEN');
    expect(sessionStorage.__STORE__['JWT']).toEqual('SUCCESS_TOKEN');
    expect(Auth.isAuthenticated()).toEqual(true);
    expect(sessionStorage.getItem).toHaveBeenLastCalledWith('JWT');

    done();
  });

  it('Calling logout will remove the stored token', () => {
    sessionStorage.__STORE__['JWT'] = 'EXPIRED_TOKEN';
    expect(Auth.isAuthenticated()).toEqual(true);

    const mockFn = jest.fn();

    Auth.logout(mockFn);
    expect(Auth.isAuthenticated()).toEqual(false);
    expect(mockFn).toBeCalledTimes(1);
    expect(sessionStorage.removeItem).toHaveBeenCalledWith('JWT');
    expect(sessionStorage.__STORE__['JWT']).toBeUndefined()
  });
});
