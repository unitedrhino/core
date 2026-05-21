SELECT column_name, column_type
FROM information_schema.columns
WHERE table_schema = 'iThings'
  AND table_name = 'sys_tenant_app'
  AND column_name LIKE 'google_%';
