-- name: GetLocations :many
select * from locations;

-- name: UpdateLocation :exec
update locations set name=?,qr_state=? where id=?;
