package repo

const queryGetUserByEmail = `
	SELECT id, email, password_hash, status, created_at
	FROM users
	WHERE email = $1 AND deleted_at IS NULL
`

const queryGetUserByID = `
	SELECT id, email, password_hash, status, created_at
	FROM users
	WHERE id = $1 AND deleted_at IS NULL
`

const queryGetUserRoles = `
	SELECT r.name
	FROM roles r
	INNER JOIN user_roles ur ON ur.role_id = r.id
	WHERE ur.user_id = $1
`

const queryCreateRefreshToken = `
	INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
	VALUES (gen_random_uuid(), $1, $2, $3)
	RETURNING id, created_at
`

const queryGetRefreshTokenByHash = `
	SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
	FROM refresh_tokens
	WHERE token_hash = $1
`

const queryRevokeRefreshToken = `
	UPDATE refresh_tokens SET revoked_at = NOW()
	WHERE id = $1 AND revoked_at IS NULL
`

const queryRevokeAllUserRefreshTokens = `
	UPDATE refresh_tokens SET revoked_at = NOW()
	WHERE user_id = $1 AND revoked_at IS NULL
`
