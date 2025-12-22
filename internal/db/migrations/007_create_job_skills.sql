CREATE TABLE job_skills (
    job_id INT NOT NULL,
    skill_id INT NOT NULL,

    PRIMARY KEY (job_id, skill_id),

    CONSTRAINT fk_job_skill_job
        FOREIGN KEY (job_id)
        REFERENCES jobs(job_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_job_skill_skill
        FOREIGN KEY (skill_id)
        REFERENCES skills(skill_id)
        ON DELETE CASCADE
);
