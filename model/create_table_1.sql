CREATE DATABASE `MutliTable`;

USE `MutliTable`;

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user`
(
    `id`          bigint(20)                             NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识用户记录',
    `user_id`     bigint(20)                             NOT NULL COMMENT '用户ID，用于业务中的用户唯一标识',
    `username`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名，唯一且不区分大小写',
    `password`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户密码，存储的是哈希值',
    `email`       varchar(64) COLLATE utf8mb4_general_ci COMMENT '用户邮箱，可为空',
    `create_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP COMMENT '记录的创建时间',
    `update_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录的最后更新时间',
    `delete_time` bigint                           NULL DEFAULT 0 COMMENT '逻辑删除时间，NULL表示未删除',

    PRIMARY KEY (`id`) COMMENT '主键索引',

    -- 联合唯一索引：确保未删除的用户名唯一
    UNIQUE INDEX `idx_username_delete_time` (`username`, `delete_time`) USING BTREE COMMENT '联合索引：用户名和删除时间确保未删除的用户名唯一',

    -- 联合唯一索引：确保未删除的用户ID唯一
    UNIQUE INDEX `idx_user_id_delete_time` (`user_id`, `delete_time`) USING BTREE COMMENT '联合索引：用户ID和删除时间确保未删除的用户ID唯一',

     -- 联合唯一索引：确保未删除的邮箱唯一
    UNIQUE INDEX `idx_email_delete_time` (`email`, `delete_time`) USING BTREE COMMENT '联合索引：用户ID和删除时间确保未删除的用户ID唯一'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT = '用户信息表：存储用户基本信息及状态';

DROP TABLE IF EXISTS `documents`;
CREATE TABLE `documents` (
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识文档记录',
    `title`       varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '文档标题',
    `owner_id`    bigint(20) NOT NULL COMMENT '文档所有者ID，关联用户记录',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录的创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录的最后更新时间',
    `delete_time` bigint(20) NULL DEFAULT 0 COMMENT '逻辑删除时间，0表示未删除',

    PRIMARY KEY (`id`) COMMENT '主键索引',
    INDEX `idx_owner_id` (`owner_id`) USING BTREE COMMENT '索引：快速查询文档所属用户',

    CONSTRAINT `fk_documents_owner` FOREIGN KEY (`owner_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_general_ci
  COMMENT='文档表：存储在线编辑文档的基本信息';



DROP TABLE IF EXISTS `document_permissions`;
CREATE TABLE `document_permissions` (
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识权限记录',
    `document_id` bigint(20) NOT NULL COMMENT '关联文档ID，关联 documents 表',
    `user_id`     bigint(20) NOT NULL COMMENT '关联用户ID，关联 user 表',
    `permission`  ENUM('owner','editor','viewer') NOT NULL DEFAULT 'viewer' COMMENT '用户在文档中的权限',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录的创建时间',

    PRIMARY KEY (`id`) COMMENT '主键索引',
    UNIQUE INDEX `uk_document_user` (`document_id`, `user_id`) USING BTREE COMMENT '联合唯一索引，确保单个用户在同一文档中的权限唯一',

    CONSTRAINT `fk_permissions_document` FOREIGN KEY (`document_id`) REFERENCES `documents`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_permissions_user` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_general_ci
  COMMENT='文档权限表：存储文档共享及权限控制信息';

DROP TABLE IF EXISTS `sheets`;
CREATE TABLE `sheets` (
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识工作表记录',
    `document_id` bigint(20) NOT NULL COMMENT '所属文档ID，关联 documents 表',
    `name`        varchar(100) COLLATE utf8mb4_general_ci NOT NULL COMMENT '工作表名称',
    `order_index` int NOT NULL DEFAULT 0 COMMENT '工作表排序索引',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录的创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录的最后更新时间',
    `delete_time` bigint(20) NULL DEFAULT 0 COMMENT '逻辑删除时间，0表示未删除',

    PRIMARY KEY (`id`) COMMENT '主键索引',
    INDEX `idx_document_id` (`document_id`) USING BTREE COMMENT '索引：快速查询所属文档的工作表',

    CONSTRAINT `fk_sheets_document` FOREIGN KEY (`document_id`) REFERENCES `documents`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_general_ci
  COMMENT='工作表表：存储在线编辑文档中各个工作表的基本信息';

DROP TABLE IF EXISTS `cells`;
CREATE TABLE `cells` (
    `id`         bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识单元格记录',
    `sheet_id`   bigint(20) NOT NULL COMMENT '所属工作表ID，关联 sheets 表',
    `row`        int NOT NULL COMMENT '单元格所在行号',
    `column`     int NOT NULL COMMENT '单元格所在列号',
    `value`      text COLLATE utf8mb4_general_ci COMMENT '单元格内容',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '单元格最后更新时间',

    PRIMARY KEY (`id`) COMMENT '主键索引',
    UNIQUE INDEX `uk_sheet_row_column` (`sheet_id`, `row`, `column`) USING BTREE COMMENT '联合唯一索引，确保同一工作表中单元格唯一',

    CONSTRAINT `fk_cells_sheet` FOREIGN KEY (`sheet_id`) REFERENCES `sheets`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_general_ci
  COMMENT='单元格表：存储工作表中各个单元格的数据';

DROP TABLE IF EXISTS `cell_history`;
CREATE TABLE `cell_history` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `cell_id` INT NOT NULL,
  `sheet_id` INT NOT NULL,
  `document_id` INT NOT NULL,
  `user_id` INT NOT NULL,
  `old_value` TEXT,
  `new_value` TEXT,
  `changed_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`cell_id`) REFERENCES `cells`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`sheet_id`) REFERENCES `sheets`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`document_id`) REFERENCES `documents`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



