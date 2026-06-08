create table tokens (
    hash bytea PRIMARY KEY,
    user_id bigint not null REFERENCES users on delete CASCADE,
    expiry TIMESTAMP(0) with time zone NOT NULL,
    scope text NOT NULL
);
