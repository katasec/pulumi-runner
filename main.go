package main

import "fmt"

func main() {
	fmt.Println("test")

	// Create a new pulumi program
	p := createPulumiProgram()

	// Run Pulumi Up
	p.Preview()
}

func createPulumiProgram() *RemoteProgram {
	args := &RemoteProgramArgs{
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
	return NewRemoteProgram(args)
}
