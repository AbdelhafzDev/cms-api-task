package repo

const queryCreate = `
	INSERT INTO programs (id, title, description, program_type, duration, thumbnail, video_url, status, category_id, language_id, created_by, updated_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
`

const queryUpdate = `
	UPDATE programs
	SET title = $1, description = $2, program_type = $3, duration = $4,
	    thumbnail = $5, video_url = $6, status = $7, category_id = $8, language_id = $9,
	    updated_by = $10, updated_at = NOW()
	WHERE id = $11
`

const queryDelete = `
	UPDATE programs SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
`

const queryGetByID = `
	SELECT p.id, p.title, p.description, p.program_type, p.duration,
	       p.published_at, p.thumbnail, p.video_url, p.external_id, p.status,
	       p.category_id, p.language_id, p.import_source_id,
	       p.created_by, p.updated_by, p.created_at, p.updated_at,
	       c.name AS category_name,
	       l.code AS language_code
	FROM programs p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN languages l ON l.id = p.language_id
	WHERE p.id = $1 AND p.deleted_at IS NULL
`

const queryListFirst = `
	SELECT p.id, p.title, p.description, p.program_type, p.duration,
	       p.published_at, p.thumbnail, p.video_url, p.external_id, p.status,
	       p.category_id, p.language_id, p.import_source_id,
	       p.created_by, p.updated_by, p.created_at, p.updated_at,
	       c.name AS category_name,
	       l.code AS language_code
	FROM programs p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN languages l ON l.id = p.language_id
	WHERE p.deleted_at IS NULL
	ORDER BY p.created_at DESC, p.id DESC
	LIMIT $1
`

const queryListAfterCursor = `
	SELECT p.id, p.title, p.description, p.program_type, p.duration,
	       p.published_at, p.thumbnail, p.video_url, p.external_id, p.status,
	       p.category_id, p.language_id, p.import_source_id,
	       p.created_by, p.updated_by, p.created_at, p.updated_at,
	       c.name AS category_name,
	       l.code AS language_code
	FROM programs p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN languages l ON l.id = p.language_id
	WHERE p.deleted_at IS NULL AND (p.created_at, p.id) < ($2, $3)
	ORDER BY p.created_at DESC, p.id DESC
	LIMIT $1
`
