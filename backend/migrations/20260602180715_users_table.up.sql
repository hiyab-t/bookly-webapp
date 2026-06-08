CREATE extension if not EXISTS citext;

create table if not EXISTS users (
    user_id bigserial PRIMARY KEY,
    create_at TIMESTAMP(0) with time zone not null DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hashed bytea not null,
    active bool,
    version int not null DEFAULT 1
)