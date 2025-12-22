CREATE TABLE employers (
    employer_id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    employer_name VARCHAR(255),
    employer_email VARCHAR(255),

    CONSTRAINT fk_employer_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);
