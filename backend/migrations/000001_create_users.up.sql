CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    -- Cognito sub
    auth_user_id VARCHAR(255) NOT NULL,

    -- login email
    email VARCHAR(255) NOT NULL,

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT uq_users_auth_user_id
        UNIQUE (auth_user_id),

    CONSTRAINT uq_users_email
        UNIQUE (email)
);