import React from 'react'
import { Route, Redirect } from 'react-router-dom'
import Dashboard from './components/Dashboard'
import Main from './components/Main'
import NodesView from './views/NodesListing'
import NodeView from './views/Node'
import AppsView from './views/AppsListing'
import AppView from './views/App'
import LoginForm from './components/Login'

import Dns from './views/dns/Main'
import DnsEdit from './views/dns/Edit'

import Policies from './views/policies/Main'
import PoliciesEdit from './views/policies/Edit'

import ProtectedRoute from './components/ProtectedRoute'
import Auth from './components/Auth'

export default props => (
  <div>
    <Route
      exact
      path='/'
      render={() => (
        Auth.isAuthenticated()
          ? <Redirect to="/home" />
          : <Redirect to="/login" />
      )}
    />

    <Route exact path='/login' component={LoginForm} />
    <ProtectedRoute exact path='/home' component={Main} />
    <ProtectedRoute exact path='/dashboard' component={Dashboard} />

    <ProtectedRoute exact path='/nodes' component={NodesView} />
    <ProtectedRoute path='/nodes/:id' component={NodeView} />

    <ProtectedRoute exact path='/apps' component={AppsView} />
    <ProtectedRoute path='/apps/:id' component={AppView} />

    <ProtectedRoute exact path='/policies' component={Policies} />
    <ProtectedRoute exact path='/policies/add' component={PoliciesEdit} />
    <ProtectedRoute path='/policies/:id/edit' component={PoliciesEdit} />

    <ProtectedRoute exact path='/dns' component={Dns} />
    <ProtectedRoute exact path='/dns/add' component={DnsEdit} />
    <ProtectedRoute path='/dns/:id/edit' component={DnsEdit} />
  </div>
)
