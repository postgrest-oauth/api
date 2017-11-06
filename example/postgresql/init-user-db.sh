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
      email               text CHECK ( email ~* '^.+@.+\..+$' ),
      phone               varchar(25),
      password            text NOT NULL DEFAULT md5(random()::text) CHECK (length(password) < 512),
      role                varchar NOT NULL DEFAULT 'member',
      unique              (email, phone)
    );

CREATE OR REPLACE FUNCTION oauth2.create_owner(email text, phone text, password text, OUT id varchar, OUT role varchar)
AS \$\$
        INSERT INTO oauth2.owners(email, phone, password) VALUES (email, phone, crypt(password, gen_salt('bf'))) RETURNING id::varchar, role;
\$\$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION oauth2.check_owner(username text, password text, OUT id varchar, OUT role varchar)
AS \$\$
SELECT id::varchar, role::varchar FROM oauth2.owners
    WHERE (email = check_owner.username OR phone = check_owner.username)
        AND owners.password = crypt(check_owner.password, owners.password);
\$\$ LANGUAGE SQL;

CREATE TABLE IF NOT EXISTS
    oauth2.clients (
      id                  text NOT NULL PRIMARY KEY,
      secret              UUID NOT NULL DEFAULT gen_random_uuid(),
      redirect_uri        text NOT NULL UNIQUE
    );

INSERT INTO oauth2.clients(id, redirect_uri) VALUES('mobile', 'https://mobile.uri');
INSERT INTO oauth2.clients(id, redirect_uri) VALUES('spa', 'https://spa.uri');

CREATE OR REPLACE FUNCTION oauth2.check_client(client_id text, client_secret text, OUT redirect_uri text)
AS \$\$
SELECT redirect_uri FROM oauth2.clients
    WHERE id = check_client.client_id;
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

EOSQL

