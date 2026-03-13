#!/usr/bin/env bash
set -euo pipefail

PORT=3000

usage() {
  echo "Usage: $0 {start|stop|restart|status}"
  exit 1
}

get_pid() {
  lsof -ti :"$PORT" 2>/dev/null || true
}

kill_port() {
  local pid
  pid=$(get_pid)
  if [ -n "$pid" ]; then
    echo "Killing process on port $PORT (PID: $pid)"
    kill -9 $pid 2>/dev/null || true
    sleep 0.5
  fi
}

do_start() {
  local pid
  pid=$(get_pid)
  if [ -n "$pid" ]; then
    echo "Port $PORT is occupied (PID: $pid), killing it..."
    kill_port
  fi
  echo "Starting dev server on port $PORT..."
  nohup npx vite --port "$PORT" --host > /tmp/adb-claw-dev.log 2>&1 &
  disown
  sleep 1
  local new_pid
  new_pid=$(get_pid)
  if [ -n "$new_pid" ]; then
    echo "Dev server started (PID: $new_pid), log: /tmp/adb-claw-dev.log"
  else
    echo "Failed to start. Check /tmp/adb-claw-dev.log"
    exit 1
  fi
}

do_stop() {
  local pid
  pid=$(get_pid)
  if [ -n "$pid" ]; then
    kill_port
    echo "Stopped."
  else
    echo "No process on port $PORT."
  fi
}

do_status() {
  local pid
  pid=$(get_pid)
  if [ -n "$pid" ]; then
    echo "Running on port $PORT (PID: $pid)"
  else
    echo "Not running."
  fi
}

cd "$(dirname "$0")"

case "${1:-}" in
  start)   do_start ;;
  stop)    do_stop ;;
  restart) do_stop; do_start ;;
  status)  do_status ;;
  *)       usage ;;
esac
