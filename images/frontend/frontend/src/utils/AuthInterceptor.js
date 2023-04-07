/*
This is a common pattern I use for auth, so code is basically the same as
another project I created:
- https://github.com/EECS4481Project/frontend/blob/main/src/auth/RequestInterceptor.js
*/
import axios from 'axios';
import { deleteAuthJWT, getAuthJWT, setAuthJWT } from './utils';
import constants from '../constants'
// Interceptor to be used for all requests that require authorization
const authorizedAxios = axios.create();

// Set auth token in request
authorizedAxios.interceptors.request.use((config) => {
  const authToken = getAuthJWT();
  if (authToken) {
    config.headers = { ...config.headers, [constants.AUTH_HEADER]: authToken };
  }
  return config;
}, (err) => {
  Promise.reject(err);
});

// Set auth token if found in response
authorizedAxios.interceptors.response.use((res) => {
  const authToken = res.headers[constants.AUTH_HEADER];
  if (authToken) {
    setAuthJWT(authToken);
  }
  return res;
}, (err) => {
  // Unauthorized, delete auth token
  if (err.response.status === 401) {
    deleteAuthJWT()
  }
  return Promise.reject(err);
});

export default authorizedAxios;