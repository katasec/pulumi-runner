package pulumirunner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

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

		// Clean up pipe symbol from the Pulumi yaml config file
		// This is an workaround for a bug
		stack := s.Name()
		workDir := s.Workspace().WorkDir()
		fName := fmt.Sprintf("Pulumi.%s.yaml", stack)
		cfgFile := path.Join(workDir, fName)

		replaceInFile(cfgFile, "arkdata: |", "arkdata:")
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

func validateRemoteArgs(args *RemoteProgramArgs) {

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

func validateLocalArgs(args *InlineProgramArgs) {

	if args.ProjectName == "" {
		exitMessage("ProjectName cannot be empty")
	}

	if args.StackName == "" {
		exitMessage("StackName cannot be empty")
	}
}

func exitMessage(message string) {
	utils.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func replaceInFile(filepath string, src string, dst string) {
	fmt.Printf("Replace %s with % s in file: %s", src, dst, filepath)
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte(src), []byte(dst), -1)

	if err = ioutil.WriteFile(filepath, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
