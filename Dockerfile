FROM golang:alpine AS gobuild
ADD . /go/src/app
RUN apk update && apk add git
RUN cd /go/src/app && go get && go build -o oauth2_proxy-token

FROM alpine

COPY --from=gobuild /go/src/app/oauth2_proxy-token /app/
WORKDIR /app

RUN adduser -D oauth2_proxy

USER oauth2_proxy

ENTRYPOINT ["./oauth2_proxy-token"]