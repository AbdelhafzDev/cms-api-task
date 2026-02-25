package repo

const queryListSources = `
	SELECT id, name, source_type, base_url, is_active, created_at, updated_at
	FROM import_sources
	ORDER BY id ASC
`

const queryGetSourceByID = `
	SELECT id, name, source_type, base_url, is_active, created_at, updated_at
	FROM import_sources
	WHERE id = $1
`

const queryCreateLog = `
	INSERT INTO import_logs (id, source_id, triggered_by, status, records_imported, error_message, started_at, finished_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
`

const queryUpdateLog = `
	UPDATE import_logs
	SET status = $1,
	    records_imported = $2,
	    error_message = $3,
	    finished_at = $4
	WHERE id = $5
`
