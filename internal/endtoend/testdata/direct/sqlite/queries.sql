-- name: Insert :exec
INSERT INTO records (id, value, state) VALUES (?1, ?2, 7);

-- name: Get :one
SELECT value, note, state FROM records WHERE id = ?1;
