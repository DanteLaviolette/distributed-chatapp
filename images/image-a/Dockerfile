FROM golang:1.20.2-alpine

WORKDIR /backend

COPY backend/* ./

COPY ./backend/ ./
RUN go mod download

RUN go build -o /backend

ENV PORT 8000
ENV GO_ENV production

EXPOSE $PORT

CMD ["/backend"]