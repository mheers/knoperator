ARG go="golang:1.23.3-alpine3.20"
ARG base="alpine:3.20.3"

# third stage for building the backend
FROM --platform=$BUILDPLATFORM ${go} AS builder

RUN apk add bash git gcc libc-dev openssh

# Copy the code from the host and compile it

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

ADD . ./

RUN [ "$(uname)" = Darwin ] && system=darwin || system=linux; \
    ./ci/go-build.sh --os ${system} --arch $(echo $TARGETPLATFORM  | cut -d/ -f2)

# final stage
FROM ${base}

COPY --from=builder /go/src/app/goapp /bin/knoperator

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENTRYPOINT [ "/bin/knoperator" ]
