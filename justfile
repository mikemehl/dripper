build:
  go build

run:
  go run .

watch:
  git ls-files "**/*.go" | fzz watch -n "just build && just run"

tail:
  tail -f debug.log | bat --pager=never -l log

dap:
  dlv dap --listen=127.0.0.1:34567 --log-output=dap
