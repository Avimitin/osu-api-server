#BUILD
FROM golang:alpine AS build
COPY . /go/src/github.com/avimitin/osuapiserver
WORKDIR /go/src/github.com/avimitin/osuapiserver
RUN go build -o /bin/oas-linux-amd64 -ldflags '-s -w' cmd/cmd.go

#RUN
FROM alpine:3
COPY --from=build /bin/oas-linux-amd64 /bin/oas-linux-amd64
ENV OSU_CONF_PATH=/data
ENTRYPOINT ["/bin/oas-linux-amd64"]
