CREATE TABLE IF NOT EXISTS `communities`
(
    `id`          BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`        VARCHAR(255) NOT NULL,
    `description` VARCHAR(511) NOT NULL,
    `location`    GEOMETRY,
    `image_file`  VARCHAR(255) NOT NULL COMMENT '画像のファイル名',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `restaurants`
(
    `id`         BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`       VARCHAR(255) NOT NULL,
    `location`   GEOMETRY,
    `url`        VARCHAR(255) NOT NULL,
    `image_url`  VARCHAR(255) COMMENT '画像のURL',
    `source`     VARCHAR(255) NOT NULL COMMENT 'レストラン情報の取得元',
    `source_id`  VARCHAR(255) NOT NULL COMMENT 'レストラン情報の取得元のレストランID',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `communities_restaurants`
(
    `id`            BIGINT AUTO_INCREMENT PRIMARY KEY,
    `community_id`  BIGINT   NOT NULL,
    `restaurant_id` BIGINT   NOT NULL,
    `created_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id),
    FOREIGN KEY `fk_restaurant_id_restaurants_id` (restaurant_id) REFERENCES `restaurants` (id),
    UNIQUE `u_community_id_restaurant_id` (community_id, restaurant_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `users`
(
    `id`                 BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`               VARCHAR(255) NOT NULL,
    `profile_image_file` VARCHAR(511) COMMENT 'プロフィール画像のファイル名',
    `created_at`         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `affiliation`
(
    `id`           BIGINT AUTO_INCREMENT PRIMARY KEY,
    `community_id` BIGINT   NOT NULL,
    `user_id`      BIGINT   NOT NULL,
    `created_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id),
    FOREIGN KEY `fk_user_id_users_id` (user_id) REFERENCES `users` (id),
    UNIQUE `u_community_id_user_id` (community_id, user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `bookmarks`
(
    `id`           BIGINT AUTO_INCREMENT PRIMARY KEY,
    `community_id` BIGINT   NOT NULL,
    `user_id`      BIGINT   NOT NULL,
    `created_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id),
    FOREIGN KEY `fk_user_id_users_id` (user_id) REFERENCES `users` (id),
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
    `created_at`    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY `fk_community_id_communities_id` (community_id) REFERENCES `communities` (id),
    FOREIGN KEY `fk_restaurant_id_restaurants_id` (restaurant_id) REFERENCES `restaurants` (id),
    UNIQUE `u_community_id_restaurant_id` (community_id, restaurant_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

