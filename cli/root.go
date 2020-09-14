package cmd

import (
	// import expvar for exposing metrics and health check in the DefaultServeMux
	_ "expvar"

	"net/http"

	"github.com/ashtonian/glew/static/cmd/cli/tools"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "base-project",
}

func init() {
	RootCmd.AddCommand(apis.RootCmd)
	RootCmd.AddCommand(tools.RootCmd)

	// http://localhost:8001/debug/vars
	go http.ListenAndServe(":8001", http.DefaultServeMux) // TODO: mv to gin std so its documented and goes through gin configs debug/vars
}
