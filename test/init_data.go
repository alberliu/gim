package test

import (
	"gim/pkg/db"
)

func initData() {
	db.Init()

	if err := db.DB.Exec("DELETE FROM `user` WHERE `id` IN (1, 2, 3)").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("DELETE FROM `device` WHERE `id` IN (11, 12, 2, 3)").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("DELETE FROM `group` WHERE `id` = 1").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("DELETE FROM `group_member` WHERE `group_id` = 1 AND `user_id` IN (1, 2, 3)").Error; err != nil {
		panic(err)
	}

	// 插入 user 数据
	if err := db.DB.Exec("INSERT INTO `user` (`id`, `created_at`, `updated_at`, `phone_number`, `nickname`, `sex`, `avatar_url`, `extra`) VALUES (1, now(), now(), '1', 'alber1', 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("INSERT INTO `user` (`id`, `created_at`, `updated_at`, `phone_number`, `nickname`, `sex`, `avatar_url`, `extra`) VALUES (2, now(), now(), '2', 'alber2', 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("INSERT INTO `user` (`id`, `created_at`, `updated_at`, `phone_number`, `nickname`, `sex`, `avatar_url`, `extra`) VALUES (3, now(), now(), '3', 'alber3', 0, '', '')").Error; err != nil {
		panic(err)
	}

	// 插入 device 数据
	if err := db.DB.Exec("INSERT INTO `device` (`id`, `created_at`, `updated_at`, `user_id`, `type`, `brand`, `model`, `system_version`, `sdk_version`, `push_token`, `connect_ip`, `client_addr`) VALUES (11, now(), now(), 1, 1, 'xiaomi', 'xiaomi 15', '15.0.0', '1.0.0', 'xiaomi push id', '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("INSERT INTO `device` (`id`, `created_at`, `updated_at`, `user_id`, `type`, `brand`, `model`, `system_version`, `sdk_version`, `push_token`, `connect_ip`, `client_addr`) VALUES (12, now(), now(), 0, 1, 'xiaomi', 'xiaomi 15', '15.0.0', '1.0.0', 'xiaomi push id', '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("INSERT INTO `device` (`id`, `created_at`, `updated_at`, `user_id`, `type`, `brand`, `model`, `system_version`, `sdk_version`, `push_token`, `connect_ip`, `client_addr`) VALUES (2, now(), now(), 0, 1, 'xiaomi', 'xiaomi 15', '15.0.0', '1.0.0', 'xiaomi push id', '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("INSERT INTO `device` (`id`, `created_at`, `updated_at`, `user_id`, `type`, `brand`, `model`, `system_version`, `sdk_version`, `push_token`, `connect_ip`, `client_addr`) VALUES (3, now(), now(), 0, 1, 'xiaomi', 'xiaomi 15', '15.0.0', '1.0.0', 'xiaomi push id', '', '')").Error; err != nil {
		panic(err)
	}

	// 插入 group 数据
	if err := db.DB.Exec("INSERT INTO `group` (`id`, `created_at`, `updated_at`, `name`, `avatar_url`, `introduction`, `extra`) VALUES (1, now(), now(), 'group', 'group', 'group', 'group')").Error; err != nil {
		panic(err)
	}

	// 插入 group_member 数据
	if err := db.DB.Exec("INSERT INTO `group_member` (`group_id`, `user_id`, `created_at`, `updated_at`, `nickname`, `type`, `status`, `extra`) VALUES (1, 1, now(), now(), '', 0, 0, '')").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("INSERT INTO `group_member` (`group_id`, `user_id`, `created_at`, `updated_at`, `nickname`, `type`, `status`, `extra`) VALUES (1, 2, now(), now(), '', 0, 0, '')").Error; err != nil {
		panic(err)
	}
	if err := db.DB.Exec("INSERT INTO `group_member` (`group_id`, `user_id`, `created_at`, `updated_at`, `nickname`, `type`, `status`, `extra`) VALUES (1, 3, now(), now(), '', 0, 0, '')").Error; err != nil {
		panic(err)
	}
}
