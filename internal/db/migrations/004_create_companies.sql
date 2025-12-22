CREATE TABLE companies (
    company_id SERIAL PRIMARY KEY,
    employer_id INT NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    company_product VARCHAR(255),

    CONSTRAINT fk_company_employer
        FOREIGN KEY (employer_id)
        REFERENCES employers(employer_id)
        ON DELETE CASCADE
);
