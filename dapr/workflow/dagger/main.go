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

func PushGo(
	ctr *Container,
	ctx context.Context,
	imageName string,
	ch chan *ContainerIMageRef,
	isLast bool,
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
	fmt.Println("Is Last: ", isLast)
	ch <- &ContainerIMageRef{
		AppName:  imageName,
		ImageRef: ref,
	}
	if isLast {
		// time.Sleep(2 * time.Second)
		close(ch)
	}
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
// func (m *Workflow) Services(
// 	ctx context.Context,
// 	dir *Directory,
// 	appName string,
// ) (*Directory, error) {
// 	appDir := dir.WithDirectory(".", dir).Directory(appName)
// 	return appDir, nil
// }

// func (m *Workflow) BuildPush(
// 	ctx context.Context,
// 	// Service directory where the app directories are located
// 	serviceDir *Directory,
// 	// Dir path where the dapr.yaml file is located
// 	daprDir *Directory,
// 	// Path to the dapr.yaml(for k8s) file
// 	// +default="dapr-k8s.yaml"
// 	path string,
// ) (string, error) {

// 	// Return a list of app names
// 	appNames, err := m.AppName(ctx, serviceDir)
// 	if err != nil {
// 		fmt.Println("Failed to get app names")
// 		panic(err)
// 	}

// 	var refs []string

// 	// y0 := dag.Yq(daprDir)
// 	// var c0 *Container
// 	c0 := dag.Yq(daprDir).Container()

// 	// Build app Images and push to ttl.sh registry
// 	for _, appName := range appNames {
// 		appDir := serviceDir.Directory(appName)
// 		fmt.Println("App Name: ", appName)
// 		fmt.Println("App Directory: ", *appDir)

// 		ctr := m.Build(ctx, appDir)
// 		ref, err := m.Push(ctx, ctr, appName)
// 		if err != nil {
// 			fmt.Println("Failed to push image:", appName)
// 			return "", err
// 		}
// 		refs = append(refs, ref)

// 		c0 = m.UpdateDaprK8sYaml(
// 			ctx,
// 			daprDir,
// 			path,
// 			appName,
// 			ref,
// 			c0,
// 		)
// 	}
// 	fmt.Println(refs)

// 	stdout, err := c0.WithoutEntrypoint().
// 		WithExec([]string{
// 			"cat",
// 			path,
// 		}).Stdout(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	return stdout, nil

// }

// func (m *Workflow) BuildPushGo(
// 	ctx context.Context,
// 	// Service directory where the app directories are located
// 	serviceDir *Directory,
// 	// Dir path where the dapr.yaml file is located
// 	daprDir *Directory,
// 	// Path to the dapr.yaml(for k8s) file
// 	// +default="dapr-k8s.yaml"
// 	path string,
// ) (string, error) {

// 	// Return a list of app names
// 	appNames, err := m.AppName(ctx, serviceDir)
// 	if err != nil {
// 		fmt.Println("Failed to get app names")
// 		panic(err)
// 	}

// 	var refs []string

// 	// y0 := dag.Yq(daprDir)
// 	// var c0 *Container
// 	c0 := dag.Yq(daprDir).Container()

// 	var wg sync.WaitGroup
// 	wg.Add(len(appNames))

// 	// Build app Images and push to ttl.sh registry
// 	for _, appName := range appNames {
// 		appDir := serviceDir.Directory(appName)
// 		fmt.Println("App Name: ", appName)
// 		fmt.Println("App Directory: ", *appDir)

// 		go func(appName string, dir *Directory) {
// 			defer wg.Done()
// 			ctr := m.Build(ctx, appDir)
// 			ref, err := m.Push(ctx, ctr, appName)
// 			if err != nil {
// 				fmt.Println("Failed to push image:", appName)
// 				panic(err)
// 			}
// 			refs = append(refs, ref)

// 			c0 = m.UpdateDaprK8sYaml(
// 				ctx,
// 				daprDir,
// 				path,
// 				appName,
// 				ref,
// 				c0,
// 			)
// 		}(appName, appDir)
// 	}

// 	wg.Wait()
// 	fmt.Println(refs)

// 	stdout, err := c0.WithoutEntrypoint().
// 		WithExec([]string{
// 			"cat",
// 			path,
// 		}).Stdout(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	return stdout, nil

// }

// BuildPushGo builds and pushes the app images to ttl.sh registry
// and returns the updated dapr-k8s.yaml file
func (m *Workflow) BuildPushGo(
	ctx context.Context,
	// Service directory where the app directories are located
	serviceDir *Directory,
	// Dir path where the dapr.yaml file is located
	daprDir *Directory,
	// Path to the dapr.yaml(for k8s) file
	// +default="dapr-k8s.yaml"
	path string,
) (*File, error) {

	// Return a list of app names
	appNames, err := m.AppName(ctx, serviceDir)
	if err != nil {
		fmt.Println("Failed to get app names")
		panic(err)
	}

	// var refs []string

	// y0 := dag.Yq(daprDir)
	// var c0 *Container
	c0 := dag.Yq(daprDir).Container()

	ch := make(chan *ContainerIMageRef, len(appNames))
	// var once sync.Once
	// defer close(ch)

	// var wg sync.WaitGroup
	// Build app Images and push to ttl.sh registry
	for _, appName := range appNames {
		// i := i
		// wg.Add(1)
		appDir := serviceDir.Directory(appName)
		fmt.Println("App Name: ", appName)
		// isLast := (i == len(appNames)-1)
		go func(appName string, dir *Directory) {

			// Build and Push the app image
			ctr := m.Build(ctx, appDir)
			ref, err := m.Push(ctx, ctr, appName)
			// _, err := PushGo(ctr, ctx, appName, ch, isLast)

			if err != nil {
				fmt.Println("Failed to push image:", appName)
				panic(err)
			}
			// refs = append(refs, ref)

			ch <- &ContainerIMageRef{
				AppName:  appName,
				ImageRef: ref,
			}
			// if i == len(appNames)-1 {
			// 	// time.Sleep(2 * time.Second)
			// 	close(ch)
			// }
		}(appName, appDir)

	}

	// var once sync.Once
	// Wait for all images to be pushed
	for {
		if len(ch) < len(appNames) {
			fmt.Println("Waiting for all images to be pushed")
			fmt.Printf("images pushed: %v/%v", len(ch), len(appNames))
			time.Sleep(1 * time.Second)
			continue
		} else {
			fmt.Printf("images pushed: %v/%v", len(ch), len(appNames))
			close(ch)
		}
		break
	}

	for {
		c, ok := <-ch
		if !ok {
			break
		}
		c0 = m.UpdateDaprK8sYaml(
			ctx,
			daprDir,
			path,
			c.AppName,
			c.ImageRef,
			c0,
		)
		// wg.Done()
	}
	// wg.Wait()

	yamlFile := c0.
		File(path)

	return yamlFile, nil
}

func (m *Workflow) UpdateDaprK8sYaml(
	ctx context.Context,
	dir *Directory,
	path string,
	appName string,
	appImage string,
	ctr *Container,
) *Container {
	// Get the container image of the app
	exprImage1 := fmt.Sprintf(".apps[]| select(.appID == \"%s\").containerImage", appName)
	stdout, err := ctr.
		WithoutEntrypoint().
		WithExec([]string{
			"yq",
			exprImage1,
			path,
		}).Stdout(ctx)
	if err != nil {
		panic(err)
	}

	imgOld := strings.TrimSuffix(stdout, "\n")

	exprImage2 := fmt.Sprintf("s#%s#%s#", imgOld, appImage)
	c := ctr.
		// WithMountedFile(path, dir.File(path)).
		WithoutEntrypoint().
		WithExec([]string{
			"sed",
			"-i",
			exprImage2,
			path,
		}).
		WithExec([]string{
			"cat",
			path})

	return c
	// stdout, err := y1.Stdout(ctx)
	// if err != nil {
	// 	return "", err
	// }
	// return stdout, nil
}
