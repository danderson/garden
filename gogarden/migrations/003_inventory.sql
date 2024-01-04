create table location (
       id integer primary key,
       name text not null,
       qr_state integer not null
);

create table planted (
       id integer primary key,
       name text not null,
       location integer not null references location(id) deferrable initially deferred,
       planted integer not null,
       removed integer not null
);
