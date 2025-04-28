FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN apk update && apk upgrade
RUN apk add upx gcc musl-dev

WORKDIR /src
COPY . .

RUN go build \
  -ldflags "-s -w" \
  -o /bin/aac \
  . \
  && strip /bin/aac \
  && upx -q -9 /bin/aac

FROM scratch

COPY --from=builder /bin/aac /aac

ENTRYPOINT ["/aac"]