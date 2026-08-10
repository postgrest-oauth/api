PostgREST OAuth2 server
=======================

> [!IMPORTANT]
> **This project is archived and is not maintained.** Do not use it for new projects.
> For authentication in front of PostgREST, see the [PostgREST external authentication docs](https://docs.postgrest.org/en/stable/explanations/external_auth.html),
> or use a maintained auth server such as [Supabase Auth](https://github.com/supabase/auth) or [Keycloak](https://www.keycloak.org/).

A standalone OAuth2 server for [PostgREST](https://postgrest.org), written in Go. It adds user accounts to a PostgREST application: signup, signin, verification by code, password reset and Facebook login. It issues JWT access tokens signed with the same secret that PostgREST validates, so PostgREST accepts them directly and switches the database role to the `role` claim from the token.

The server has no users table of its own. User storage stays in your database: the server calls a small contract of SQL functions in the `oauth2` schema (`create_owner`, `check_owner`, `verify_owner`, `password_request` and so on), and your implementation of these functions decides how users are stored. A working implementation is in the [example](https://github.com/postgrest-oauth/example) repository.

Features
--------

* OAuth2 Authorization Code and Client Credentials flows
* Signup and signin with email or phone
* Account verification with a one-time code, with re-verification
* Password reset flow
* Facebook signup and signin
* JWT access and refresh tokens compatible with PostgREST
* Optional [Hasura](https://hasura.io)-compatible claims in the token
* Configured only with environment variables, ships as a Docker image
* Integration tests with a Postman collection and Newman
* A ready React frontend for the whole flow: [postgrest-oauth/ui](https://github.com/postgrest-oauth/ui)

Endpoints
---------

```
GET  /authorize            Authorization Code flow entry
POST /token                exchange code / refresh token / client credentials
POST /signup               create an account
POST /signin               password signin
POST /verify               confirm account with a code
POST /re-verify            request a new verification code
POST /password/request     start password reset
POST /password/reset       finish password reset
GET  /logout               drop the session
GET  /facebook/url         get Facebook login URL
POST /facebook/enter       signup or signin with Facebook
```

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

**OAUTH_CORS_ALLOW_ORIGIN**

Default: http://localhost:3685,http://localhost:3001

Allowed CORS origins

**HASURA_ALLOWED_ROLES**

Example: "editor,user"

If specified, support for Hasura will be enabled and Hasura specific info will be added to the token:

```
  "https://hasura.io/jwt/claims": {
    "x-hasura-allowed-roles": ["editor","user"],
    "x-hasura-default-role": "user",
    "x-hasura-user-id": "123"
  }
```

More info: https://hasura.io/docs/1.0/graphql/manual/auth/authentication/jwt.html

Facebook Signup/Signin
======================

Prepare
-------

1. Go to [developers.facebook.com](https://developers.facebook.com) and create an app, add Facebook Login product ([tutorial](https://youtu.be/MpLCBEdhg3Y))
2. Add OAUTH_FACEBOOK_CLIENT_ID and OAUTH_FACEBOOK_CLIENT_SECRET environmental variables

Configure your app
---------------

Add 2 functions to your database
```SQL
CREATE OR REPLACE FUNCTION oauth2.create_facebook_owner(obj json, phone varchar, OUT id varchar, OUT role varchar, OUT jti varchar)
AS $$
        INSERT INTO api.users(email, phone, role, facebook_id)
        VALUES
         (
         obj->>'email'::varchar,
         phone,
         'verified',
         obj->>'id'::varchar
         )
        RETURNING id::varchar, role::varchar, jti::varchar;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.check_owner_facebook(facebook_id varchar, OUT id varchar, OUT role varchar, OUT jti varchar)
AS $$
SELECT id::varchar, role::varchar, jti::varchar FROM api.users
    WHERE facebook_id = check_owner_facebook.facebook_id;
$$ LANGUAGE SQL;
```

Get facebook button URL
```
GET http://localhost:3684/facebook/url?redirect_uri=http://localhost:3685/
```

After user clicks it he'll be returned to your app with `code` and `state`. Pass them to `/api/enter` route

```
POST http://localhost:3684/facebook/enter
Content-Type: application/x-www-form-urlencoded

code={CODE}&state={STATE}

```

If user don't exist it will be created. If it exists he'll be signed in. Now you can redirect your app to `/authorize`  

Testing with Newman
===================
```
$ cd ./example
$ docker-compose -f docker-compose.yml up -d --build
$ cd ../
$ newman run --bail --ignore-redirects --global-var host=localhost ./postgrest-oauth-server.postman_collection.json

```
