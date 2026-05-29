CREATE TABLE tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    -- external/public identifier
    public_id CHAR(36) NOT NULL,

    user_id BIGINT NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT NULL,

    status VARCHAR(16)
        NOT NULL
        DEFAULT 'TODO',

    due_date DATETIME NULL,

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT uq_tasks_public_id
        UNIQUE (public_id),

    CONSTRAINT fk_tasks_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================
-- indexes
-- =========================

CREATE INDEX idx_tasks_user_id
    ON tasks(user_id);

CREATE INDEX idx_tasks_user_created_at
    ON tasks(
        user_id,
        created_at DESC,
        id DESC
    );

CREATE INDEX idx_tasks_user_status_created
    ON tasks(
        user_id,
        status,
        created_at DESC,
        id DESC
    );

CREATE INDEX idx_tasks_user_due_date
    ON tasks(
        user_id,
        due_date,
        id DESC
    );