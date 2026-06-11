-- 去掉 platform / common 伪租户码，统一为 default
-- 部署新代码前在维护窗口执行；执行前请备份相关表

UPDATE sys_user_info SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE sys_notify_template SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE sys_notify_channel SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE sys_notify_config_template SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');

UPDATE dm_product_info SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE dm_product_schema SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE dm_protocol_script SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common', '__common __');
UPDATE dm_protocol_script_device SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE dm_ota_firmware_info SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
UPDATE dm_ota_firmware_job SET tenant_code = 'default' WHERE tenant_code IN ('platform', 'common');
