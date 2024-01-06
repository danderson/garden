-- name: ListLocations :many
select * from locations order by name;

-- name: UpdateLocation :exec
update locations set name=?,qr_state=? where id=?;

-- name: ListSeeds :many
select * from seeds order by name;
