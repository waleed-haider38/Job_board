CREATE TABLE job_seekers (
    job_seeker_id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    full_name VARCHAR(255),
    resume_url TEXT,

    CONSTRAINT fk_job_seeker_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);
