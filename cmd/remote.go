package main

import pulumirunner "github.com/katasec/pulumi-runner"

func createRemoteProgram() (*pulumirunner.RemoteProgram, error) {
	args := &pulumirunner.RemoteProgramArgs{
		ProjectName: "ArkInit",
		GitURL:      "https://github.com/katasec/ArkInit.git",
		Branch:      "refs/remotes/origin/main",
		ProjectPath: "Azure",
		StackName:   "dev",
		Plugins: []map[string]string{
			{
				"name":    "azure-native",
				"version": "v1.89.1",
			},
		},
		Config: []map[string]string{
			{
				"name":  "azure-native:location",
				"value": "westus2",
			},
		},
		Runtime: "dotnet",
	}

	// Create a new pulumi program
	return pulumirunner.NewRemoteProgram(args)
}
