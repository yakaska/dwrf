FROM golang:alpine as build

WORKDIR /build

COPY . .

RUN go mod download && go build -o tinyurl cmd/tinyurl/main.go

FROM alpine:latest

ENV CONFIG_PATH /home/root/dev.yaml

COPY --from=build /build/tinyurl /usr/local/bin/tinyurl

COPY --from=build /build/config/dev.yaml /home/root/dev.yaml

EXPOSE 4000

ENTRYPOINT ["/usr/local/bin/tinyurl"]
