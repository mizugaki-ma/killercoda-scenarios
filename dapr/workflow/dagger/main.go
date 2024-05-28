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
	"strings"
	"time"
)

type Workflow struct{}

type ContainerIMageRef struct {
	AppName  string
	ImageRef string
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

func (m *Workflow) UpdateDaprK8sYaml(
	ctx context.Context,
	daprFile *File,
	appName string,
	appImage string,
	ctr *Container,
) {
	// Get the container image of the app
	exprImage1 := fmt.Sprintf(".apps[]| select(.appID == \"%s\").containerImage", appName)

	filePath, _ := daprFile.Name(ctx)
	stdout, err := ctr.
		WithoutEntrypoint().
		WithExec([]string{
			"yq",
			exprImage1,
			filePath,
		}).Stdout(ctx)
	if err != nil {
		panic(err)
	}
	appImageOld := strings.TrimSuffix(stdout, "\n")

	// Update the container image of the app
	exprImage2 := fmt.Sprintf("s#%s#%s#", appImageOld, appImage)
	ctr.
		// WithMountedFile(path, dir.File(path)).
		WithoutEntrypoint().
		WithExec([]string{
			"sed",
			"-i",
			exprImage2,
			filePath,
		}).
		WithExec([]string{
			"cat",
			filePath})
}

// BuildPushGo builds and pushes the app images to ttl.sh registry
// and returns the updated dapr-k8s.yaml file
func (m *Workflow) BuildPushGo(
	ctx context.Context,
	// Service directory where the app directories are located
	serviceDir *Directory,
	// Dir path where the dapr.yaml file is located
	daprDir *Directory,
	// Path to the template of dapr.yaml file (for k8s use case)
	// +default="dapr-k8s-tpl.yaml"
	templatePath string,
) (*File, error) {

	// Return a list of app names
	appNames, err := serviceDir.Entries(ctx)
	if err != nil {
		fmt.Println("Failed to get app names")
		panic(err)
	}

	// Create a channel to receive the image refs
	ch := make(chan *ContainerIMageRef, len(appNames))

	// Build app Images and push to ttl.sh registry concurrently
	for _, appName := range appNames {
		appDir := serviceDir.Directory(appName)
		go func(appName string, appDir *Directory) {
			// Build and Push the app image
			ctr := m.Build(ctx, appDir)
			ref, err := m.Push(ctx, ctr, appName)
			if err != nil {
				fmt.Println("Failed to push image:", appName)
				panic(err)
			}
			ch <- &ContainerIMageRef{
				AppName:  appName,
				ImageRef: ref,
			}
		}(appName, appDir)
	}

	// Create a new Container to update the dapr.yaml file
	c0 := dag.Yq(daprDir).Container()
	daprYaml := daprDir.File(templatePath)

	// Wait for all images to be pushed
	for {
		if len(ch) < len(appNames) {
			fmt.Println("Waiting for all images to be pushed")
			fmt.Printf("image pushed: %v/%v", len(ch), len(appNames))
			time.Sleep(1 * time.Second)
			continue
		} else {
			fmt.Printf("images pushed: %v/%v", len(ch), len(appNames))
			close(ch)
		}
		break
	}

	// Update the dapr.yaml file sequentially
	for {
		c, ok := <-ch
		if !ok {
			break
		}

		m.UpdateDaprK8sYaml(
			ctx,
			daprYaml,
			c.AppName,
			c.ImageRef,
			c0,
		)
	}

	// Return the updated dapr.yaml file
	return c0.File(templatePath), nil
}
