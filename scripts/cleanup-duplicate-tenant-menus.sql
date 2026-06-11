-- 清理租户菜单重复数据（template_id 已失效或同 path 多份副本）
-- 执行前请备份 sys_tenant_app_menu、sys_role_menu；部署含修复的 syssvr 后也可由启动导入自动清理

-- 1. 查看孤儿租户菜单（template_id 在模块菜单中不存在）
SELECT tam.id, tam.template_id, tam.tenant_code, tam.app_code, tam.module_code, tam.name, tam.path
FROM sys_tenant_app_menu tam
LEFT JOIN sys_module_menu mm ON tam.template_id = mm.id AND mm.deleted_time = 0
WHERE tam.deleted_time = 0 AND mm.id IS NULL
ORDER BY tam.module_code, tam.path, tam.id;

-- 2. 查看同 path 重复的一级菜单
SELECT tenant_code, app_code, module_code, path, parent_id, COUNT(*) AS cnt
FROM sys_tenant_app_menu
WHERE deleted_time = 0 AND parent_id = 1
GROUP BY tenant_code, app_code, module_code, path, parent_id
HAVING COUNT(*) > 1;

-- 3. 清理角色菜单中已失效的 menu_id 引用
DELETE FROM sys_role_menu
WHERE menu_id NOT IN (SELECT id FROM sys_tenant_app_menu WHERE deleted_time = 0);

-- 4. 软删除孤儿租户菜单
UPDATE sys_tenant_app_menu tam
SET deleted_time = UNIX_TIMESTAMP(NOW(3)) * 1000
WHERE tam.deleted_time = 0
  AND NOT EXISTS (
    SELECT 1 FROM sys_module_menu mm
    WHERE mm.id = tam.template_id AND mm.deleted_time = 0
  );
