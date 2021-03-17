FROM golang:1.15-buster
COPY . /go/src/github.com/avimitin/osuapiserver
WORKDIR /go/src/github.com/avimitin/osuapiserver
RUN go build -o /bin/osuapi-linux -ldflags '-s -w' cmd/cmd.go
ENV OSU_CONF_PATH=/data
ENTRYPOINT ["bin/osuapi-linux"]
