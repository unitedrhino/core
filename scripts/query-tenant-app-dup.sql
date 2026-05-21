SELECT id, tenant_code, app_code, deleted_time, login_types,
       apple_app_id, apple_team_id, apple_key_id, LENGTH(IFNULL(apple_private_key,'')) AS pk_len
FROM sys_tenant_app
WHERE tenant_code = 'default' AND app_code = 'client-app-android'
ORDER BY id;
