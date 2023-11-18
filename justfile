build:
  go build

run:
  go run .

watch:
  git ls-files "**/*.go" | fzz watch -n "just build && just run"

tail:
  tail -f dripper.log | bat --pager=never -l log

tail-destructive: clear-log tail

clear-log:
  echo "" > dripper.log

dap:
  dlv dap --listen=127.0.0.1:34567 --log-output=dap
