FROM golang:1.14-alpine AS build_base

# setup of thorcast
RUN apk add bash git
WORKDIR /go/src/github.com/kylep342/thorcast-server
# ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

# compile the binary
FROM build_base AS server_builder

COPY pkg/app/ .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o thorserver
RUN go build -o thorserver

# execute the binary
FROM alpine
RUN apk add ca-certificates
# RUN apk add ca-certificates tzdata
# ENV TZ UTC
COPY --from=server_builder /go/src/github.com/kylep342/thorcast-server/thorserver /bin/thorserver
EXPOSE 8000

# run thorcast
ENTRYPOINT ["/bin/thorserver"]