-- Planting timing information. start/end are relative to first/last
-- frost dates, transplant is a positive day count.

alter table seeds add column latin_name text not null default "";

alter table seeds add column needs_stratification integer;

alter table seeds add column sun_type text check (sun_type in ('Full', 'Partial', 'Shade', 'Unknown')) not null default 'Unknown';

alter table seeds add column soil_type text check (soil_type in ('Dry', 'Wet', 'Both', 'Unknown')) not null default 'Unknown';
