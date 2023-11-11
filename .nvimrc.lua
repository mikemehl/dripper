local dap = require("dap")
dap.adapters.delve = {
	type = "server",
	host = "127.0.0.1",
	port = "34567",
	executable = {
		command = "dlv",
		args = {
			"dap",
			"--listen=127.0.0.1:34567",
			"--log",
			"--log-output=dap",
		},
	},
}

local overseer = require("overseer")
overseer
	.new_task({
		name = "build",
		cmd = "just build",
		components = {
			{ "restart_on_save" },
			{ "unique", { replace = true } },
			"default",
		},
	})
	:start()
