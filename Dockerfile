FROM golang:1.9.2-alpine3.6
MAINTAINER Ivan Kuznetsov <kuzma.wm@gmail.com>

ENV OAUTH_DB_CONN_STRING="postgres://user:pass@postgresql:5432/test?sslmode=disable" \
    OAUTH_ACCESS_TOKEN_JWT_SECRET="morethan32symbolssecretkey!!!!!!" \
    OAUTH_ACCESS_TOKEN_TTL=7200 \
    OAUTH_REFRESH_TOKEN_JWT_SECRET="notlesshan32symbolssecretkey!!!!" \
    OAUTH_COOKIE_HASH_KEY="supersecret" \
    OAUTH_COOKIE_BLOCK_KEY="16charssecret!!!" \
    OAUTH_TEMPLATE_NAME="index.html" \
    OAUTH_TEMPLATE_PATH="./"

WORKDIR /go/src/github.com/wildsurfer/postgrest-oauth-server
RUN apk add --no-cache openssl git
RUN wget -O /usr/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64
RUN chmod +x /usr/bin/dep
COPY . .
RUN dep ensure -vendor-only && go build
#RUN dep status
RUN apk del openssl git && rm -vf /usr/bin/dep
CMD ./postgrest-oauth-server \
    -dbConnString "${OAUTH_DB_CONN_STRING}" \
    -accessTokenJWTSecret "${OAUTH_ACCESS_TOKEN_JWT_SECRET}" \
    -accessTokenTTL ${OAUTH_ACCESS_TOKEN_TTL} \
    -refreshTokenJWTSecret "${OAUTH_REFRESH_TOKEN_JWT_SECRET}" \
    -cookieBlockKey "${OAUTH_COOKIE_BLOCK_KEY}" \
    -cookieHashKey "${OAUTH_COOKIE_HASH_KEY}" \
    -templateName "${OAUTH_TEMPLATE_NAME}" \
    -templatePath "${OAUTH_TEMPLATE_PATH}"
