package constant

const (
	TimeOut          = 5
	HttpServerPort   = "9000"
	TcpServerPort    = "9001"
	AppServerPort    = "8000"
	AdminServerPort  = "8001"
	Host             = "127.0.0.1"
	TcpBufferSize    = 128
	MaxConnFailCount = 3
	HttpServName     = "app_server"
	TcpServName      = "admin_server"
	AppServName      = "app_server"
	AdminServName    = "admin_server"
	ServiceSliceCap  = 2
)

const (
	TimeOutErr    = "timeout err!"
	ListenFailed  = "listen failed!"
	ConnectFailed = "connect failed!"
	ReadFailed    = "read failed!"
)
