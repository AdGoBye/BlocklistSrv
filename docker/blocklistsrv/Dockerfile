# syntax=docker/dockerfile:1
FROM golang:1.22-alpine
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /bin/blocklistsrv .
ENTRYPOINT [ "/bin/blocklistsrv" ]