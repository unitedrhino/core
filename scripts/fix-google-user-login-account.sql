-- 修复历史 Google 用户：补全 user_name，使邮箱可作为密码登录 account
-- 执行前请备份；在对应租户库执行

UPDATE sys_user_info
SET user_name = email
WHERE (user_name IS NULL OR user_name = '')
  AND email IS NOT NULL AND email != ''
  AND google_user_id IS NOT NULL AND google_user_id != '';
