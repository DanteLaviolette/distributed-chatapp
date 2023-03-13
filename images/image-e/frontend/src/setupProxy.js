const { createProxyMiddleware } = require("http-proxy-middleware");
 
module.exports = function(app) {
    app.use(createProxyMiddleware('/api_login/*', { target: 'http://localhost:8000/' }))
    app.use(createProxyMiddleware('/api_register/*', { target: 'http://localhost:8001/' }))
}