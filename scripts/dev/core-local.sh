#!/usr/bin/env bash
# core-local.sh 管理 oldCore 本机联调环境。
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/../.." && pwd)"
CMD_DIR="${REPO_ROOT}/cmd"
RUN_DIR="${CMD_DIR}/run"
LOG_DIR="${CMD_DIR}/logs"
PID_FILE="${RUN_DIR}/coresvr.pid"
TUNNEL_SOCKET="${RUN_DIR}/db-tunnel.sock"
LOG_FILE="${LOG_DIR}/coresvr.log"
BIN_FILE="${CMD_DIR}/coresvr"
OSS_DIR="${CMD_DIR}/oss"
REDIS_PID_FILE="${RUN_DIR}/redis.pid"
NATS_PID_FILE="${RUN_DIR}/nats.pid"
REDIS_LOG_FILE="${LOG_DIR}/redis.log"
NATS_LOG_FILE="${LOG_DIR}/nats.log"
NATS_STORE_DIR="${RUN_DIR}/nats-jetstream"
CORE_SCREEN_SESSION="${CORE_SCREEN_SESSION:-oldcore-coresvr}"
NATS_SCREEN_SESSION="${CORE_NATS_SCREEN_SESSION:-oldcore-nats}"

REMOTE_HOST="${CORE_REMOTE_HOST:-root@new.ykhl.vip}"
REMOTE_CONFIG="${CORE_REMOTE_CONFIG:-/root/run/core/etc/common.yaml}"
LOCAL_DB_HOST="${CORE_LOCAL_DB_HOST:-127.0.0.1}"
LOCAL_DB_PORT="${CORE_LOCAL_DB_PORT:-13306}"
REMOTE_DB_HOST="${CORE_REMOTE_DB_HOST:-127.0.0.1}"
REMOTE_DB_PORT="${CORE_REMOTE_DB_PORT:-3306}"
REDIS_CONTAINER="${CORE_REDIS_CONTAINER:-oldcore-redis}"
NATS_CONTAINER="${CORE_NATS_CONTAINER:-oldcore-nats}"

# usage 输出脚本的可用命令。
usage() {
  cat <<'EOF'
Usage:
  scripts/dev/core-local.sh doctor
  scripts/dev/core-local.sh start
  scripts/dev/core-local.sh stop
  scripts/dev/core-local.sh restart
  scripts/dev/core-local.sh status
  scripts/dev/core-local.sh logs

Environment:
  CORE_REMOTE_HOST=root@new.ykhl.vip
  CORE_REMOTE_CONFIG=/root/run/core/etc/common.yaml
  CORE_LOCAL_DB_PORT=13306
  CORE_STOP_INFRA=1       stop redis/nats on stop
EOF
}

# require_cmd 确认必要命令存在。
require_cmd() {
  local name="$1"
  if ! command -v "${name}" >/dev/null 2>&1; then
    echo "missing command: ${name}" >&2
    return 1
  fi
}

# ensure_dirs 创建运行态目录。
ensure_dirs() {
  mkdir -p "${RUN_DIR}" "${LOG_DIR}" "${CMD_DIR}/etc"
}

# ensure_oss_dirs 创建本地 OSS 存储目录，避免依赖首次创建目录权限异常。
ensure_oss_dirs() {
  mkdir -p \
    "${OSS_DIR}/ithings-public" \
    "${OSS_DIR}/ithings-private" \
    "${OSS_DIR}/ithings-temporary"
}

# port_open 检测本地端口是否已监听。
port_open() {
  local host="$1"
  local port="$2"
  nc -z "${host}" "${port}" >/dev/null 2>&1
}

# pid_for_port 返回监听指定 TCP 端口的进程号。
pid_for_port() {
  local port="$1"
  lsof -nP -tiTCP:"${port}" -sTCP:LISTEN 2>/dev/null | head -n 1
}

# wait_for_port 等待本地端口可用。
wait_for_port() {
  local host="$1"
  local port="$2"
  local timeout="$3"
  local i
  for ((i = 0; i < timeout; i++)); do
    if port_open "${host}" "${port}"; then
      return 0
    fi
    sleep 1
  done
  return 1
}

# docker_container_exists 判断 Docker 容器是否存在。
docker_container_exists() {
  docker container inspect "$1" >/dev/null 2>&1
}

# docker_container_running 判断 Docker 容器是否运行中。
docker_container_running() {
  [[ "$(docker inspect -f '{{.State.Running}}' "$1" 2>/dev/null || true)" == "true" ]]
}

# docker_available 判断当前 Docker daemon 是否可用。
docker_available() {
  command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1
}

# quit_screen_session 关闭脚本创建的 detached screen 会话。
quit_screen_session() {
  local name="$1"
  screen -S "${name}" -X quit >/dev/null 2>&1 || true
}

# ensure_redis_native 使用本地 redis-server 启动 Redis。
ensure_redis_native() {
  require_cmd redis-server
  redis-server \
    --bind 127.0.0.1 \
    --port 6379 \
    --save "" \
    --appendonly no \
    --daemonize yes \
    --pidfile "${REDIS_PID_FILE}" \
    --logfile "${REDIS_LOG_FILE}" \
    --dir "${RUN_DIR}" >/dev/null
  wait_for_port 127.0.0.1 6379 20
  echo "redis: ready on 127.0.0.1:6379 (native)"
}

# ensure_redis 启动本机 Redis。
ensure_redis() {
  if port_open 127.0.0.1 6379; then
    echo "redis: already listening on 127.0.0.1:6379"
    return
  fi
  if docker_available; then
    if docker_container_exists "${REDIS_CONTAINER}"; then
      docker start "${REDIS_CONTAINER}" >/dev/null
    else
      docker run -d --name "${REDIS_CONTAINER}" -p 127.0.0.1:6379:6379 redis:7-alpine >/dev/null
    fi
    wait_for_port 127.0.0.1 6379 20
    echo "redis: ready on 127.0.0.1:6379 (docker)"
  else
    ensure_redis_native
  fi
}

# ensure_nats_native 使用本地 nats-server 启动 NATS。
ensure_nats_native() {
  require_cmd nats-server
  require_cmd screen
  mkdir -p "${NATS_STORE_DIR}"
  quit_screen_session "${NATS_SCREEN_SESSION}"
  screen -dmS "${NATS_SCREEN_SESSION}" bash -c '
    store_dir="$1"
    log_file="$2"
    exec nats-server \
      --jetstream \
      --store_dir "$store_dir" \
      --addr 127.0.0.1 \
      --port 4222 \
      --http_port 8222 >>"$log_file" 2>&1
  ' bash "${NATS_STORE_DIR}" "${NATS_LOG_FILE}"
  wait_for_port 127.0.0.1 4222 20
  pid_for_port 4222 >"${NATS_PID_FILE}"
  echo "nats: ready on 127.0.0.1:4222 (native)"
}

# ensure_nats 启动本机 NATS。
ensure_nats() {
  if port_open 127.0.0.1 4222; then
    echo "nats: already listening on 127.0.0.1:4222"
    return
  fi
  if docker_available; then
    if docker_container_exists "${NATS_CONTAINER}"; then
      docker start "${NATS_CONTAINER}" >/dev/null
    else
      docker run -d --name "${NATS_CONTAINER}" \
        -p 127.0.0.1:4222:4222 \
        -p 127.0.0.1:8222:8222 \
        nats:2-alpine -js -m 8222 >/dev/null
    fi
    wait_for_port 127.0.0.1 4222 20
    echo "nats: ready on 127.0.0.1:4222 (docker)"
  else
    ensure_nats_native
  fi
}

# fetch_remote_dsn 从远端配置读取数据库 DSN。
fetch_remote_dsn() {
  ssh -o BatchMode=yes -o ConnectTimeout=8 "${REMOTE_HOST}" \
    "awk '/^[[:space:]]*DSN:/{sub(/^[[:space:]]*DSN:[[:space:]]*/, \"\"); print; exit}' '${REMOTE_CONFIG}'"
}

# localize_dsn 将远端 DSN 的 TCP 地址替换成本地 SSH 隧道地址。
localize_dsn() {
  local dsn="$1"
  local local_addr="${LOCAL_DB_HOST}:${LOCAL_DB_PORT}"
  local rewritten
  rewritten="$(printf '%s\n' "${dsn}" | sed -E "s/@tcp\\([^)]*\\)/@tcp(${local_addr})/")"
  if [[ "${rewritten}" == "${dsn}" ]]; then
    echo "remote DSN does not contain @tcp(...)" >&2
    return 1
  fi
  printf '%s\n' "${rewritten}"
}

# ensure_db_tunnel 建立本机到远端数据库的 SSH 隧道。
ensure_db_tunnel() {
  if port_open "${LOCAL_DB_HOST}" "${LOCAL_DB_PORT}"; then
    echo "db tunnel: already listening on ${LOCAL_DB_HOST}:${LOCAL_DB_PORT}"
    return
  fi
  rm -f "${TUNNEL_SOCKET}"
  ssh -M -S "${TUNNEL_SOCKET}" -fnNT \
    -L "${LOCAL_DB_HOST}:${LOCAL_DB_PORT}:${REMOTE_DB_HOST}:${REMOTE_DB_PORT}" \
    -o ExitOnForwardFailure=yes \
    -o ServerAliveInterval=30 \
    -o ServerAliveCountMax=3 \
    "${REMOTE_HOST}"
  wait_for_port "${LOCAL_DB_HOST}" "${LOCAL_DB_PORT}" 10
  echo "db tunnel: ready on ${LOCAL_DB_HOST}:${LOCAL_DB_PORT}"
}

# stop_db_tunnel 关闭由脚本创建的 SSH 隧道。
stop_db_tunnel() {
  if [[ -S "${TUNNEL_SOCKET}" ]]; then
    ssh -S "${TUNNEL_SOCKET}" -O exit "${REMOTE_HOST}" >/dev/null 2>&1 || true
    rm -f "${TUNNEL_SOCKET}"
  fi
}

# build_core 编译本机 coresvr。
build_core() {
  (cd "${REPO_ROOT}" && go build -tags no_k8s -o "${BIN_FILE}" ./service/apisvr)
}

# sync_etc 同步运行所需配置文件。
sync_etc() {
  rsync -a --delete "${REPO_ROOT}/service/apisvr/etc/" "${CMD_DIR}/etc/"
}

# core_pid 返回当前记录的 coresvr 进程号。
core_pid() {
  if [[ -f "${PID_FILE}" ]]; then
    cat "${PID_FILE}"
  elif port_open 127.0.0.1 7777; then
    pid_for_port 7777
  fi
}

# core_running 判断当前记录的 coresvr 是否仍在运行。
core_running() {
  local pid
  pid="$(core_pid || true)"
  [[ -n "${pid}" ]] && kill -0 "${pid}" >/dev/null 2>&1
}

# stop_core 停止本机 coresvr。
stop_core() {
  local pid
  pid="$(core_pid || true)"
  quit_screen_session "${CORE_SCREEN_SESSION}"
  if [[ -n "${pid}" ]] && kill -0 "${pid}" >/dev/null 2>&1; then
    kill "${pid}" >/dev/null 2>&1 || true
    local i
    for ((i = 0; i < 10; i++)); do
      if ! kill -0 "${pid}" >/dev/null 2>&1; then
        break
      fi
      sleep 1
    done
    if kill -0 "${pid}" >/dev/null 2>&1; then
      kill -9 "${pid}" >/dev/null 2>&1 || true
    fi
  fi
  rm -f "${PID_FILE}"
}

# start_core 启动本机 coresvr。
start_core() {
  local db_dsn="$1"
  stop_core
  : > "${LOG_FILE}"
  screen -dmS "${CORE_SCREEN_SESSION}" bash -c '
    workdir="$1"
    log_file="$2"
    shift 2
    cd "$workdir"
    exec env "$@" ./coresvr core >>"$log_file" 2>&1
  ' bash "${CMD_DIR}" "${LOG_FILE}" \
      confSuffix=Local \
      CORE_DB_DSN="${db_dsn}" \
      CORE_LOCAL_SKIP_STARTUP_SIDE_EFFECTS=1 \
      CORE_LOCAL_SKIP_AUTO_MIGRATE=1
  if ! wait_for_port 127.0.0.1 7777 30; then
    echo "coresvr did not open 127.0.0.1:7777" >&2
    tail -80 "${LOG_FILE}" >&2 || true
    return 1
  fi
  pid_for_port 7777 >"${PID_FILE}"
  echo "coresvr: ready on 0.0.0.0:7777"
}

# doctor 检查本机联调依赖是否满足。
doctor() {
  require_cmd go
  require_cmd ssh
  require_cmd nc
  require_cmd curl
  require_cmd rsync
  if docker_available; then
    echo "infra: docker"
  else
    require_cmd redis-server
    require_cmd nats-server
    require_cmd screen
    echo "infra: native redis-server/nats-server"
  fi
  ssh -o BatchMode=yes -o ConnectTimeout=8 "${REMOTE_HOST}" "test -r '${REMOTE_CONFIG}'"
  if [[ -z "$(fetch_remote_dsn)" ]]; then
    echo "remote DSN not found in ${REMOTE_CONFIG}" >&2
    return 1
  fi
  echo "doctor: OK"
}

# start 启动完整本地联调链路。
start() {
  ensure_dirs
  doctor
  ensure_redis
  ensure_nats
  ensure_db_tunnel
  local remote_dsn
  local local_dsn
  remote_dsn="$(fetch_remote_dsn)"
  local_dsn="$(localize_dsn "${remote_dsn}")"
  build_core
  sync_etc
  ensure_oss_dirs
  start_core "${local_dsn}"
  curl -fsS http://127.0.0.1:7777/api/v1/system/common/debug >/dev/null
  echo "debug endpoint: OK"
}

# stop 停止本地联调链路。
stop() {
  ensure_dirs
  stop_core
  stop_db_tunnel
  if [[ "${CORE_STOP_INFRA:-}" == "1" ]]; then
    if docker_available; then
      docker stop "${REDIS_CONTAINER}" >/dev/null 2>&1 || true
      docker stop "${NATS_CONTAINER}" >/dev/null 2>&1 || true
    fi
    if [[ -f "${REDIS_PID_FILE}" ]]; then
      kill "$(cat "${REDIS_PID_FILE}")" >/dev/null 2>&1 || true
      rm -f "${REDIS_PID_FILE}"
    fi
    if [[ -f "${NATS_PID_FILE}" ]]; then
      quit_screen_session "${NATS_SCREEN_SESSION}"
      kill "$(cat "${NATS_PID_FILE}")" >/dev/null 2>&1 || true
      rm -f "${NATS_PID_FILE}"
    fi
  fi
  echo "stopped"
}

# status 输出本地联调链路状态。
status() {
  ensure_dirs
  if core_running; then
    echo "coresvr: running pid $(core_pid)"
  else
    echo "coresvr: stopped"
  fi
  for port in "${LOCAL_DB_PORT}" 6379 4222 7777; do
    if port_open 127.0.0.1 "${port}"; then
      echo "port ${port}: listening"
    else
      echo "port ${port}: closed"
    fi
  done
  if docker_available; then
    if docker_container_running "${REDIS_CONTAINER}"; then
      echo "redis container: running"
    fi
    if docker_container_running "${NATS_CONTAINER}"; then
      echo "nats container: running"
    fi
  fi
  if [[ -f "${REDIS_PID_FILE}" ]]; then
    echo "redis native pid: $(cat "${REDIS_PID_FILE}")"
  fi
  if [[ -f "${NATS_PID_FILE}" ]]; then
    echo "nats native pid: $(cat "${NATS_PID_FILE}")"
  fi
}

# logs 跟随 coresvr 日志。
logs() {
  ensure_dirs
  touch "${LOG_FILE}"
  tail -f "${LOG_FILE}"
}

case "${1:-}" in
  doctor)
    doctor
    ;;
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    stop
    start
    ;;
  status)
    status
    ;;
  logs)
    logs
    ;;
  -h|--help|"")
    usage
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac
