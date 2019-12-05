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

import React from 'react'
import { Switch, Route, Redirect } from 'react-router-dom'
import NodesView from './views/NodesListing'
import NodeView from './views/Node'
import AppsView from './views/AppsListing'
import AppView from './views/AppView'
import LoginForm from './components/Login'
import Landing from './components/Landing'

// import Dns from './views/dns/Main'
// import DnsEdit from './views/dns/Edit'

import Policies from './views/policies/Main'
import PoliciesEdit from './views/policies/Edit'

import ProtectedRoute from './components/ProtectedRoute'

export default props => (
  <div>
    <Switch>
      <Route
        exact
        path='/'
        render={() => <Redirect to="/landing" />}
      />

      <Route exact path='/login' component={LoginForm} />

      <Route exact path='/landing' component={Landing} />

      <ProtectedRoute exact path='/nodes' component={NodesView} />
      <ProtectedRoute path='/nodes/:id' component={NodeView} />

      <ProtectedRoute exact path='/apps' component={AppsView} />
      <ProtectedRoute path='/apps/:id' component={AppView} />

      <ProtectedRoute exact path='/policies' component={Policies} />
      <ProtectedRoute exact path='/policies/add' component={PoliciesEdit} />
      <ProtectedRoute path='/policies/:id/edit' component={PoliciesEdit} />

      <Route render={() => <div>404 Not Found</div>} />
    </Switch>
  </div>
)
