CREATE TABLE IF NOT EXISTS repository(
    id bigserial primary key,
    name varchar,
    local_path varchar,
    clone_url varchar,
    hook_id bigint
);
