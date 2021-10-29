CREATE TABLE IF NOT EXISTS `invite_tokens`
(
    `token_digest` VARCHAR(255) PRIMARY KEY,
    `expires_at`   TIMESTAMP NOT NULL,
    `community_id` BIGINT    NOT NULL,
    `issued_at`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;
