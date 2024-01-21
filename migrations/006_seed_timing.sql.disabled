-- Planting timing information. start/end are relative to first/last
-- frost dates, transplant is a positive day count.

create table seed_windows (
  id integer primary key autoincrement,
  seed_id integer not null references seed (id) on delete cascade,
  how text check(how in ('I', 'C', 'D')) not null,
  datum text check(datum in ('F', 'L')) not null,
  start integer not null,
  end integer not null,
  transplant integer,
  check ((how='I' and transplant is not null) or how!='I')
);
