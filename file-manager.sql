DROP TABLE IF EXISTS `v1_users`;
CREATE TABLE `v1_users` (
    `id` int(8) NOT NULL AUTO_INCREMENT,
    `email` varchar (100),
    `password` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `status` int(8) DEFAULT 0,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY(`id`),
    UNIQUE(`email`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `v1_users` (`email`, `password`) VALUES
('2296176046@qq.com', '9PmwyHJhbjcbLOG7:f14272ab12b6a886a4e3e8f7429bcfc94c6b0e78')