FROM node:alpine as build

# create directory for src
WORKDIR /usr/src/frontend

# copy src to container
COPY ./frontend/ ./
# install dependencies & build
RUN npm install && npm run build

# pull nginx image
FROM nginx

# move built site to nginx html dir
WORKDIR /usr/share/nginx/html
COPY --from=build /usr/src/frontend/dist/ ./

# move nginx config to expected location
COPY ./nginx/app.conf /etc/nginx/conf.d/default.conf

# expose port
EXPOSE 30222