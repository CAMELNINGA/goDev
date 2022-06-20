CREATE TABLE IF NOT EXISTS path(
    id serial primary key,
    display_name text,
    deleted bool default false ,
    created_at timestamp default now()
);


CREATE TABLE  file(
    id serial primary key ,
    user_id integer references users(id),
    path_id integer references path(id),
    paths text not null
);

CREATE TABLE users_paths(
    id serial primary key,
    path_id integer references path(id),
    user_id integer references users(id),
    create_at timestamp default now(),
    deleted bool default false ,
    delete_at timestamp
);

ALTER TABLE users ADD COLUMN path_id integer references path(id);