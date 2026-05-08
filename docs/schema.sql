CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    email VARCHAR(255) NOT NULL UNIQUE,

    display_name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    user_id BIGINT NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    status ENUM('TODO', 'DOING', 'DONE')
        NOT NULL DEFAULT 'TODO',

    due_date DATETIME NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_tasks_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_tasks_status
    ON tasks(status);

CREATE INDEX idx_tasks_user_id
    ON tasks(user_id);