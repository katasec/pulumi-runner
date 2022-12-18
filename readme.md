# Overview 

Helper function to run a remote pulumi program. Here's an example:

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

## Output Options

By default, output is sent to `stdout`. Stream output to a file using the following to create an `io.Writer` point to a log file:

```
logger := utils.ConfigureLogger("log1.txt")
```

And use the above logger whilst creating `RemoteProgramArgs`:

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