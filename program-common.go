package pulumirunner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

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

	}

	return s, nil
}

func refreshStack(w io.Writer, ctx context.Context, s auto.Stack) error {
	utils.Fprintln(w, "Starting refresh")

	result, err := s.Refresh(ctx)
	if err != nil {
		utils.Fprintln(w, fmt.Sprintf("Failed to refresh stack: %v\n", err))
		return err
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

func validateInlineArgs(args *InlineProgramArgs) {

	if args.ProjectName == "" {
		exitMessage("ProjectName cannot be empty")
	}

	if args.StackName == "" {
		exitMessage("StackName cannot be empty")
	}
}

func validateLocalArgs(args *LocalProgramArgs) {

	if args.ProjectName == "" {
		exitMessage("ProjectName cannot be empty")
	}

	if args.StackName == "" {
		exitMessage("StackName cannot be empty")
	}
}

func exitMessage(message string) {
	utils.Fprintln(os.Stderr, message)
	//os.Exit(1)
}

func replaceInFile(filepath string, src string, dst string) error {
	fmt.Printf("Replacing source: %s\n Destination: %s\n File: %s\n", src, dst, filepath)
	input, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	output := bytes.Replace(input, []byte(src), []byte(dst), 1)

	if err = os.WriteFile(filepath, output, 0666); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
