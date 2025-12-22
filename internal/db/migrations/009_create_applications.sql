CREATE TABLE applications (
    application_id SERIAL PRIMARY KEY,
    job_id INT NOT NULL,
    job_seeker_id INT NOT NULL,
    cover_letter TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    applied_at TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_application_job
        FOREIGN KEY (job_id)
        REFERENCES jobs(job_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_application_seeker
        FOREIGN KEY (job_seeker_id)
        REFERENCES job_seekers(job_seeker_id)
        ON DELETE CASCADE,

    CONSTRAINT unique_job_application
        UNIQUE (job_id, job_seeker_id)
);
