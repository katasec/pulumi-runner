package pulumirunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

type RemoteProgram struct {
	RemoteProgramArgs
	ctx context.Context
}

type RemoteProgramArgs struct {
	ProjectName string              // Name of the Pulumi project to create/destroy
	GitURL      string              // For example github.com/katasec/project.git
	Branch      string              // Git branch of the repo to checkout
	Tag         string              // Git tag
	ProjectPath string              // A sub folder under  github.com/katasec/project.git for e.g. folder1
	StackName   string              // Name of your pulumi stack. For e.g. "dev" or "prod"
	Plugins     []map[string]string // Plugins required for your Pulumi program, Specified as "name" and "version" in string map
	Config      []map[string]string // Config for your pulumi program, specified as "name" and "value" in string map
	Runtime     string              // Pulumi runtime for e.g. go, dotnet etc.
	Writer      io.Writer

	stack auto.Stack
}

// NewRemoteProgram Initalizes a RemoteProgram. This clones a remote Git repository
// in to a local folder and sets it up as a Pulumi workspace
func NewRemoteProgram(args *RemoteProgramArgs) *RemoteProgram {

	// Validate args
	validateArgs(args)

	// Set stdout as default output if unspecified
	if args.Writer == nil {
		args.Writer = os.Stdout
	}
	w := args.Writer
	utils.Fprintln(w, "****** Creating New Remote Pulumi Program")
	// Initialize Context
	ctx := context.Background()

	// Clone Git repo to a temp folder
	utils.Fprintln(w, "Cloning repo: "+args.GitURL)
	workDir := utils.CloneRemote(w, args.GitURL)

	// Create Stack using the tmp git clone folder as a local workspace
	s, err := auto.UpsertStackLocalSource(ctx, args.StackName, filepath.Join(workDir, args.ProjectPath))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to create or select stack: %v\n", err))
		os.Exit(1)
	}
	args.stack = s

	// Setup Pulumi config for stack
	setConfig(w, ctx, s, args.Config)

	return &RemoteProgram{
		ctx:               ctx,
		RemoteProgramArgs: *args,
	}
}

func setConfig(w io.Writer, ctx context.Context, s auto.Stack, config []map[string]string) (auto.Stack, error) {
	// Set stack config if specified:
	if config != nil {
		// set stack configuration using name/value from map
		for _, key := range config {
			err := s.SetConfig(ctx, key["name"], auto.ConfigValue{Value: key["value"]})
			if err != nil {
				return s, err
			}
		}

		utils.Fprintln(w, "Successfully set config")
	}

	return s, nil
}

func refreshStack(w io.Writer, ctx context.Context, s auto.Stack) error {
	utils.Fprintln(w, "Starting refresh")

	result, err := s.Refresh(ctx)
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to refresh stack: %v\n", err))
		os.Exit(1)
	}

	utils.Fprintln(w, fmt.Sprintf("Refresh succeeded!, Result:%s \n", result.Summary.Result))

	return nil
}

func validateArgs(args *RemoteProgramArgs) {

	if args.ProjectName == "" {
		exitMessage("ProjectName cannot be empty")
	}

	if args.GitURL == "" {
		exitMessage("GitURL cannot be empty")
	}

	if args.StackName == "" {
		exitMessage("StackName cannot be empty")
	}
}

func exitMessage(message string) {
	utils.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func (r *RemoteProgram) Up() {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Up")
	// Refresh before Pulumi Up
	refreshStack(w, r.ctx, r.stack)

	// Run destroy
	_, err := r.stack.Up(r.ctx, optup.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to update stack: %v", err))
		os.Exit(1)
	} else {
		utils.Fprintln(w, "Stack successfully updated")
	}
}

func (r *RemoteProgram) Destroy() {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Destroy")
	// Refresh before Pulumi Destroy
	refreshStack(w, r.ctx, r.stack)

	// Run destroy
	_, err := r.stack.Destroy(r.ctx, optdestroy.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to destroy stack: %v", err))
		os.Exit(1)
	} else {
		utils.Fprintln(w, "Stack successfully destroyed")
	}
}
