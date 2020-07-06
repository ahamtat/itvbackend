CREATE TABLE requests (
    id bigserial not null primary key,
    uuid uuid,
    method varchar not null,
    url varchar not null,
    fetch_headers varchar,
    body varchar,
    status integer,
    response_headers varchar,
    length integer
);