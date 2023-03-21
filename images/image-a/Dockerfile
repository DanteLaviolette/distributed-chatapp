FROM golang:1.20.2-alpine

# create directory for src 
WORKDIR /usr/src/backend

# copy src to container
COPY ./backend/ ./

# install dependencies & build
RUN go mod download && go build -o /backend

# set env vars
ENV PORT 8000
ENV GO_ENV production

# expose port & start
EXPOSE $PORT
CMD ["/backend"]