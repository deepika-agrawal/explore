CREATE TABLE user_decisions (
    recipient_user_id VARCHAR(255) NOT NULL,
    actor_user_id VARCHAR(255) NOT NULL,
    liked_recipient BOOLEAN NOT NULL,
    decision_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (recipient_user_id, actor_user_id)
);

CREATE INDEX idx_recipient_user_id ON user_decisions (recipient_user_id);
CREATE INDEX idx_actor_user_id ON user_decisions (actor_user_id);

INSERT INTO user_decisions (recipient_user_id, actor_user_id, liked_recipient, decision_timestamp) VALUES
(1, 2, TRUE, NOW() - INTERVAL '2 DAYS'),
(2, 1, FALSE, NOW() - INTERVAL '1 DAY'),
(3, 1, TRUE, NOW() - INTERVAL '1 DAY'),
(1, 3, TRUE, NOW()),
(4, 1, TRUE, NOW());