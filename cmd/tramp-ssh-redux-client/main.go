package main

import (
	"github.com/raymond-w-ko/tramp-ssh-redux/pkg/utils"
	"github.com/raymond-w-ko/tramp-ssh-redux/internal/client"
)

func main() {
	err := utils.ChdirToSelfExecutablePath()
	if (err != nil) {
		utils.FatalMessageAsJson("Could not chdir to self executable path", err)
	}
	client.SetupConnectionToProxy()
	cmd := client.ReadInitialCommandFromArgs()
	client.SendCommandToProxyAndWriteOutputToStdout(cmd)
}
