FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=amd64

RUN apk update && apk upgrade
RUN apk add upx gcc musl-dev

WORKDIR /src
COPY . .

RUN go build \
  -ldflags "-s -w -extldflags '-static'" \
  -o /bin/aac \
  . \
  && strip /bin/aac \
  && upx -q -9 /bin/aac

FROM scratch

COPY --from=builder /bin/aac /aac

ENTRYPOINT ["/aac"]