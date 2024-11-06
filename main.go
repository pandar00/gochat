package main

import (
	"github.com/pandar00/gochat/pkg/claudeai"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "gochat",
}

func main() {
	Cmd.AddCommand(claudeai.Cmd)
	Cmd.Execute()
}
