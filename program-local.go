package pulumirunner

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

// LocalProgram A struct for running pulumi programs in the local file system
type LocalProgram struct {
	LocalProgramArgs
	ctx context.Context
}
type LocalProgramArgs struct {
	ProjectName string              // Name of the Pulumi project to create/destroy
	StackName   string              // Name of your pulumi stack. For e.g. "dev" or "prod"
	Plugins     []map[string]string // Plugins required for your Pulumi program, Specified as "name" and "version" in string map
	Config      []map[string]string // Config for your pulumi program, specified as "name" and "value" in string map
	Runtime     string              // Pulumi runtime for e.g. go, dotnet etc.
	Writer      io.Writer
	WorkDir     string // Path to the pulumi program in your local file system
	stack       auto.Stack
}

// NewLocalProgram Initalizes a stack using a pulimi program in the local file system
func NewLocalProgram(args *LocalProgramArgs) (*LocalProgram, error) {
	// Validate args
	validateLocalArgs(args)

	// Set stdout as default output if unspecified
	if args.Writer == nil {
		args.Writer = os.Stdout
	}
	w := args.Writer
	utils.Fprintln(w, "****** Creating New Inline Pulumi Program")
	// Initialize Context
	ctx := context.Background()

	// Create Stack using the tmp git clone folder as a local workspace
	s, err := auto.UpsertStackLocalSource(ctx, args.StackName, args.WorkDir)
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to create or select stack: %v\n", err))
		return nil, err
	}
	args.stack = s

	// Setup Pulumi config for stack
	setConfig(w, ctx, s, args.Config)

	return &LocalProgram{
		ctx:              ctx,
		LocalProgramArgs: *args,
	}, nil
}

func (r *LocalProgram) Up() error {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Up")

	// Run Preview
	_, err := r.stack.Up(r.ctx, optup.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to update stack: %v", err))
	} else {
		utils.Fprintln(w, "Stack successfully updated")
	}

	return err
}

func (r *LocalProgram) Destroy() error {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Destroy")

	// Run Preview
	_, err := r.stack.Destroy(r.ctx, optdestroy.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to destroy stack: %v", err))
	} else {
		utils.Fprintln(w, "Stack successfully destroyed")
	}

	return err
}
