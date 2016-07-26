DROP TABLE IF EXISTS que_jobs;
CREATE TABLE que_jobs (
	priority    SMALLINT    NOT NULL DEFAULT 100,
	run_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
	job_id      BIGSERIAL   NOT NULL,
	job_class   TEXT        NOT NULL,
	args        JSON        NOT NULL DEFAULT '[]'::json,
	error_count INTEGER     NOT NULL DEFAULT 0,
	last_error  TEXT,
	queue       TEXT        NOT NULL DEFAULT '',
	CONSTRAINT que_jobs_pkey PRIMARY KEY (queue, priority, run_at, job_id)
);
