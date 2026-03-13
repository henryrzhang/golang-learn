#!/bin/bash
# 测试数据脚本：注册用户 -> 获取 token -> 插入 20 条剧集数据
# 依赖：curl（jq 可选，用于解析 JSON）
# 用法：./scripts/seed_test_data.sh [BASE_URL]

set -e
BASE_URL="${1:-http://localhost:8080}"
API="${BASE_URL}/api"

EMAIL="test@example.com"
PASSWORD="test123456"
NAME="测试用户"
PHONE="13800138000"

echo "=== 1. 注册/登录测试用户 ==="
REGISTER_RESP=$(curl -s -w "\n%{http_code}" -X POST "${API}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"${NAME}\",\"email\":\"${EMAIL}\",\"phone\":\"${PHONE}\",\"password\":\"${PASSWORD}\"}")

HTTP_CODE=$(echo "$REGISTER_RESP" | tail -n1)
BODY=$(echo "$REGISTER_RESP" | sed '$d')

if [ "$HTTP_CODE" = "409" ]; then
  echo "用户已存在，尝试登录..."
  LOGIN_RESP=$(curl -s -X POST "${API}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}")
  BODY="$LOGIN_RESP"
fi

# 解析 token（兼容 jq 或 grep）
if command -v jq &>/dev/null; then
  TOKEN=$(echo "$BODY" | jq -r '.token // .body.token // empty')
else
  TOKEN=$(echo "$BODY" | grep -oE '"token"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed 's/"token"[^"]*"\([^"]*\)".*/\1/')
fi
if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "获取 token 失败，响应: $BODY"
  exit 1
fi
echo "Token 获取成功"

echo ""
echo "=== 2. 插入 20 条剧集测试数据 ==="
for i in $(seq 1 20); do
  DRAMA_NO="DRAMA-$(printf '%03d' $i)"
  TITLE="测试剧集 #${i}"
  OUTLINE="这是第 ${i} 条测试剧集的大纲内容。"
  TASK_NO="TASK-$(printf '%04d' $i)"

  RESP=$(curl -s -w "\n%{http_code}" -X POST "${API}/drama" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${TOKEN}" \
    -d "{
      \"drama_no\": \"${DRAMA_NO}\",
      \"title\": \"${TITLE}\",
      \"outline\": \"${OUTLINE}\",
      \"characters\": \"\",
      \"character_relation_desc\": \"\",
      \"task_no\": \"${TASK_NO}\"
    }")

  CODE=$(echo "$RESP" | tail -n1)
  if [ "$CODE" = "200" ] || [ "$CODE" = "201" ]; then
    echo "  [${i}/20] ${DRAMA_NO} - ${TITLE} ✓"
  else
    echo "  [${i}/20] ${DRAMA_NO} 失败 (HTTP $CODE)"
  fi
done

echo ""
echo "=== 完成 ==="
