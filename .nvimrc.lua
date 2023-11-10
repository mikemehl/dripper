local dap = require("dap")
dap.adapters.delve = {
	type = "server",
	host = "127.0.0.1",
	port = "34567",
}
