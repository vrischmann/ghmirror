CREATE TABLE IF NOT EXISTS repository(
    id bigint primary key,
    name varchar,
    local_path varchar,
    clone_url varchar,
    hook_id bigint
);

CREATE TABLE IF NOT EXISTS owner_blacklist(
    id serial primary key,
    name varchar
);

CREATE INDEX owner_blacklist_name_idx ON owner_blacklist(name);

CREATE TABLE IF NOT EXISTS repository_blacklist(
    id serial primary key,
    organization varchar,
    name varchar
);

CREATE INDEX repository_blacklist_idx ON repository_blacklist(organization, name);
