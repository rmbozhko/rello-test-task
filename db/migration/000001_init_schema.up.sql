create table if not exists building (
    id serial primary key,
    name text unique not null,
    address text
);

-- building
-- id: Primary key, integer, auto-increment
-- name: String, unique
-- address: Text

create table if not exists apartment (
    id serial primary key,
    number text,
    floor int,
    sq_meters int,
    building_id int not null references building (id) on delete cascade
);

-- apartment
-- id: Primary key, integer, auto-increment
-- building_id: Foreign key, integer (references building.id)
-- number: String
-- floor: Integer
-- sq_meters: Integer