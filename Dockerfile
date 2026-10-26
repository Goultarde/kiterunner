FROM golang:1.25-alpine AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/kr ./cmd/kiterunner

FROM alpine:3.22

RUN adduser -D -H -u 10001 kiterunner
COPY --from=build /out/kr /usr/local/bin/kr

USER kiterunner
ENTRYPOINT ["/usr/local/bin/kr"]
