package repo

const queryGetByID = `
	SELECT p.id, p.title, p.description, p.program_type, p.duration,
	       p.published_at, p.thumbnail, p.video_url, p.status,
	       p.category_id, p.language_id, p.created_at, p.updated_at,
	       c.name AS category_name,
	       l.code AS language_code
	FROM programs p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN languages l ON l.id = p.language_id
	WHERE p.id = $1 AND p.status = 'active' AND p.deleted_at IS NULL
`

const queryListFirst = `
	SELECT p.id, p.title, p.description, p.program_type, p.duration,
	       p.published_at, p.thumbnail, p.video_url, p.status,
	       p.category_id, p.language_id, p.created_at, p.updated_at,
	       c.name AS category_name,
	       l.code AS language_code
	FROM programs p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN languages l ON l.id = p.language_id
	WHERE p.status = 'active' AND p.deleted_at IS NULL
	ORDER BY p.published_at DESC, p.id DESC
	LIMIT $1
`

const queryListAfterCursor = `
	SELECT p.id, p.title, p.description, p.program_type, p.duration,
	       p.published_at, p.thumbnail, p.video_url, p.status,
	       p.category_id, p.language_id, p.created_at, p.updated_at,
	       c.name AS category_name,
	       l.code AS language_code
	FROM programs p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN languages l ON l.id = p.language_id
	WHERE p.status = 'active' AND p.deleted_at IS NULL
	  AND (p.published_at, p.id) < ($2, $3)
	ORDER BY p.published_at DESC, p.id DESC
	LIMIT $1
`
