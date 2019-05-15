import React from 'react'
import {Route, BrowserRouter as Router} from 'react-router-dom'
import Dashboard from './components/Dashboard'
import Main from './components/Main'
import NodesView from './views/Nodes'
import AppsView from './views/Apps'
import VnfsView from './views/Vnfs'
import LoginForm from './components/Login'
import Wizard from './components/Wizard'
import ProtectedRoute from './components/ProtectedRoute'

export default props => (
  <Router>
    <div>
      <Route exact path='/login' component={LoginForm} />
      <ProtectedRoute exact path='/' component={Main} />
      <ProtectedRoute exact path='/home' component={Main} />
      <ProtectedRoute exact path='/dashboard' component={Dashboard} />
      <ProtectedRoute exact path='/nodes' component={NodesView} />
      <ProtectedRoute exact path='/vnfs' component={VnfsView} />
      <ProtectedRoute exact path='/apps' component={AppsView} />
      <ProtectedRoute exact path='/traffic-policies' component={Wizard} />
      <ProtectedRoute exact path='/dns-configs' component={Wizard} />
    </div>
  </Router>
)
