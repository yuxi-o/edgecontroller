import React from 'react'
import { Route, Redirect, BrowserRouter as Router } from 'react-router-dom'
import Dashboard from './components/Dashboard'
import Main from './components/Main'
import NodesView from './views/Nodes'
import NodeView from './views/Node'
import AppsView from './views/Apps'
import AppView from './views/App'
import VnfsView from './views/Vnfs'
import VnfView from './views/Vnf'
import LoginForm from './components/Login'
import Wizard from './components/Wizard'
import ProtectedRoute from './components/ProtectedRoute'
import Auth from './components/Auth'

export default props => (
  <Router>
    <div>
      <Route exact path='/' render={() => (
        (Auth.isAuthenticated()) ? (
          <Redirect to="/home"/>
        ) : (
          <Redirect to="/login"/>
        )
      )}/>

      <Route exact path='/login' component={LoginForm} />
      <ProtectedRoute exact path='/home' component={Main} />
      <ProtectedRoute exact path='/dashboard' component={Dashboard} />
      <ProtectedRoute exact path='/nodes' component={NodesView} />
      <ProtectedRoute path='/nodes/:id' component={NodeView} />
      <ProtectedRoute exact path='/vnfs' component={VnfsView} />
      <ProtectedRoute path='/vnfs/:id' component={VnfView} />
      <ProtectedRoute exact path='/apps' component={AppsView} />
      <ProtectedRoute path='/apps/:id' component={AppView} />
      <ProtectedRoute exact path='/traffic-policies' component={Wizard} />
      <ProtectedRoute exact path='/dns-configs' component={Wizard} />
    </div>
  </Router>
)
