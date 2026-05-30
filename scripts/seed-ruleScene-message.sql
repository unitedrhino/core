-- 已有环境补全「场景联动通知 ruleScene + 站内信 message」模板与配置
-- 执行前请备份；可重复执行（INSERT IGNORE / ON DUPLICATE KEY 语义依赖表唯一索引）

-- 1. 补全 ruleScene 支持类型（若 JSON 列尚未含 message）
UPDATE sys_notify_config
SET support_types = JSON_ARRAY('message','sms','email','phoneCall','dingWebhook','wxEWebHook','wxMini','dingTalk','dingMini'),
    enable_types  = JSON_ARRAY('message')
WHERE code = 'ruleScene'
  AND deleted_time = 0
  AND (JSON_CONTAINS(support_types, '"message"') = 0 OR JSON_LENGTH(enable_types) = 0);

-- 2. 通用站内信模板（id=1 与初始化种子一致）
INSERT INTO sys_notify_template (id, tenant_code, name, notify_code, type, code, sign_name, subject, body, `desc`, channel_id, deleted_time, created_time, updated_time)
SELECT 1, 'common', '场景联动站内信', 'ruleScene', 'message', 'ruleScene_message', '', '{{.title}}', '{{.body}}', '场景联动站内信默认模板', 0, 0, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (
    SELECT 1 FROM sys_notify_template WHERE notify_code = 'ruleScene' AND type = 'message' AND deleted_time = 0
);

-- 3. 默认平台绑定模板（tenant_code 按实际默认平台调整，一般为 default）
INSERT INTO sys_notify_config_template (tenant_code, notify_code, type, template_id, deleted_time, created_time, updated_time)
SELECT 'default', 'ruleScene', 'message', t.id, 0, NOW(), NOW()
FROM sys_notify_template t
WHERE t.notify_code = 'ruleScene' AND t.type = 'message' AND t.deleted_time = 0
  AND NOT EXISTS (
    SELECT 1 FROM sys_notify_config_template c
    WHERE c.notify_code = 'ruleScene' AND c.type = 'message' AND c.deleted_time = 0
  )
LIMIT 1;
