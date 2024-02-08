package main

import (
	pulumirunner "github.com/katasec/pulumi-runner"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createInlineProgram() (*pulumirunner.InlineProgram, error) {
	args := &pulumirunner.InlineProgramArgs{
		ProjectName: "Inline",
		StackName:   "dev",
		Runtime:     "go",
		PulumiFn:    inlineFunc,
	}

	// Create a new pulumi program
	return pulumirunner.NewInlineProgram(args)
}

func inlineFunc(ctx *pulumi.Context) error {

	ctx.Export("message", pulumi.String("hello world"))

	return nil
}
