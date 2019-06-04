import axios from 'axios';

class ApiClient {
  _CONTROLLER_API_ = process.env.REACT_APP_CONTROLLER_API;
  _CUPS_API = process.env.REACT_APP_CUPS_API;
  _CUPS_UI_URL = process.env.REACT_APP_CUPS_API;

  axiosConfig = {
    baseURL: (process.env.NODE_ENV === 'production') ? this._CONTROLLER_API_ : '/api',
    timeout: 10000,
    headers: {
      contentType: 'application/json',
      accept: 'application/json',
      'Authorization': `Bearer ${this.getJWT()}`,
    }
  };

  axiosInstance = axios.create(this.axiosConfig);

  /**
   *
   * @returns string - The AuthToken if present
   */
  getJWT() {
    return sessionStorage.getItem('JWT');
  }

  /**
   *
   * @param Token
   */
  setJWT(Token) {
    sessionStorage.setItem('JWT', Token);

    // Whenever setJWT is called we need to update the axios instance
    this.updateAxios();
  }

  updateAxios() {
    this.axiosConfig = {...this.axiosConfig, headers: {
        contentType: 'application/json',
        accept: 'application/json',
        'Authorization': `Bearer ${this.getJWT()}`,
      }};

    this.axiosInstance = axios.create(this.axiosConfig);
  }

  /**
   * @async
   * @param username
   * @param password
   * @returns {Promise<AxiosResponse>}
   */
  async login(username, password) {
    return await this.axiosInstance.post('/auth', {username, password});
  }

  /**
   * @async
   * @param {string} path - The Path
   * @param {*} options - Axios options
   * @returns {Promise<AxiosResponse>}
   */
  async get(path, options = {}) {
    return await this.axiosInstance.get(path, options);
  }

  /**
   * @async
   * @param {string} path - The Path
   * @param {*} data - Post Body
   * @returns {Promise<AxiosResponse>}
   */
  async post(path, data = {}) {
    return await this.axiosInstance.post(path, data);
  }

  /**
   * @async
   * @param {string} path - The Path
   * @param {*} options - Axios options
   * @returns {Promise<AxiosResponse>}
   */
  async delete(path, options = {}) {
    return await this.axiosInstance.delete(path, options);
  }
}

export default new ApiClient();
