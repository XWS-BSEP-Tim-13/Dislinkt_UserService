package main

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/startup"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
