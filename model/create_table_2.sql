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

-- 班级表
DROP TABLE IF EXISTS `class`;
CREATE TABLE `class` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '班级名称',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` bigint NULL DEFAULT 0 COMMENT '逻辑删除时间戳',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT = '班级表:存储班级信息';

-- 工作表
DROP TABLE IF EXISTS `sheet`;
CREATE TABLE `sheet` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '工作表名称',
  `week` int NOT NULL COMMENT '周数',
  `row` int NOT NULL COMMENT '行数',
  `col` int NOT NULL COMMENT '列数',
  `creator_id` bigint(20) NOT NULL COMMENT '创建者ID（关联user.id）',
  `class_id` bigint(20) NOT NULL COMMENT '班级ID（关联class.id）',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` bigint NULL DEFAULT 0 COMMENT '逻辑删除时间戳',
  PRIMARY KEY (`id`),
  INDEX `idx_creator` (`creator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='工作表主表';

-- 单元格表（核心数据）
DROP TABLE IF EXISTS `cell`;
CREATE TABLE `cell` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sheet_id` bigint(20) NOT NULL COMMENT '所属工作表ID',
  `row_index` int NOT NULL COMMENT '行号（从1开始）',
  `col_index` int NOT NULL COMMENT '列号（从1开始）',
  -- `content` text COLLATE utf8mb4_general_ci,
  `item_id` bigint(20) DEFAULT NULL COMMENT '关联的可拖放元素ID',
  `last_modified_by` bigint(20) DEFAULT NULL COMMENT '最后修改者ID',
  `version` int NOT NULL DEFAULT 1 COMMENT '版本号',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` bigint NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_cell_position` (`sheet_id`, `row_index`, `col_index`, `delete_time`),
  INDEX `fk_item` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='单元格数据表';

-- 可拖放元素表(课程)
DROP TABLE IF EXISTS `draggable_item`;
CREATE TABLE `draggable_item` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `content` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '课程名称',
  `week_type` ENUM('single', 'double', 'all') NOT NULL COMMENT '周类型：单周/双周/全上',
  `classroom` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '上课教室',
  `creator_id` bigint(20) NOT NULL COMMENT '创建者ID',
  `teacher` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '任课老师',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` bigint NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `idx_creator` (`creator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='可拖放元素库';


-- 多班级复用
DROP TABLE IF EXISTS `draggable_class_sheet`;
CREATE TABLE `draggable_class_sheet` (
  `item_id` bigint(20) NOT NULL COMMENT '可拖动元素ID',
  `class_id` bigint(20) NOT NULL COMMENT '班级ID',
  PRIMARY KEY (`item_id`, `class_id`),
  INDEX `idx_sheet` (`class_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='元素-课表关联表';
-- 权限管理表
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `sheet_id` bigint(20) NOT NULL COMMENT '工作表ID',
  -- `access_level` enum('READ','WRITE','ADMIN') NOT NULL DEFAULT 'READ',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` bigint NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_user_sheet` (`user_id`, `sheet_id`, `delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='权限控制表';
