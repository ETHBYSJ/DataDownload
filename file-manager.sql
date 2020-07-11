DROP TABLE IF EXISTS `v1_users`;
CREATE TABLE `v1_users` (
    `id` int(8) NOT NULL AUTO_INCREMENT,
    `email` varchar (100),
    `password` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `status` tinyint(1) DEFAULT 1,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    `user_type` varchar(8) DEFAULT 'User',
    `language`  varchar(8)  DEFAULT 'CN',
    PRIMARY KEY(`id`),
    UNIQUE(`email`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `v1_users` (`email`, `password`) VALUES
('2296176046@qq.com', '9PmwyHJhbjcbLOG7:f14272ab12b6a886a4e3e8f7429bcfc94c6b0e78');

DROP TABLE IF EXISTS `v1_files`;
CREATE TABLE `v1_files` (
    `id` int(8) NOT NULL AUTO_INCREMENT,
    `name`  varchar(255) NOT NULL,
    `is_dir` tinyint(1) NOT NULL DEFAULT 0,
    `path`  varchar(255) NOT NULL,
    `owner_id` int(8) NOT NULL,
    `privacy` tinyint(1) NOT NULL DEFAULT 0,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `size`      int(8) NOT NULL DEFAULT 0,
    `review`    tinyint(1) NOT NULL DEFAULT 0,
    PRIMARY KEY(`id`),
    FOREIGN KEY (`owner_id`) REFERENCES `v1_users`(`id`)

) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
ALTER TABLE `v1_files` ADD UNIQUE INDEX idx(`name`, `path`);