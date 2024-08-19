-- name: ListMemos :many
SELECT id, user_id, content, created_at FROM memos WHERE user_id = $1;

-- name: GetMemo :one
SELECT id, user_id, content, created_at FROM memos WHERE user_id = $1 AND id = $2;

-- name: CreateMemo :exec
INSERT INTO memos (id, user_id, content, created_at) VALUES ($1, $2, $3, $4);

-- name: UpdateMemo :exec
UPDATE memos SET content = $2 WHERE id = $1;