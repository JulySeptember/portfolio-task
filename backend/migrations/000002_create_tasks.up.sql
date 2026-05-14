CREATE TABLE tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    user_id BIGINT NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT NULL,

    status ENUM(
        'TODO',
        'DOING',
        'DONE'
    ) NOT NULL DEFAULT 'TODO',

    due_date DATETIME NULL,

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_tasks_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================
-- indexes
-- =========================

-- user task lookup
CREATE INDEX idx_tasks_user_id
    ON tasks(user_id);

-- status filtering
CREATE INDEX idx_tasks_status
    ON tasks(status);

-- user + status filtering
CREATE INDEX idx_tasks_user_status
    ON tasks(user_id, status);

-- due date filtering/sorting
CREATE INDEX idx_tasks_due_date
    ON tasks(due_date);

-- task list ordering
-- supports:
-- WHERE user_id = ?
-- ORDER BY created_at DESC, id DESC
CREATE INDEX idx_tasks_user_created_at
    ON tasks(user_id, created_at DESC, id DESC);