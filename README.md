# tomatillo
A pomodoro app with good features

## usage

```
add         Add a new task
    --name
    --estimate
update      Update the actual pomodoros of a task
    --id
done        Mark a task as done
    --id
edit        Edit the estimate of a task
report      Generate a report
    --type monthly
    --type yearly
    --type weekly (default)
delete      Delete a task
    --id 
load        Load tasks from a file
    --file
simple      Generate a simple report of todays work
version     Print the version of the application
```


Add a task

```bash
tomatillo add --name "Prepare the widget" --estimate 2

tomatillo add -n "Add another task" -e 4
```

This uses the tomatillo binary, a sqlite database and a timer app. Install the timer app with

```bash
brew install caarlos0/tap/timer
```

which is a simple and reliable timer

Add this to your zshrc

```bash
work() {
  # usage: work 10m 42, work 10m on task 42. Default is 25m

  duration="${1:-25m}"
  task_id="$2"

  timer "$duration" && echo -e "\a"
  # Update the tomatillo task if a task id is provided
  if [ -n "$task_id" ]; then
    go run . update --id="$task_id"
  fi
}

rest() {
  # usage: rest 10m, rest 60s etc. Default is 5m
  timer "${1:-5m}" && echo -e "\a"
}
```

And run a command such as:

```bash
work 25m 192
```

And it will update your current task. 

Which means work for 25 minutes on task number 192

### Generate a weekly report

How productive have you been this week? 

```bash
tomatillo report
```

```
Weekly Report
Day         | Done | Actual
------------|------|----------------
2024-09-08  | 0    |
2024-09-09  | 3    | 🍅🍅🍅🍅🍅🍅🍅🍅🍅🍅
2024-09-10  | 2    | 🍅🍅🍅🍅🍅🍅🍅
2024-09-11  | 0    |
2024-09-12  | 0    |
2024-09-13  | 0    |
2024-09-14  | 0    |
```

How about for the year? 

```bash
tomatillo report --type yearly
```

```
July
Su Mo Tu We Th Fr Sa
   --  3  4 --  4 10
 1 --  2 --  3  2 --
 1 -- -- -- -- --  4
 2  6  5  5 -- -- --
 2  2  3  1

August
Su Mo Tu We Th Fr Sa
            --  1 --
 3  1  3  1  5 --  1
 1 -- -- --  7  5  6
 2  7  2  2 --  1 --
 2 --  5  4 -- -- --


September
Su Mo Tu We Th Fr Sa
--  3 --  3 -- -- --
-- 10  7 -- -- -- --
-- -- -- -- -- -- --
-- -- -- -- -- -- --
-- --
```


## local testing

Testing the build pipeline by running `act` to simulate the Github Actions workflow

```bash
act -j test
```

## Building

Build the binary and move it to a location in your path

```bash
go build -o tomatillo

mv ./tomatillo ~/go/bin/
```



