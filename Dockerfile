FROM golang:1.12-alpine AS build_base

# setup of thorcast
RUN apk add bash git
WORKDIR /go/src/github.com/kylep342/thorcast-server
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

#
FROM build_base AS server_builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o thorserver

#
FROM alpine
RUN apk add ca-certificates
COPY --from=server_builder /go/src/github.com/kylep342/thorcast-server/thorserver /bin/thorserver
EXPOSE 8000

# run thorcast
ENTRYPOINT ["/bin/thorserver"]