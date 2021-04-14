FROM golang:alpine as go
WORKDIR /app
ENV GO111MODULE=on

COPY ./backend .
COPY ./frontend/build ./build
RUN go build -o client .
EXPOSE 8080

CMD ["/app/client"]
