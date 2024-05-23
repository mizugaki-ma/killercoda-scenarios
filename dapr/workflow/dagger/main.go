// A generated module for Workflow functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
)

type Workflow struct{}

// Returns a container that echoes whatever string argument is provided
func (m *Workflow) ContainerEcho(stringArg string) *Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Workflow) GrepDir(
	ctx context.Context,
	directoryArg *Directory,
	pattern string,
) (string, error) {

	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

func (m *Workflow) Build(
	ctx context.Context,
	buildSrc *Directory,
) *Container {

	builder := dag.Container().
		From("python:3-alpine").
		WithDirectory("/app", buildSrc).
		WithWorkdir("/app").
		WithFile("/app/requirements.txt", buildSrc.File("requirements.txt")).
		WithExec([]string{"pip", "install", "-r", "requirements.txt"})

	runner := dag.Container().
		From("python:3-alpine").
		WithDirectory("/app", buildSrc).
		WithDirectory("/usr/local/lib/python3.12/site-packages", builder.Directory("/usr/local/lib/python3.12/site-packages")).
		WithFile(".", buildSrc.File("app.py")).
		WithEntrypoint([]string{"python3", "/app/app.py"})

	return runner
}

func (m *Workflow) Push(
	ctx context.Context,
	ctr *Container,
	imageName string,
) (string, error) {

	ref, err := ctr.Publish(
		ctx,
		fmt.Sprintf(
			"ttl.sh/%s-%.0f",
			imageName,
			math.Floor(rand.Float64()*10000000),
		),
	)
	if err != nil {
		fmt.Println("Failed to push image:", ref)
		panic(err)
	}
	fmt.Println("Successfully pushed image:", ref)
	return ref, nil
}

// Get the list of app names in the directory
func (m *Workflow) AppName(
	ctx context.Context,
	dir *Directory,
) ([]string, error) {
	apps, err := dir.Entries(ctx)
	if err != nil {
		fmt.Printf("failed to list directory entries: %v", err)
	}
	return apps, nil
}

// Get the directory of the app
func (m *Workflow) Services(
	ctx context.Context,
	dir *Directory,
	appName string,
) (*Directory, error) {
	appDir := dir.WithDirectory(".", dir).Directory(appName)
	return appDir, nil
}

func (m *Workflow) BuildPush(
	ctx context.Context,
	serviceDir *Directory,
) ([]string, error) {

	// Return a list of app names
	appNames, err := m.AppName(ctx, serviceDir)
	if err != nil {
		fmt.Println("Failed to get app names")
		panic(err)
	}

	// // Return a list of app directory objects
	// appDirs, err := m.Services(ctx, serviceDir)
	// if err != nil {
	// 	fmt.Println("Failed to get app directories")
	// 	panic(err)
	// }

	var refs []string
	// Build app Images and push to ttl.sh registry
	for _, appName := range appNames {
		appDir := serviceDir.Directory(appName)
		fmt.Println("App Name: ", appName)
		fmt.Println("App Directory: ", *appDir)

		ctr := m.Build(ctx, appDir)
		ref, err := m.Push(ctx, ctr, appName)
		if err != nil {
			fmt.Println("Failed to push image:", appName)
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}
