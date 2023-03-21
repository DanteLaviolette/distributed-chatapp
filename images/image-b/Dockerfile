FROM node:alpine as build

# create directory for src
WORKDIR /usr/src/frontend

# copy src to container
COPY ./frontend/ ./
# build react site from src
RUN npm install && npm run build

# pull nginx image
FROM nginx

# Move built site to nginx html dir
WORKDIR /usr/share/nginx/html
COPY --from=build /usr/src/frontend/dist/ ./

# Move nginx config to expected location
COPY ./nginx/app.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 30222