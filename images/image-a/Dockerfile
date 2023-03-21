FROM golang:1.20.2-alpine

# create directory for src 
WORKDIR /usr/src/backend

# copy src to container
COPY ./backend/ ./

# install dependencies & build
RUN go mod download && go build -o ./backend

# set env vars
ENV PORT=8000
ENV GO_ENV production
# IRL our secrets would be provided by the executing env (ie. AWS)
# although as we aren't hosting
# the app, we'll store them in plaintext here :(
ENV AUTH_PRIVATE_KEY="sDD6oQm2I0u4I/rNutvCK+LemfgBLO6n8nYlMz1P9XyEqWaFmvPz+NOXS2cra8m+VUjnLgr/NzNPKsxRD7YMA"
ENV REFRESH_PRIVATE_KEY="PkcS1H49SOR5f3HlmYJFtrOVgr9Rcf6jmNRuUuTjnXJhbmHPzbo/uA9h/hRUd757QCHOacvxRO1cNPvqDcb44Q"

# NOTE: MONGODB_URL, REDIS_URL, REDIS_PASSWORD, REDIS_DB,
# REDIS_MESSAGING_CHANNEL, REDIS_USER_COUNT_CHANNEL, 
# ANONYMOUS_USERS_REDIS_KEY & AUTHORIZED_USERS_REDIS_KEY
# should be set by kubernetes config to decouple the
# services.

# expose port & start
EXPOSE $PORT
CMD ["./backend"]