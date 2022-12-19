# Overview 

Helper function to run a remote pulumi program. Check out the examples below

## Remote Program

```go
package main

func main() {
    // Create a new pulumi program
    p := createPulumiProgram()

    // Run Pulumi Up
    p.Up()

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

```

## Inline Program

```go

package main

func main() {
    // Create a new pulumi program
    p := createInlineProgram()

    // Run Pulumi Up
    p.Up()

}

func createInlineProgram(pulumiFunc pulumi.RunFunc) *pulumirunner.InlineProgram {

	homedir, _ := os.UserHomeDir()
	logger := utils.ConfigureLogger(homedir + "/ark.log")

	args := &pulumirunner.InlineProgramArgs{
		ProjectName: "ark-init",
		StackName:   "dev",
		Writer:   logger,
		PulumiFn: pulumiFunc,
	}

	return pulumirunner.NewInlineProgram(args)
}

func pulumiFunc(ctx *pulumi.Context) error {
	ctx.Export("Message", pulumi.String("Hello from Pulumi!"))
	return nil
}

```

## Output Options

By default, output is sent to `stdout`. Stream output to a file using the following to create an `io.Writer` point to a log file:

```
logger := utils.ConfigureLogger("log1.txt")
```

And use the above logger as a `Writer` whilst creating `RemoteProgramArgs`:

```go
	logger := utils.ConfigureLogger("log1.txt")
	args := &RemoteProgramArgs{
		ProjectName: "ArkInit",
		GitURL:      "https://github.com/katasec/ArkInit.git",
        .
        . //other options here
        .

		Runtime: "dotnet",
		Writer:  logger,
	}
```