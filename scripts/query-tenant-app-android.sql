-- 核对 Apple 配置是否可被 syssvr 正确加载（方案 A 排查用）
SELECT app_code, login_types,
       apple_app_id, apple_team_id, apple_key_id, apple_redirect_uri,
       LENGTH(IFNULL(apple_private_key,'')) AS pk_len,
       LEFT(IFNULL(apple_private_key,''), 30) AS pk_head,
       CASE WHEN apple_private_key LIKE '-----BEGIN%' THEN 'PEM_OK' ELSE 'PEM_BAD' END AS pem_check
FROM sys_tenant_app
WHERE tenant_code = 'default' AND app_code IN ('client-app-android', 'client-app-ios');
