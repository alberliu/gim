-- ----------------------------
-- Table structure for t_device
-- ----------------------------
CREATE DATABASE im DEFAULT CHARSET 'utf8';

DROP TABLE IF EXISTS `t_device`;
CREATE TABLE `t_device` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '设备id',
  `user_id` bigint(20) unsigned DEFAULT '0' COMMENT '账户id',
  `token` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '设备登录的token',
  `type` tinyint(3) NOT NULL COMMENT '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web',
  `brand` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '手机厂商',
  `model` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '机型',
  `system_version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '系统版本',
  `app_version` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'app版本',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在线状态，0：离线；1：在线',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='设备';

-- ----------------------------
-- Table structure for t_device_send_sequence
-- ----------------------------
DROP TABLE IF EXISTS `t_device_send_sequence`;
CREATE TABLE `t_device_send_sequence` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `device_id` bigint(20) unsigned NOT NULL COMMENT '设备id',
  `send_sequence` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '设备发送序列号',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_device_id` (`device_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='设备消息发送序列号';

-- ----------------------------
-- Table structure for t_device_sync_sequence
-- ----------------------------
DROP TABLE IF EXISTS `t_device_sync_sequence`;
CREATE TABLE `t_device_sync_sequence` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `device_id` bigint(20) unsigned NOT NULL COMMENT '设备id',
  `sync_sequence` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '已同步序列号',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_device_id` (`device_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='设备消息同步序列号';

-- ----------------------------
-- Table structure for t_friend
-- ----------------------------
DROP TABLE IF EXISTS `t_friend`;
CREATE TABLE `t_friend` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '账户id',
  `friend_id` bigint(20) unsigned NOT NULL COMMENT '好友账户id',
  `label` char(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '备注，标签',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_id`,`friend_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='好友关系';

-- ----------------------------
-- Table structure for t_group
-- ----------------------------
DROP TABLE IF EXISTS `t_group`;
CREATE TABLE `t_group` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '群组id',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '组名',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='群组';

-- ----------------------------
-- Table structure for t_group_user
-- ----------------------------
DROP TABLE IF EXISTS `t_group_user`;
CREATE TABLE `t_group_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `group_id` bigint(20) unsigned NOT NULL COMMENT '组id',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `label` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户在群组的昵称',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`,`user_id`) USING BTREE,
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='群组成员关系';

-- ----------------------------
-- Table structure for t_lid
-- ----------------------------
DROP TABLE IF EXISTS `t_lid`;
CREATE TABLE `t_lid` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '业务id',
  `max_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '最大id',
  `step` int(10) unsigned NOT NULL DEFAULT '1000' COMMENT '步长',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '描述',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_business_id` (`business_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='分布式自增主键';

-- ----------------------------
-- Records of t_lid
-- ----------------------------
BEGIN;
INSERT INTO `t_lid` VALUES (1, 'message_id', 0, 1000, '消息id', '2018-11-11 09:14:09', '2018-11-12 18:30:41');
COMMIT;

-- ----------------------------
-- Table structure for t_message
-- ----------------------------
DROP TABLE IF EXISTS `t_message`;
CREATE TABLE `t_message` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `message_id` bigint(20) unsigned NOT NULL COMMENT '消息id',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `sender_type` tinyint(3) NOT NULL COMMENT '发送者类型',
  `sender_id` bigint(20) unsigned NOT NULL COMMENT '发送者id',
  `sender_device_id` bigint(20) unsigned NOT NULL COMMENT '发送设备id',
  `receiver_type` tinyint(3) NOT NULL COMMENT '接收者类型,1:个人；2：群组',
  `receiver_id` bigint(20) unsigned NOT NULL COMMENT '接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id',
  `type` tinyint(3) NOT NULL COMMENT '消息类型,0：文本；1：语音；2：图片',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '内容',
  `sequence` bigint(20) unsigned NOT NULL COMMENT '消息序列号',
  `send_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '消息发送时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id_sequence` (`user_id`,`sequence`),
  KEY `idx_message_id` (`message_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='消息';

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `number` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '手机号',
  `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '昵称',
  `sex` tinyint(4) NOT NULL COMMENT '性别，0:未知；1:男；2:女',
  `avatar` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户头像',
  `password` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '密码',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_number` (`number`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='账户';

-- ----------------------------
-- Table structure for t_user_sequence
-- ----------------------------
DROP TABLE IF EXISTS `t_user_sequence`;
CREATE TABLE `t_user_sequence` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `sequence` bigint(20) unsigned NOT NULL COMMENT '自增张序列',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户消息的自增长序列';