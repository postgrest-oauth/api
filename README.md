README
------

[![Build Status](https://travis-ci.org/postgrest-oauth/api.svg?branch=master)](https://travis-ci.org/postgrest-oauth/api)

Environment Variables
=====================

**OAUTH_DB_CONN_STRING**

Default: "postgres://user:pass@postgresql:5432/test?sslmode=disable"

See http://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING for more information about connection string parameters.

**OAUTH_ACCESS_TOKEN_JaW_SECRET**

Default: "morethan32symbolssecretkey!!!!!!"

Random string. Should be >= to 32 symbols. This is important.

**OAUTH_ACCESS_TOKEN_TTL=7200**

Default: 7200

Access token life cycle in seconds

**OAUTH_REFRESH_TOKEN_SECRET**

Default: "notlesshan32symbolssecretkey!!!!"

Random string. Should be >= to 32 symbols. This is important.

**OAUTH_COOKIE_HASH_KEY**

Default: "supersecret"

Random string.

**OAUTH_COOKIE_BLOCK_KEY**

Default: "16charssecret!!!"

Random string. Should be equal to 16, 24 or 32 symbols. This is important.


**OAUTH_VALIDATE_REDIRECT_URI**

Default: true

This setting should be `true` when you use this in production. When set to `false` you can use any **redirect_uri**. Handy for development. 

**OAUTH_CODE_UI**

Default: http://localhost:3685

This is a URL of UI that is used for Authorization Code Flow. 

Testing with Newman
===================
```
$ cd ./example
$ docker-compose -f docker-compose.yml up -d --build
$ cd ../
$ newman run --bail --ignore-redirects --global-var host=localhost ./postgrest-oauth-server.postman_collection.json

```
