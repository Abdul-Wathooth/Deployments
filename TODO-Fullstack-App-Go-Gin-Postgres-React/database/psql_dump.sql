CREATE table list (
    id serial primary key,
    item char(140),
    done bool,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);