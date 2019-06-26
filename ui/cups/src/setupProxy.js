const proxy = require('http-proxy-middleware');

module.exports = function (app) {
  console.log("PROXY " + process.env.REACT_APP_CUPS_API);
  app.use(proxy('/api', { target: process.env.REACT_APP_CUPS_API, pathRewrite: { '/api': '' } }));
};
