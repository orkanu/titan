package utils

type Action string

const (
	FETCH        Action = "fetch"
	CLEAN        Action = "clean"
	INSTALL      Action = "install"
	BUILD        Action = "build"
	ALL          Action = "all"
	PROXY_SERVER Action = "proxy-server"
)
