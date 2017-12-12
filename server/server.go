package server

import "github.com/yaotian/gowechat/server/context"

//Server struct
type Server struct {
	*context.Context
}

//NewServer init
func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context
	return srv
}
