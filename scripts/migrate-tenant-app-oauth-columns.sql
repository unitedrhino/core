-- 扩展 sys_tenant_app 中 Google/GitHub/Apple OAuth 字段长度
-- 原因：Google Client ID 常超过 50 字符（如 *.apps.googleusercontent.com）

ALTER TABLE `sys_tenant_app`
  MODIFY COLUMN `google_app_id` VARCHAR(128) NOT NULL DEFAULT '',
  MODIFY COLUMN `google_app_key` VARCHAR(128) NOT NULL DEFAULT '',
  MODIFY COLUMN `google_app_secret` VARCHAR(512) NOT NULL DEFAULT '',
  MODIFY COLUMN `github_app_id` VARCHAR(128) NOT NULL DEFAULT '',
  MODIFY COLUMN `github_app_key` VARCHAR(128) NOT NULL DEFAULT '',
  MODIFY COLUMN `github_app_secret` VARCHAR(512) NOT NULL DEFAULT '',
  MODIFY COLUMN `apple_app_id` VARCHAR(128) NOT NULL DEFAULT '',
  MODIFY COLUMN `apple_team_id` VARCHAR(50) NOT NULL DEFAULT '',
  MODIFY COLUMN `apple_key_id` VARCHAR(50) NOT NULL DEFAULT '',
  MODIFY COLUMN `apple_private_key` TEXT,
  MODIFY COLUMN `apple_redirect_uri` VARCHAR(512) NOT NULL DEFAULT '';
