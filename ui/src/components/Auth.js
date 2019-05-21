class Auth {
  static login(email, password) {
    //TODO: Hook up JWT Auth here
    if (email !== password) {
      return false;
    }

    sessionStorage.setItem('AuthToken', email);
    return this.isAuthenticated();
  }

  static logout(cb) {
    return cb(sessionStorage.removeItem('AuthToken'));
  }

  static isAuthenticated() {
    return !!sessionStorage.getItem('AuthToken');
  }
}
export default Auth;
