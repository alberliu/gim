SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(50) NOT NULL COMMENT '名称',
  `private_key` varchar(1024) NOT NULL COMMENT '私钥',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='app';

INSERT INTO `gim`.`app`(`id`, `name`, `private_key`, `create_time`, `update_time`) VALUES (1, '测试', '-----BEGIN RSA PRIVATE KEY-----\nMIICWwIBAAKBgQDcGsUIIAINHfRTdMmgGwLrjzfMNSrtgIf4EGsNaYwmC1GjF/bM\nh0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdTnCDPPZ7oV7p1B9Pud+6zPaco\nqDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Zy682X1+R1lRK8D+vmQIDAQAB\nAoGAeWAZvz1HZExca5k/hpbeqV+0+VtobMgwMs96+U53BpO/VRzl8Cu3CpNyb7HY\n64L9YQ+J5QgpPhqkgIO0dMu/0RIXsmhvr2gcxmKObcqT3JQ6S4rjHTln49I2sYTz\n7JEH4TcplKjSjHyq5MhHfA+CV2/AB2BO6G8limu7SheXuvECQQDwOpZrZDeTOOBk\nz1vercawd+J9ll/FZYttnrWYTI1sSF1sNfZ7dUXPyYPQFZ0LQ1bhZGmWBZ6a6wd9\nR+PKlmJvAkEA6o32c/WEXxW2zeh18sOO4wqUiBYq3L3hFObhcsUAY8jfykQefW8q\nyPuuL02jLIajFWd0itjvIrzWnVmoUuXydwJAXGLrvllIVkIlah+lATprkypH3Gyc\nYFnxCTNkOzIVoXMjGp6WMFylgIfLPZdSUiaPnxby1FNM7987fh7Lp/m12QJAK9iL\n2JNtwkSR3p305oOuAz0oFORn8MnB+KFMRaMT9pNHWk0vke0lB1sc7ZTKyvkEJW0o\neQgic9DvIYzwDUcU8wJAIkKROzuzLi9AvLnLUrSdI6998lmeYO9x7pwZPukz3era\nzncjRK3pbVkv0KrKfczuJiRlZ7dUzVO0b6QJr8TRAA==\n-----END RSA PRIVATE KEY-----', '2019-10-15 16:49:39', '2019-10-15 16:49:39');
-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `device_id` bigint NOT NULL COMMENT '设备id',
  `app_id` bigint unsigned NOT NULL COMMENT 'app_id',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '账户id',
  `type` tinyint NOT NULL COMMENT '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web',
  `brand` varchar(20) NOT NULL COMMENT '手机厂商',
  `model` varchar(20) NOT NULL COMMENT '机型',
  `system_version` varchar(10) NOT NULL COMMENT '系统版本',
  `sdk_version` varchar(10) NOT NULL COMMENT 'app版本',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '在线状态，0：离线；1：在线',
  `conn_addr` varchar(25) NOT NULL COMMENT '连接层服务器地址',
  `conn_fd` bigint NOT NULL COMMENT 'TCP连接对应的文件描述符',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_device_id` (`device_id`) USING BTREE,
  KEY `idx_app_id_user_id` (`app_id`,`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='设备';

-- ----------------------------
-- Table structure for device_ack
-- ----------------------------
DROP TABLE IF EXISTS `device_ack`;
CREATE TABLE `device_ack` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `device_id` bigint unsigned NOT NULL COMMENT '设备id',
  `ack` bigint unsigned NOT NULL DEFAULT '0' COMMENT '收到消息确认号',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_device_id` (`device_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='设备消息同步序列号';

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint NOT NULL COMMENT 'app_id',
  `group_id` bigint NOT NULL COMMENT '群组id',
  `name` varchar(50) NOT NULL COMMENT '群组名称',
  `introduction` varchar(255) NOT NULL COMMENT '群组简介',
  `user_num` int NOT NULL DEFAULT '0' COMMENT '群组人数',
  `type` tinyint NOT NULL COMMENT '群组类型，1：小群；2：大群',
  `extra` varchar(1024) NOT NULL COMMENT '附加属性',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_app_id_group_id` (`app_id`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='群组';

-- ----------------------------
-- Table structure for group_user
-- ----------------------------
DROP TABLE IF EXISTS `group_user`;
CREATE TABLE `group_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint NOT NULL COMMENT 'app_id',
  `group_id` bigint unsigned NOT NULL COMMENT '组id',
  `user_id` bigint unsigned NOT NULL COMMENT '用户id',
  `label` varchar(20) NOT NULL COMMENT '用户在群组的昵称',
  `extra` varchar(1024) NOT NULL COMMENT '附加属性',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_app_id_group_id_user_id` (`app_id`,`group_id`,`user_id`) USING BTREE,
  KEY `idx_app_id_user_id` (`app_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='群组成员关系';

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint unsigned NOT NULL COMMENT 'app_id',
  `object_type` tinyint NOT NULL COMMENT '所属类型，1：用户；2：群组',
  `object_id` bigint unsigned NOT NULL COMMENT '所属类型的id',
  `request_id` bigint NOT NULL COMMENT '请求id',
  `sender_type` tinyint NOT NULL COMMENT '发送者类型',
  `sender_id` bigint unsigned NOT NULL COMMENT '发送者id',
  `sender_device_id` bigint unsigned NOT NULL COMMENT '发送设备id',
  `receiver_type` tinyint NOT NULL COMMENT '接收者类型,1:个人；2：群组',
  `receiver_id` bigint unsigned NOT NULL COMMENT '接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id',
  `to_user_ids` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '需要@的用户id列表，多个用户用，隔开',
  `type` tinyint NOT NULL COMMENT '消息类型',
  `content` varchar(4094) COLLATE utf8mb4_bin NOT NULL COMMENT '消息内容',
  `seq` bigint unsigned NOT NULL COMMENT '消息序列号',
  `send_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息发送时间',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '消息状态，0：未处理1：消息撤回',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_app_id_object_seq` (`app_id`,`object_type`,`object_id`,`seq`) USING BTREE,
  KEY `idx_request_id` (`request_id`) USING BTREE
) ENGINE=InnoDB COLLATE=utf8mb4_bin COMMENT='消息';

-- ----------------------------
-- Table structure for uid
-- ----------------------------
DROP TABLE IF EXISTS `uid`;
CREATE TABLE `uid` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `business_id` varchar(128) NOT NULL COMMENT '业务id',
  `max_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '最大id',
  `step` int unsigned NOT NULL DEFAULT '1000' COMMENT '步长',
  `description` varchar(255) NOT NULL COMMENT '描述',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_business_id` (`business_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='分布式自增主键';

INSERT INTO `gim`.`uid`(`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`) VALUES (1, 'device_id', 0, 10, '设备id', '2019-10-15 16:42:05', '2020-04-29 23:33:40');
-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint unsigned NOT NULL COMMENT 'app_id',
  `user_id` bigint unsigned NOT NULL COMMENT '用户id',
  `nickname` varchar(20) NOT NULL COMMENT '昵称',
  `sex` tinyint NOT NULL COMMENT '性别，0:未知；1:男；2:女',
  `avatar_url` varchar(50) NOT NULL COMMENT '用户头像链接',
  `extra` varchar(1024) NOT NULL COMMENT '附加属性',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_app_id_user_id` (`app_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户';

SET FOREIGN_KEY_CHECKS = 1;
