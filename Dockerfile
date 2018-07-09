FROM golang:1.9.2-alpine3.6
WORKDIR /go/src/github.com/wildsurfer/postgrest-oauth-server
COPY . .
RUN apk add --no-cache openssl git && \
    wget -O /usr/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && \
    chmod +x /usr/bin/dep /usr/bin/dep && \
    dep ensure -vendor-only && go build

FROM alpine:3.6
MAINTAINER Ivan Kuznetsov <kuzma.wm@gmail.com>
ENV OAUTH_DB_CONN_STRING="postgres://user:pass@postgresql:5432/test?sslmode=disable" \
    OAUTH_ACCESS_TOKEN_SECRET="morethan32symbolssecretkey!!!!!!" \
    OAUTH_ACCESS_TOKEN_TTL=7200 \
    OAUTH_REFRESH_TOKEN_SECRET="notlesshan32symbolssecretkey!!!!" \
    OAUTH_COOKIE_HASH_KEY="supersecret" \
    OAUTH_COOKIE_BLOCK_KEY="16charssecret!!!" \
    OAUTH_VALIDATE_REDIRECT_URI=true \
    OAUTH_CODE_UI="http://localhost:3685"
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/wildsurfer/postgrest-oauth-server/postgrest-oauth-server .
CMD ./postgrest-oauth-server
