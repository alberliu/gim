-- MySQL dump 10.13  Distrib 8.4.4, for macos15.2 (arm64)
--
-- Host: 127.0.0.1    Database: gim
-- ------------------------------------------------------
-- Server version	8.4.3
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `gim`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `gim` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `gim`;

--
-- Table structure for table `device`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `device` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '账户id',
  `type` tinyint NOT NULL COMMENT '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web',
  `brand` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '手机厂商',
  `model` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '机型',
  `system_version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '系统版本',
  `sdk_version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'app版本',
  `brand_push_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '厂商推送ID',
  `connect_addr` varchar(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '连接层服务器地址',
  `client_addr` varchar(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '客户端地址',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='设备';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `friend`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friend` (
  `user_id` bigint unsigned NOT NULL COMMENT '用户id',
  `friend_id` bigint unsigned NOT NULL COMMENT '好友id',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `remarks` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '备注',
  `extra` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '附加属性',
  `status` tinyint NOT NULL COMMENT '状态，1：申请，2：同意',
  PRIMARY KEY (`user_id`,`friend_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='好友';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `group`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '群组名称',
  `avatar_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '群组头像',
  `introduction` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '群组简介',
  `extra` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '附加属性',
  `members` json NOT NULL COMMENT '群组成员',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='群组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `message`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `request_id` bigint NOT NULL COMMENT '请求id',
  `command` int NOT NULL COMMENT '消息类型',
  `content` blob NOT NULL COMMENT '消息内容',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '消息状态，0：未处理1：消息撤回',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='消息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `seq`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `seq` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `object_type` tinyint NOT NULL COMMENT '对象类型,1:用户；2：群组',
  `object_id` bigint unsigned NOT NULL COMMENT '对象id',
  `seq` bigint unsigned NOT NULL COMMENT '序列号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_object` (`object_type`,`object_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='序列号';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `phone_number` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '手机号',
  `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '昵称',
  `sex` tinyint NOT NULL COMMENT '性别，0:未知；1:男；2:女',
  `avatar_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户头像链接',
  `extra` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '附加属性',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone_number` (`phone_number`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_message`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_message` (
  `user_id` bigint unsigned NOT NULL COMMENT '所属类型的id',
  `seq` bigint unsigned NOT NULL COMMENT '消息序列号',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `message_id` bigint unsigned NOT NULL COMMENT '消息ID',
  PRIMARY KEY (`user_id`,`seq`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户消息';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-08-31 21:51:04
