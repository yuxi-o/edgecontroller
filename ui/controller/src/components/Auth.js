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

import ApiClient from '../api/ApiClient';

class Auth {
  static async login(email, password) {

    try {
      const authResp = await ApiClient.login(email, password);

      if (!authResp.data.token) {
        return false;
      }

      ApiClient.setJWT(authResp.data.token);
      return { success: true };
    } catch (err) {
      if (err.response && "data" in err.response) {
        return { success: false, errorText: err.response.data };
      }

      return { success: false, errorText: "Login Failed Try again Later" };
    }
  }

  static logout(cb) {
    return cb(sessionStorage.removeItem('JWT'));
  }

  static isAuthenticated() {
    return !!sessionStorage.getItem('JWT');
  }
}
export default Auth;
