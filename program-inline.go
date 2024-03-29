package pulumirunner

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type InlineProgram struct {
	InlineProgramArgs
	ctx context.Context
}
type InlineProgramArgs struct {
	ProjectName string              // Name of the Pulumi project to create/destroy
	StackName   string              // Name of your pulumi stack. For e.g. "dev" or "prod"
	Plugins     []map[string]string // Plugins required for your Pulumi program, Specified as "name" and "version" in string map
	Config      []map[string]string // Config for your pulumi program, specified as "name" and "value" in string map
	Runtime     string              // Pulumi runtime for e.g. go, dotnet etc.
	Writer      io.Writer
	PulumiFn    pulumi.RunFunc // Your pulumi program you want to run, passed in as a function
	stack       auto.Stack
}

// NewInlineProgram Initalizes a stack using an inline program passed as a func
func NewInlineProgram(args *InlineProgramArgs) (*InlineProgram, error) {
	// Validate args
	validateInlineArgs(args)

	// Set stdout as default output if unspecified
	if args.Writer == nil {
		args.Writer = os.Stdout
	}
	w := args.Writer
	utils.Fprintln(w, "****** Creating New Inline Pulumi Program")
	// Initialize Context
	ctx := context.Background()

	// Create Stack using the tmp git clone folder as a local workspace
	s, err := auto.UpsertStackInlineSource(ctx, args.StackName, args.ProjectName, args.PulumiFn)
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to create or select stack: %v\n", err))
		return nil, err
	}
	args.stack = s

	// Setup Pulumi config for stack
	setConfig(w, ctx, s, args.Config)

	return &InlineProgram{
		ctx:               ctx,
		InlineProgramArgs: *args,
	}, nil
}

func (r *InlineProgram) Preview() error {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Preview")

	// Run Preview
	_, err := r.stack.Preview(r.ctx, optpreview.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Error on previewing stack: %v", err))
	}

	return err
}

func (r *InlineProgram) Up() error {

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

func (r *InlineProgram) Destroy() error {

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
