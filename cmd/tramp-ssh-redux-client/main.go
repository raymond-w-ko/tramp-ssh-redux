package main

import (
	"github.com/raymond-w-ko/tramp-ssh-redux/pkg/utils"
	"github.com/raymond-w-ko/tramp-ssh-redux/internal/client"
)

func main() {
	utils.ChdirToSelfExecutablePath()
	client.SetupConnectionToProxy()
}
