#!/bin/bash
# 清除所有用户 token（Redis 在 Docker 中）

CONTAINER=$(docker ps --format '{{.Names}}' | grep -i redis | head -1)

if [ -z "$CONTAINER" ]; then
  echo "未找到 Redis 容器，请手动指定容器名："
  echo "  CONTAINER=redis-container $0"
  exit 1
fi

echo "使用 Redis 容器: $CONTAINER"

COUNT=$(docker exec "$CONTAINER" redis-cli KEYS "userToken:*" | wc -l)
echo "找到 $COUNT 个 token"

if [ "$COUNT" -gt 0 ]; then
  docker exec "$CONTAINER" redis-cli KEYS "userToken:*" | xargs docker exec -i "$CONTAINER" redis-cli DEL
  echo "✓ 已清除所有用户 token"
else
  echo "没有需要清除的 token"
fi
