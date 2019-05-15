class Auth {
  login(email, password) {
    //TODO: Hook up JWT Auth here
    if (email !== password) {
      return this.authenticated = false;
    }

    sessionStorage.setItem('AuthToken', email);
    return this.authenticated = true;
  }

  logout(cb) {
    return cb(sessionStorage.removeItem('AuthToken'));
  }

  isAuthenticated() {
    return !!sessionStorage.getItem('AuthToken');
  }
}
export default new Auth();
