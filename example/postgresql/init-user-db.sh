#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 \
    --username "$POSTGRES_USER" \
    --dbname "$POSTGRES_DB" \
    --host localhost <<-EOSQL

CREATE EXTENSION pgcrypto;

CREATE SCHEMA IF NOT EXISTS oauth2;
CREATE TABLE IF NOT EXISTS
    oauth2.owners (
      id                  serial PRIMARY KEY NOT NULL,
      email               text DEFAULT NULL UNIQUE CHECK ( email ~* '^.+@.+\..+$' ),
      phone               text DEFAULT NULL UNIQUE,
      password            text NOT NULL DEFAULT md5(random()::text) CHECK (length(password) < 512),
      role                varchar NOT NULL DEFAULT 'member',
      jti                 timestamp without time zone NOT NULL DEFAULT now(),
      CHECK(email IS NOT NULL OR phone IS NOT NULL)
    );

CREATE OR REPLACE FUNCTION oauth2.create_owner(email text, phone text, password text, verification_code text, verification_route text, OUT id varchar, OUT role varchar, OUT jti varchar)
AS \$\$
        INSERT INTO oauth2.owners(email, phone, password) VALUES (NULLIF(email, ''), NULLIF(phone, ''), crypt(password, gen_salt('bf'))) RETURNING id::varchar, role, jti::varchar;
\$\$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.check_owner(username text, password text, OUT id varchar, OUT role varchar, OUT jti varchar)
AS \$\$
SELECT id::varchar, role::varchar, jti::varchar FROM oauth2.owners
    WHERE (email = check_owner.username OR phone = check_owner.username)
        AND owners.password = crypt(check_owner.password, owners.password);
\$\$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.owner_role_and_jti_by_id(id text, OUT role varchar, OUT jti varchar)
AS \$\$
SELECT role::varchar, jti::varchar FROM oauth2.owners
    WHERE (id = owner_role_and_jti_by_id.id::bigint);
\$\$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.verify_owner(user_id varchar) RETURNS void
AS \$\$
UPDATE oauth2.owners SET role='verified' WHERE oauth2.owners.id = user_id::int;
\$\$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.password_request(username text, verification_code text, verification_route text, OUT id varchar)
AS \$\$
        SELECT id::varchar from oauth2.owners WHERE email = password_request.username OR phone = password_request.username;
\$\$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.password_reset(id text, password text) RETURNS void
AS \$\$
        UPDATE oauth2.owners SET password = crypt(password_reset.password, gen_salt('bf')), jti = now() WHERE id = password_reset.id::int;
\$\$ LANGUAGE SQL;

CREATE TABLE IF NOT EXISTS
    oauth2.clients (
      id                  text NOT NULL PRIMARY KEY,
      secret              text DEFAULT gen_random_uuid()::text,
      redirect_uri        text DEFAULT NULL UNIQUE,
      type                varchar NOT NULL DEFAULT 'public'
    );

INSERT INTO oauth2.clients(id, redirect_uri, type) VALUES('mobile', 'https://mobile.uri', 'public');
INSERT INTO oauth2.clients(id, redirect_uri, type) VALUES('spa', 'https://spa.uri', 'public');
INSERT INTO oauth2.clients(id, secret, type) VALUES('worker', 'secret', 'confidential');

CREATE OR REPLACE FUNCTION oauth2.check_client(client_id text, OUT redirect_uri text)
AS \$\$
SELECT redirect_uri FROM oauth2.clients
    WHERE id = check_client.client_id;
\$\$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.check_client_secret(client_id text, client_secret text, OUT type varchar)
AS \$\$
SELECT type FROM oauth2.clients
    WHERE id = check_client_secret.client_id AND secret = check_client_secret.client_secret;
\$\$ LANGUAGE SQL;

CREATE SCHEMA api;
CREATE OR REPLACE VIEW api.me AS
 SELECT
    id,
    email,
    phone,
    role
 FROM oauth2.owners WHERE
    oauth2.owners.id = current_setting('request.jwt.claim.id', true)::int
 WITH LOCAL CHECK OPTION;

CREATE OR REPLACE VIEW api.client AS
 SELECT
    id,
    type
 FROM oauth2.clients WHERE
    oauth2.clients.id = current_setting('request.jwt.claim.client_id', true)::varchar
 WITH LOCAL CHECK OPTION;

-------------
--  Roles  --
-------------

CREATE ROLE authenticator NOINHERIT LOGIN PASSWORD '$PGRST_AUTHENTICATOR_PASSWORD';

CREATE ROLE "guest" NOLOGIN;
GRANT "guest" TO "authenticator";

CREATE ROLE "member" NOLOGIN;
GRANT "member" TO "authenticator";
GRANT USAGE ON SCHEMA api TO "member";
GRANT SELECT ON TABLE api.me TO "member";

CREATE ROLE "msrv-worker" NOLOGIN;
GRANT "msrv-worker" TO "authenticator";
GRANT USAGE ON SCHEMA api TO "msrv-worker";
GRANT SELECT ON TABLE api.client TO "msrv-worker";

EOSQL
