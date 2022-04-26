CREATE TABLE schema_migrations
(
    version bigint  not null
        primary key,
    dirty   boolean not null
);

CREATE TABLE users(
    ID SERIAL PRIMARY KEY,
    TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    USERNAME TEXT,
    CHAT_ID INT,
    MESSAGE TEXT,
    ANSWER TEXT
                  );

CREATE TABLE app_log
(
    user_id  integer references users (id),
    start_dt timestamp default now(),
    header   text,
    body     text,
    status   int
);
