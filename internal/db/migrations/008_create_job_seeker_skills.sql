CREATE TABLE job_seeker_skills (
    job_seeker_id INT NOT NULL,
    skill_id INT NOT NULL,

    PRIMARY KEY (job_seeker_id, skill_id),

    CONSTRAINT fk_js_skill_seeker
        FOREIGN KEY (job_seeker_id)
        REFERENCES job_seekers(job_seeker_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_js_skill_skill
        FOREIGN KEY (skill_id)
        REFERENCES skills(skill_id)
        ON DELETE CASCADE
);
