package repo

const queryClaimPendingJobs = `
	WITH claimable AS (
		SELECT id
		FROM search_index_jobs
		WHERE status IN ('pending', 'failed')
		  AND scheduled_at <= NOW()
		ORDER BY scheduled_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	)
	UPDATE search_index_jobs j
	SET status = 'processing', updated_at = NOW()
	FROM claimable c
	WHERE j.id = c.id
	RETURNING j.id, j.program_id, j.action, j.status, j.attempts, j.max_attempts,
	          j.last_error, j.scheduled_at, j.processed_at, j.created_at, j.updated_at
`

const queryMarkCompleted = `
	UPDATE search_index_jobs
	SET status = 'completed', processed_at = NOW(), updated_at = NOW()
	WHERE id = $1
`

const queryMarkFailed = `
	UPDATE search_index_jobs
	SET status = 'failed',
	    attempts = attempts + 1,
	    last_error = $2,
	    scheduled_at = $3,
	    updated_at = NOW()
	WHERE id = $1
`

const queryMarkDead = `
	UPDATE search_index_jobs
	SET status = 'dead',
	    attempts = attempts + 1,
	    last_error = $2,
	    updated_at = NOW()
	WHERE id = $1
`

const queryGetProgramForIndex = `
	SELECT p.id,
	       p.title,
	       p.description,
	       p.program_type,
	       p.status,
	       p.duration::TEXT AS duration,
	       p.published_at,
	       c.name AS category,
	       l.code AS language,
	       p.thumbnail,
	       p.video_url,
	       p.created_at
	FROM programs p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN languages l ON l.id = p.language_id
	WHERE p.id = $1 AND p.deleted_at IS NULL
`
