SET NAMES utf8mb4;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `gim` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `gim`;

CREATE TABLE `device` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '账户id',
  `type` tinyint NOT NULL COMMENT '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web',
  `brand` varchar(20) NOT NULL COMMENT '手机厂商',
  `model` varchar(20) NOT NULL COMMENT '机型',
  `system_version` varchar(10) NOT NULL COMMENT '系统版本',
  `sdk_version` varchar(10) NOT NULL COMMENT 'app版本',
  `push_token` varchar(1024) NOT NULL COMMENT '厂商推送token',
  `connect_ip` varchar(25) NOT NULL COMMENT '连接层服务器IP',
  `client_addr` varchar(25) NOT NULL COMMENT '客户端地址',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10000 COMMENT='设备';

CREATE TABLE `device_ack` (
  `device_id` bigint unsigned NOT NULL COMMENT '设备ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `ack` bigint unsigned NOT NULL COMMENT '序列号',
  PRIMARY KEY (`device_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB COMMENT='ack';

CREATE TABLE `friend` (
  `user_id` bigint unsigned NOT NULL COMMENT '用户id',
  `friend_id` bigint unsigned NOT NULL COMMENT '好友id',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `remarks` varchar(20) NOT NULL COMMENT '备注',
  `extra` varchar(1024) NOT NULL COMMENT '附加属性',
  `status` tinyint NOT NULL COMMENT '状态，1：申请，2：同意',
  PRIMARY KEY (`user_id`,`friend_id`) USING BTREE
) ENGINE=InnoDB COMMENT='好友';

CREATE TABLE `group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `name` varchar(50) NOT NULL COMMENT '群组名称',
  `avatar_url` varchar(255) NOT NULL COMMENT '群组头像',
  `introduction` varchar(255) NOT NULL COMMENT '群组简介',
  `extra` varchar(1024) NOT NULL COMMENT '附加属性',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 COMMENT='群组';

CREATE TABLE `group_member` (
  `group_id` bigint unsigned NOT NULL COMMENT '群组ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `nickname` varchar(255) NOT NULL COMMENT '昵称',
  `type` tinyint NOT NULL COMMENT '类型',
  `status` tinyint NOT NULL COMMENT '状态',
  `extra` varchar(1024) NOT NULL COMMENT '附加属性',
  PRIMARY KEY (`group_id`,`user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB COMMENT='群组成员';

CREATE TABLE `message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `request_id` varchar(128) NOT NULL COMMENT '请求id',
  `command` int NOT NULL COMMENT '消息类型',
  `content` blob NOT NULL COMMENT '消息内容',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 COMMENT='消息'
/*!50100 PARTITION BY HASH (`id`)
PARTITIONS 8 */;

CREATE TABLE `seq` (
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `seq` bigint unsigned NOT NULL COMMENT '序列号',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB COMMENT='序列号';

CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `phone_number` varchar(20) NOT NULL COMMENT '手机号',
  `nickname` varchar(64) NOT NULL COMMENT '昵称',
  `sex` tinyint NOT NULL COMMENT '性别，0:未知；1:男；2:女',
  `avatar_url` varchar(256) NOT NULL COMMENT '用户头像链接',
  `extra` varchar(1024) NOT NULL COMMENT '附加属性',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone_number` (`phone_number`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 COMMENT='用户';

CREATE TABLE `user_message` (
  `user_id` bigint unsigned NOT NULL COMMENT '所属类型的id',
  `seq` bigint unsigned NOT NULL COMMENT '消息序列号',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `message_id` bigint unsigned NOT NULL COMMENT '消息ID',
  PRIMARY KEY (`user_id`,`seq`) USING BTREE
) ENGINE=InnoDB COMMENT='用户消息'
/*!50100 PARTITION BY HASH (`user_id`)
PARTITIONS 8 */;
