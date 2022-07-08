DROP INDEX IF EXISTS trgm_idx_users_country;

DROP INDEX IF EXISTS idx_users_country;

DROP INDEX IF EXISTS trgm_idx_users_email;

DROP INDEX IF EXISTS email_idx;

DROP INDEX IF EXISTS trgm_idx_users_nickname;

DROP INDEX IF EXISTS idx_users_nickname;

DROP INDEX IF EXISTS trgm_idx_users_last_name;

DROP INDEX IF EXISTS idx_users_last_name;

DROP INDEX IF EXISTS trgm_idx_users_first_name;

DROP INDEX IF EXISTS idx_users_first_name;

DROP TABLE IF EXISTS users;

DROP DOMAIN IF EXISTS email;

DROP EXTENSION IF EXISTS pg_trgm;

DROP EXTENSION IF EXISTS citext;

DROP EXTENSION IF EXISTS "uuid-ossp";