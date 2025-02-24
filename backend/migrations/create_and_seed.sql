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
    CHECK (max_members > 0),
    INDEX idx_study_groups_module_id (module_id)
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
