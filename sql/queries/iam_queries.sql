-- name: GetUserByPhone :one
SELECT 
    u.id,
    u.status,
    CASE 
        WHEN u.status = 'active' THEN u.phone_number
        ELSE NULL 
    END AS phone_number,
    CASE 
        WHEN u.status = 'active' THEN u.first_name 
        ELSE NULL 
    END AS first_name,
    CASE 
        WHEN u.status = 'active' THEN u.last_name 
        ELSE NULL 
    END AS last_name,
    CASE 
        WHEN u.status = 'active' THEN r.name 
        ELSE NULL 
    END AS main_role_name,
    CASE 
        WHEN u.status = 'active' THEN 
            COALESCE(
                json_agg(
                    json_build_object(
                        'company_id', uc.company_id,
                        'role_name', r2.name
                    )
                ) FILTER (WHERE uc.id IS NOT NULL),
                '[]'
            )
        ELSE '[]'::json
    END AS company_roles
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
LEFT JOIN user_companies uc ON uc.user_id = u.id AND uc.is_active = true
LEFT JOIN companies c ON uc.company_id = c.id
LEFT JOIN roles r2 ON uc.role_id = r2.id
WHERE u.phone_number = $1
GROUP BY u.id, u.status, r.name;

-- name: CreateUser :one
INSERT INTO users (phone_number, first_name, last_name, role_id, is_phone_verified)
VALUES ($1, $2, $3, $4, true)
ON CONFLICT (phone_number) DO NOTHING
RETURNING id, phone_number, first_name, last_name, role_id;

-- name: UpdateUserPhone :exec
UPDATE users 
SET 
    phone_number = $2
WHERE id = $1;

-- name: SetUserPassword :exec
UPDATE users 
SET 
    password_hash = $2,
    is_password_set = TRUE
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users 
SET 
    status = 'inactive'
WHERE id = $1;

-- name: SuspendUser :exec
UPDATE users 
SET 
    status = 'suspended'
WHERE id = $1;

-- name: Activeuser :exec
UPDATE users 
SET 
    status = 'active'
WHERE id = $1;

-- name: InsertUserSession :exec
INSERT INTO user_sessions (
    id,
    user_id,
    refresh_token_hash,
    refresh_token_created_at,
    refresh_token_expires_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetSession :one
SELECT 
    id,
    user_id,
    refresh_token_hash,
    refresh_token_expires_at,
    status
FROM user_sessions 
WHERE id = $1 AND user_id = $2;

-- name: UpdateSession :exec
UPDATE user_sessions 
SET 
    refresh_token_hash = $2,
    refresh_token_created_at = $3,
    refresh_token_expires_at = $4,
    last_active = $5
WHERE id = $1;

-- name: GetUserActiveSessions :many
SELECT id 
FROM user_sessions 
WHERE user_id = $1 AND status = 'active';

-- name: RevokeAllUserSessions :exec
UPDATE user_sessions 
SET 
    status = 'revoked',
    revoked_at = NOW(),
    revoked_by = $2
WHERE user_id = $1 AND status = 'active';

-- name: RevokeSession :exec
UPDATE user_sessions 
SET 
    status = 'revoked',
    revoked_at = NOW(),
    revoked_by = $3
WHERE id = $1 AND user_id = $2;

-- name: UpdatedAllUserSessions :exec
UPDATE user_sessions 
SET 
    status = 'updated'
WHERE user_id = $1 AND status = 'active';

-- name: GetAllRoles :many
SELECT
    name,
    permission_codes,
    is_system_role
FROM roles
ORDER BY id;