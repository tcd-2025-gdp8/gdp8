CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS modules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS study_groups (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type ENUM('public', 'closed', 'invite-only') NOT NULL,
    module_id INT NOT NULL,
    max_members INT NOT NULL,

    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE RESTRICT,

    INDEX idx_study_groups_module_id (module_id),

    CHECK (max_members > 0)
);

CREATE TABLE IF NOT EXISTS user_modules (
    user_id VARCHAR(255) NOT NULL,
    module_id INT NOT NULL,

    PRIMARY KEY (user_id, module_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE,

    INDEX idx_user_modules_user_id (user_id),
    INDEX idx_user_modules_module_id (module_id)
);

CREATE TABLE IF NOT EXISTS user_study_groups (
    user_id VARCHAR(255) NOT NULL,
    study_group_id INT NOT NULL,
    type ENUM('admin', 'member', 'invitee', 'requester') NOT NULL,

    PRIMARY KEY (user_id, study_group_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE CASCADE,

    INDEX idx_user_study_groups_user_id (user_id),
    INDEX idx_user_study_groups_study_group_id (study_group_id)
);

INSERT IGNORE INTO modules (id, code, name)
VALUES
  (1, 'CSU44052', 'Computer Graphics'),
  (2, 'CSU44061', 'Machine Learning'),
  (3, 'CSU44051', 'Human Factors'),
  (4, 'CSU44000', 'Internet Applications'),
  (5, 'CSU44012', 'Topics in Functional Programming'),
  (6, 'CSU44099', 'Final Year Project'),
  (7, 'CSU44098', 'Group Design Project'),
  (8, 'CSU44081', 'Entrepreneurship & High Tech Venture Creation');

CREATE TABLE IF NOT EXISTS notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type ENUM(
        'study-group-joined',
        'study-group-requested-to-join',
        'study-group-left',
        'study-group-accepted-invite',
        'study-group-rejected-invite',
        'study-group-invited',
        'study-group-accepted-join-request',
        'study-group-rejected-join-request',
        'study-group-removed-member',
        'study-group-chat-message'
        ) NOT NULL,
    triggering_user_id VARCHAR(255) NOT NULL,
    target_user_id VARCHAR(255) NULL,
    study_group_id INT NULL NOT NULL,
    message_id INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (triggering_user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (target_user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE CASCADE,

    CHECK (
        CASE
            WHEN type = 'study-group-joined'                THEN (target_user_id IS NULL AND message_id IS NULL)
            WHEN type = 'study-group-requested-to-join'     THEN (target_user_id IS NULL AND message_id IS NULL)
            WHEN type = 'study-group-left'                  THEN (target_user_id IS NULL AND message_id IS NULL)
            WHEN type = 'study-group-accepted-invite'       THEN (target_user_id IS NULL AND message_id IS NULL)
            WHEN type = 'study-group-rejected-invite'       THEN (target_user_id IS NULL AND message_id IS NULL)

            WHEN type = 'study-group-chat-message'          THEN (target_user_id IS NULL AND message_id IS NOT NULL)

            WHEN type = 'study-group-invited'               THEN (target_user_id IS NOT NULL AND message_id IS NULL)
            WHEN type = 'study-group-accepted-join-request' THEN (target_user_id IS NOT NULL AND message_id IS NULL)
            WHEN type = 'study-group-rejected-join-request' THEN (target_user_id IS NOT NULL AND message_id IS NULL)
            WHEN type = 'study-group-removed-member'        THEN (target_user_id IS NOT NULL AND message_id IS NULL)

            ELSE FALSE
        END = TRUE
    )
);

CREATE TABLE IF NOT EXISTS user_notifications (
    user_id VARCHAR(255) NOT NULL,
    notification_id INT NOT NULL,

    PRIMARY KEY (user_id, notification_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE,

    INDEX idx_user_notifications_user_id (user_id),
    INDEX idx_user_notifications_notification_id (notification_id)
);

CREATE TRIGGER IF NOT EXISTS delete_notification_when_no_corresponding_users
    AFTER DELETE ON user_notifications
    FOR EACH ROW
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM user_notifications WHERE notification_id = OLD.notification_id
    ) THEN
        DELETE FROM notifications WHERE id = OLD.notification_id;
    END IF;
END;

CREATE TABLE IF NOT EXISTS study_group_chat_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    study_group_id INT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,

    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,

    INDEX idx_study_group_chat_messages_timestamp (timestamp),
    INDEX idx_study_group_chat_messages_study_group_id_and_timestamp (study_group_id, timestamp)
);

CREATE TABLE IF NOT EXISTS study_session_availability_requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    study_group_id INT NOT NULL,
    creator_id VARCHAR(255) NOT NULL,
    availability_period_start TIMESTAMP NOT NULL,
    availability_period_end TIMESTAMP NOT NULL,

    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE CASCADE,

    INDEX idx_study_session_availability_requests_period_end (availability_period_end),
    INDEX idx_study_session_availability_requests_group_and_period_end (study_group_id, availability_period_end)
);

CREATE PROCEDURE IF NOT EXISTS delete_expired_study_session_availability_requests()
BEGIN
    DELETE FROM study_session_availability_requests
    WHERE availability_period_end < NOW();
END;

CREATE EVENT IF NOT EXISTS delete_expired_study_session_availability_requests_event
    ON SCHEDULE EVERY 1 DAY
    DO CALL delete_expired_study_session_availability_requests();

CALL delete_expired_study_session_availability_requests();

CREATE VIEW IF NOT EXISTS current_study_session_availability_requests AS
    SELECT *
    FROM study_session_availability_requests
    WHERE availability_period_end > NOW();

CREATE TABLE IF NOT EXISTS study_session_availability_entries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    availability_request_id INT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    availability_entry_start TIMESTAMP NOT NULL,
    availability_entry_end TIMESTAMP NOT NULL,

    FOREIGN KEY (availability_request_id) REFERENCES study_session_availability_requests(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_study_session_availability_entries_availability_request_id (availability_request_id)
);

CREATE TABLE IF NOT EXISTS study_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    study_group_id INT NOT NULL,
    creator_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    duration_minutes INT NOT NULL,

    end_time TIMESTAMP GENERATED ALWAYS AS (start_time + INTERVAL duration_minutes MINUTE) STORED,

    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE RESTRICT,

    INDEX idx_study_sessions_study_group_id (study_group_id),
    INDEX idx_study_sessions_group_id_and_end_time (study_group_id, end_time)
);

CREATE TABLE IF NOT EXISTS files (
    name VARCHAR(255) NOT NULL,
    file_owner_id VARCHAR(255) NOT NULL,
    study_group_id INT NOT NULL,
    PRIMARY KEY (name, study_group_id),
    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (file_owner_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS chatbot_memory (
    study_group_id INT NOT NULL,
    context_src VARCHAR(255) NOT NULL,
    context_data LONGTEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (study_group_id, context_src),
    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (context_src, study_group_id) REFERENCES files(name, study_group_id) ON DELETE CASCADE
);
