create table if not exists permissions (
    permission_id bigserial PRIMARY KEY,
    code text NOT NULL
);

create table if not EXISTS user_permission (
    user_id bigint not null REFERENCES users on delete CASCADE,
    permission_id bigint not null REFERENCES permissions on delete CASCADE,
    PRIMARY KEY(user_id, permission_id)
);

insert into permissions (code) VALUES ('books:read'),('books:write');