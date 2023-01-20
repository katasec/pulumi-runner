package pulumirunner

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
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

	Stack auto.Stack
}

// NewRemoteProgram Initalizes a RemoteProgram. This clones a remote Git repository
// in to a local folder and sets it up as a Pulumi workspace
func NewRemoteProgram(args *RemoteProgramArgs) (*RemoteProgram, error) {

	// Validate args
	validateRemoteArgs(args)

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
		return nil, err
	}
	args.Stack = s

	// Setup Pulumi config for stack
	setConfig(w, ctx, s, args.Config)

	return &RemoteProgram{
		ctx:               ctx,
		RemoteProgramArgs: *args,
	}, nil
}

func (r *RemoteProgram) Up() error {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Up")
	// Refresh before Pulumi Up
	refreshStack(w, r.ctx, r.Stack)

	// Run Up
	_, err := r.Stack.Up(r.ctx, optup.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to update stack: %v", err))
	} else {
		utils.Fprintln(w, "Stack successfully updated")
	}

	return err
}

func (r *RemoteProgram) Preview() error {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Preview")

	// Run Preview
	_, err := r.Stack.Preview(r.ctx, optpreview.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to update stack: %v", err))
	} else {
		utils.Fprintln(w, "Stack successfully updated")
	}

	return err
}

func (r *RemoteProgram) Destroy() error {

	// Get writer for logging
	w := r.Writer
	utils.Fprintln(w, "****** Starting Pulumi Destroy")
	// Refresh before Pulumi Destroy
	refreshStack(w, r.ctx, r.Stack)

	// Run destroy
	_, err := r.Stack.Destroy(r.ctx, optdestroy.ProgressStreams(w))
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to destroy stack: %v", err))
	} else {
		utils.Fprintln(w, "Stack successfully destroyed")
	}

	return err
}

func (r *RemoteProgram) CleanUp() {
	log.Println("Removing Temporary remote folder")
	os.RemoveAll(r.Stack.Workspace().WorkDir())
}

func (r *RemoteProgram) GetCfgFile() string {

	// Get stack name and workdir
	stack := r.RemoteProgramArgs.Stack.Name()
	workDir := r.RemoteProgramArgs.Stack.Workspace().WorkDir()

	// Generate pulumi yaml file name with stack
	fName := fmt.Sprintf("Pulumi.%s.yaml", stack)

	// Return full path to yaml config file
	cfgFile := path.Join(workDir, fName)

	return cfgFile
}

// FixConfig Cleans up an extraneous pipe symbol from the Pulumi yaml config file. This is a hack to allow
// injection of yaml data received from Ark MQ versus a string literal, as an input parameter to a remote
// pulumi program. Yaml data from Ark is stored as 'arkdata' in the config file
func (r *RemoteProgram) FixConfig() {
	file := r.GetCfgFile()
	replaceInFile(file, "arkdata: |", "arkdata:")
}

func RemoveDir(dir string) {
	os.RemoveAll(dir)
}
