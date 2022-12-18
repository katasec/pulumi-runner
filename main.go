package pulumirunner

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	fmt.Println("test")

	// Create a new pulumi program
	p := createInlineProgram()

	// Run Pulumi Up
	p.Preview()
}

func createRemoteProgram() *RemoteProgram {
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

func createInlineProgram() *InlineProgram {
	args := &InlineProgramArgs{
		ProjectName: "Inline",
		StackName:   "dev",
		Runtime:     "go",
		PulumiFn:    inlineFunc,
	}

	// Create a new pulumi program
	return NewInlineProgram(args)
}

func inlineFunc(ctx *pulumi.Context) error {

	ctx.Export("message", pulumi.String("hello world"))

	return nil
}
