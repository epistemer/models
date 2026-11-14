-- name: GetNote :one
SELECT * FROM note WHERE id = $1;
