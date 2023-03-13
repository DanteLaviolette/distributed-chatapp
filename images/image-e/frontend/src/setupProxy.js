const { createProxyMiddleware } = require("http-proxy-middleware");
 
module.exports = function(app) {
    app.use(createProxyMiddleware('/api/login', { target: 'http://localhost:8000/' }))
    app.use(createProxyMiddleware('/api/register', { target: 'http://localhost:8001/' }))
}