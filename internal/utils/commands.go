package utils

type Command string

const (
	FETCH        Command = "fetch"
	CLEAN        Command = "clean"
	INSTALL      Command = "install"
	BUILD        Command = "build"
	ALL          Command = "all"
	PROXY_SERVER Command = "proxy-server"
)
