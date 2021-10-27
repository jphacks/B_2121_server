CREATE TABLE IF NOT EXISTS `communities`
(
    `id`          BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`        VARCHAR(255) NOT NULL,
    `description` VARCHAR(511) NOT NULL,
    `latitude`    DOUBLE CHECK ( `latitude` >= -180 AND `latitude` < 180 ),
    `longitude`   DOUBLE CHECK ( `longitude` >= -90 AND `longitude` < 90 ),
    `image_file`  VARCHAR(255) NOT NULL COMMENT '画像のファイル名',
    `created_at`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `restaurants`
(
    `id`         BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`       VARCHAR(255) NOT NULL,
    `latitude`   DOUBLE CHECK ( `latitude` >= -180 AND `latitude` < 180 ),
    `longitude`  DOUBLE CHECK ( `longitude` >= -90 AND `longitude` < 90 ),
    `address`    VARCHAR(255) NOT NULL UNIQUE COMMENT '住所',
    `url`        VARCHAR(255) NOT NULL COMMENT 'レストラン情報の取得元のレストランのURL',
    `image_url`  VARCHAR(255) COMMENT '画像のURL',
    `source`     VARCHAR(255) NOT NULL COMMENT 'レストラン情報の取得元',
    `source_id`  VARCHAR(255) NOT NULL COMMENT 'レストラン情報の取得元のレストランID',
    `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE `u_source_source_id` (source, source_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `communities_restaurants`
(
    `id`            BIGINT AUTO_INCREMENT PRIMARY KEY,
    `community_id`  BIGINT    NOT NULL,
    `restaurant_id` BIGINT    NOT NULL,
    `created_at`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    FOREIGN KEY `fk_restaurant_id_restaurants_id` (restaurant_id) REFERENCES `restaurants` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    UNIQUE `u_community_id_restaurant_id` (community_id, restaurant_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `users`
(
    `id`                 BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`               VARCHAR(255) NOT NULL,
    `profile_image_file` VARCHAR(511) COMMENT 'プロフィール画像のファイル名',
    `created_at`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `affiliation`
(
    `id`           BIGINT AUTO_INCREMENT PRIMARY KEY,
    `community_id` BIGINT    NOT NULL,
    `user_id`      BIGINT    NOT NULL,
    `created_at`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    FOREIGN KEY `fk_user_id_users_id` (user_id) REFERENCES `users` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    UNIQUE `u_community_id_user_id` (community_id, user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `bookmarks`
(
    `id`           BIGINT AUTO_INCREMENT PRIMARY KEY,
    `community_id` BIGINT    NOT NULL,
    `user_id`      BIGINT    NOT NULL,
    `created_at`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    FOREIGN KEY `fk_user_id_users_id` (user_id) REFERENCES `users` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    UNIQUE `u_community_id_user_id` (community_id, user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `comments`
(
    `id`            BIGINT AUTO_INCREMENT PRIMARY KEY,
    `community_id`  BIGINT        NOT NULL,
    `restaurant_id` BIGINT        NOT NULL,
    `body`          VARCHAR(1023) NOT NULL DEFAULT '',
    `created_at`    TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    FOREIGN KEY `fk_restaurant_id_restaurants_id` (restaurant_id) REFERENCES `restaurants` (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    UNIQUE `u_community_id_restaurant_id` (community_id, restaurant_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

