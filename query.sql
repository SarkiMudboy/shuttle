-- name: GetRequest :one
SELECT * FROM request_history
WHERE request_id = ? LIMIT 1;

-- name: GetRequestHistory :many
SELECT * FROM request_history
ORDER BY request_time DESC LIMIT 20;

-- name: CreateRequest :execresult
INSERT INTO request_history (
  endpoint, headers, body, method  
) VALUES (?, ?, ?, ?);

-- name: DeleteRequest :exec
DELETE FROM request_history
WHERE request_id = ?;


-- name: GetlastRequest :one
SELECT * FROM request_history
ORDER BY request_time DESC LIMIT 1;
