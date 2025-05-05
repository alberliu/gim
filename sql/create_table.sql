SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device`
(
    `id`             bigint unsigned                 NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `created_at`     datetime                        NOT NULL COMMENT '创建时间',
    `updated_at`     datetime                        NOT NULL COMMENT '更新时间',
    `user_id`        bigint unsigned                 NOT NULL DEFAULT '0' COMMENT '账户id',
    `type`           tinyint                         NOT NULL COMMENT '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web',
    `brand`          varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '手机厂商',
    `model`          varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '机型',
    `system_version` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT '系统版本',
    `sdk_version`    varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'app版本',
    `status`         tinyint                         NOT NULL DEFAULT '0' COMMENT '在线状态，0：离线；1：在线',
    `conn_addr`      varchar(25) COLLATE utf8mb4_bin NOT NULL COMMENT '连接层服务器地址',
    `client_addr`    varchar(25) COLLATE utf8mb4_bin NOT NULL COMMENT '客户端地址',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 3
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='设备';

-- ----------------------------
-- Table structure for friend
-- ----------------------------
DROP TABLE IF EXISTS `friend`;
CREATE TABLE `friend`
(
    `user_id`    bigint unsigned                   NOT NULL COMMENT '用户id',
    `friend_id`  bigint unsigned                   NOT NULL COMMENT '好友id',
    `created_at` datetime                          NOT NULL COMMENT '创建时间',
    `updated_at` datetime                          NOT NULL COMMENT '更新时间',
    `remarks`    varchar(20) COLLATE utf8mb4_bin   NOT NULL COMMENT '备注',
    `extra`      varchar(1024) COLLATE utf8mb4_bin NOT NULL COMMENT '附加属性',
    `status`     tinyint                           NOT NULL COMMENT '状态，1：申请，2：同意',
    PRIMARY KEY (`user_id`, `friend_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='好友';

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group`
(
    `id`           bigint unsigned                   NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `created_at`   datetime                          NOT NULL COMMENT '创建时间',
    `updated_at`   datetime                          NOT NULL COMMENT '更新时间',
    `name`         varchar(50) COLLATE utf8mb4_bin   NOT NULL COMMENT '群组名称',
    `avatar_url`   varchar(255) COLLATE utf8mb4_bin  NOT NULL COMMENT '群组头像',
    `introduction` varchar(255) COLLATE utf8mb4_bin  NOT NULL COMMENT '群组简介',
    `user_num`     int                               NOT NULL DEFAULT '0' COMMENT '群组人数',
    `extra`        varchar(1024) COLLATE utf8mb4_bin NOT NULL COMMENT '附加属性',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 6
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='群组';

-- ----------------------------
-- Table structure for group_user
-- ----------------------------
DROP TABLE IF EXISTS `group_user`;
CREATE TABLE `group_user`
(
    `group_id`    bigint unsigned                   NOT NULL COMMENT '组id',
    `user_id`     bigint unsigned                   NOT NULL COMMENT '用户id',
    `created_at`  datetime                          NOT NULL COMMENT '创建时间',
    `updated_at`  datetime                          NOT NULL COMMENT '更新时间',
    `member_type` tinyint                           NOT NULL COMMENT '成员类型，1：管理员；2：普通成员',
    `remarks`     varchar(20) COLLATE utf8mb4_bin   NOT NULL COMMENT '备注',
    `extra`       varchar(1024) COLLATE utf8mb4_bin NOT NULL COMMENT '附加属性',
    `status`      tinyint                           NOT NULL COMMENT '状态',
    PRIMARY KEY (`group_id`, `user_id`),
    KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='群组成员';

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `created_at` datetime        NOT NULL COMMENT '创建时间',
    `updated_at` datetime        NOT NULL COMMENT '更新时间',
    `request_id` bigint          NOT NULL COMMENT '请求id',
    `code`       int             NOT NULL COMMENT '消息类型',
    `content`    blob            NOT NULL COMMENT '消息内容',
    `status`     tinyint         NOT NULL DEFAULT '0' COMMENT '消息状态，0：未处理1：消息撤回',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 22
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='消息';

-- ----------------------------
-- Table structure for seq
-- ----------------------------
DROP TABLE IF EXISTS `seq`;
CREATE TABLE `seq`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `created_at`  datetime        NOT NULL COMMENT '创建时间',
    `updated_at`  datetime        NOT NULL COMMENT '更新时间',
    `object_type` tinyint         NOT NULL COMMENT '对象类型,1:用户；2：群组',
    `object_id`   bigint unsigned NOT NULL COMMENT '对象id',
    `seq`         bigint unsigned NOT NULL COMMENT '序列号',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_object` (`object_type`, `object_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 3
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='序列号';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`           bigint unsigned                   NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `created_at`   datetime                          NOT NULL COMMENT '创建时间',
    `updated_at`   datetime                          NOT NULL COMMENT '更新时间',
    `phone_number` varchar(20) COLLATE utf8mb4_bin   NOT NULL COMMENT '手机号',
    `nickname`     varchar(20) COLLATE utf8mb4_bin   NOT NULL COMMENT '昵称',
    `sex`          tinyint                           NOT NULL COMMENT '性别，0:未知；1:男；2:女',
    `avatar_url`   varchar(256) COLLATE utf8mb4_bin  NOT NULL COMMENT '用户头像链接',
    `extra`        varchar(1024) COLLATE utf8mb4_bin NOT NULL COMMENT '附加属性',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_phone_number` (`phone_number`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户';

-- ----------------------------
-- Table structure for user_message
-- ----------------------------
DROP TABLE IF EXISTS `user_message`;
CREATE TABLE `user_message`
(
    `user_id`    bigint unsigned NOT NULL COMMENT '所属类型的id',
    `seq`        bigint unsigned NOT NULL COMMENT '消息序列号',
    `created_at` datetime        NOT NULL COMMENT '创建时间',
    `updated_at` datetime        NOT NULL COMMENT '更新时间',
    `message_id` bigint unsigned NOT NULL COMMENT '消息ID',
    PRIMARY KEY (`user_id`, `seq`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户消息';

SET FOREIGN_KEY_CHECKS = 1;
