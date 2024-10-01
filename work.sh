work() {
  # usage: work 10m 42, work 10m on task 42. Default is 25m

  duration="${1:-25m}"
  task_id="$2"

  if [ -n "$task_id" ]; then
    tomatillo activate --id="$task_id"
  fi

  timer "$duration" && echo -e "\a"
  # Update the tomatillo task if a task id is provided
  if [ -n "$task_id" ]; then
    tomatillo update --id="$task_id"
  fi
}

rest() {
  # usage: rest 10m, rest 60s etc. Default is 5m
  timer "${1:-5m}" && echo -e "\a"
}