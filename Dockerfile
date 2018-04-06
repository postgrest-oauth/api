FROM golang:1.9.2-alpine3.6
WORKDIR /go/src/github.com/wildsurfer/postgrest-oauth-server
COPY . .
RUN apk add --no-cache openssl git && \
    wget -O /usr/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && \
    chmod +x /usr/bin/dep /usr/bin/dep && \
    dep ensure -vendor-only && go build

FROM node:9.10.1-alpine
WORKDIR /app
COPY ./react-app /app
RUN yarn install
RUN yarn build

FROM alpine:3.6
MAINTAINER Ivan Kuznetsov <kuzma.wm@gmail.com>
ENV OAUTH_DB_CONN_STRING="postgres://user:pass@postgresql:5432/test?sslmode=disable" \
    OAUTH_ACCESS_TOKEN_JWT_SECRET="morethan32symbolssecretkey!!!!!!" \
    OAUTH_ACCESS_TOKEN_TTL=7200 \
    OAUTH_REFRESH_TOKEN_JWT_SECRET="notlesshan32symbolssecretkey!!!!" \
    OAUTH_COOKIE_HASH_KEY="supersecret" \
    OAUTH_COOKIE_BLOCK_KEY="16charssecret!!!" \
    OAUTH_VALIDATE_REDIRECT_URI=true
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/wildsurfer/postgrest-oauth-server/postgrest-oauth-server .
COPY --from=0 /go/src/github.com/wildsurfer/postgrest-oauth-server/templates/ ./templates
COPY --from=0 /go/src/github.com/wildsurfer/postgrest-oauth-server/static/ ./static
COPY --from=1 /app/build ./ui
CMD ./postgrest-oauth-server
