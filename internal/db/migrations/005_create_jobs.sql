CREATE TABLE jobs (
    job_id SERIAL PRIMARY KEY,
    company_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    job_location VARCHAR(255),
    job_type VARCHAR(50),
    salary_min INT,
    salary_max INT,
    status VARCHAR(50) DEFAULT 'published',

    CONSTRAINT fk_job_company
        FOREIGN KEY (company_id)
        REFERENCES companies(company_id)
        ON DELETE CASCADE
);
