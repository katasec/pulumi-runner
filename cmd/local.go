package main

import (
	"os"

	pulumirunner "github.com/katasec/pulumi-runner"
)

func createLocalProgram(projectName string, workdir string) (*pulumirunner.LocalProgram, error) {
	args := &pulumirunner.LocalProgramArgs{
		ProjectName: projectName,
		StackName:   "dev",
		Plugins: []map[string]string{
			{
				"name":    "azure-native",
				"version": "v2.7.0",
			},
		},
		Config: []map[string]string{
			{
				"name":  "azure-native:location",
				"value": "westus2",
			},
		},
		WorkDir: workdir,
		Writer:  os.Stdout,
	}

	// Create a new pulumi program
	return pulumirunner.NewLocalProgram(args)
}
