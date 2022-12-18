package main

import (
	"io"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

// NewRemoteProgram Initalizes a RemoteProgram. This clones a remote Git repository
// in to a local folder and sets it up as a Pulumi workspace
func NewLocalProgram(args *RemoteProgramArgs) *RemoteProgram {
	return nil
}

type LocalProgramArgs struct {
	ProjectName string              // Name of the Pulumi project to create/destroy
	StackName   string              // Name of your pulumi stack. For e.g. "dev" or "prod"
	Plugins     []map[string]string // Plugins required for your Pulumi program, Specified as "name" and "version" in string map
	Config      []map[string]string // Config for your pulumi program, specified as "name" and "value" in string map
	Runtime     string              // Pulumi runtime for e.g. go, dotnet etc.
	Writer      io.Writer

	stack auto.Stack
}
