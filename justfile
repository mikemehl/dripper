build:
  go build

run:
  go run .

tail:
  tail -f debug.log | bat --pager=never -l log
