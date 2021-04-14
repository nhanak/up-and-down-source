FROM golang:alpine as go
WORKDIR /app
ENV GO111MODULE=on

COPY . .
RUN go build -o cloud-mysql .

CMD ["/app/cloud-mysql"]
