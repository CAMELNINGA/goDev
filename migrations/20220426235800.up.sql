ALTER TABLE users DROP column message;
ALTER TABLE users DROP column answer;


CREATE TABLE shops(
                      ID SERIAL PRIMARY KEY,
                    display_name text ,
                    description text,
                    photo text,
                    created_at timestamp default now()
);

CREATE TABLE shops_count(
  ID serial primary key ,
  size text ,
  count integer,
  shop_id integer references shops(id),
  created_at timestamp default now()
);

CREATE TABLE dc_order(
    id serial primary key ,
    display_name text
);

CREATE TABLE orders(
    ID serial primary key ,
    user_id integer references users(id),
    shop_count_id integer references shops_count(id),
    address text ,
    status integer references dc_order(id),
    created_at timestamp default now()
);

CREATE TABLE ds_shops_type(
                           ID SERIAL primary key ,
                           display_name text
);

ALTER TABLE ds_shops_type ADD COLUMN dc_sh_type_id integer references ds_shops_type(id);

CREATE TABLE shops_type(
  ID SERIAL primary key ,
  dc_sh_type_id integer references ds_shops_type(id),
  shop_id integer references shops(id)
)