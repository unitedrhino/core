-- 去掉 platform / common 伪租户码，统一为 default
-- 部署新代码前在维护窗口执行；执行前请备份相关表

UPDATE sys_user_info SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE sys_notify_template SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE sys_notify_channel SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE sys_notify_config_template SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');

UPDATE dm_product_info SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
-- 注意：物模型真实表名为 dm_schema_info（DmProductSchema 结构体 TableName 返回 dm_schema_info），
-- 早期版本此处误写为 dm_product_schema 导致漏迁移，2026-06-11 已通过 fix-schema-info-tenant-code*.mjs 补救：
-- 1) 与 default 行重复（同 product_id+identifier）的 common 行软删；2) 其余 common 行改为 default；3) default 内部重复行保留最新一行软删其余
UPDATE dm_schema_info SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE sys_login_log SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE dm_protocol_script SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common', '__common __');
UPDATE dm_protocol_script_device SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE dm_ota_firmware_info SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE dm_ota_firmware_job SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
