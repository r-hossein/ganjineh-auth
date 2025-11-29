-- name: InsertError :exec
INSERT INTO errors (http_code,status_code, message, stack_trace, endpoint, method, query_params,request_body)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

