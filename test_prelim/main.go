package main

import (
	"context"
 	"fmt"
	"log"
 	"time"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
 	"github.com/containerd/containerd/namespaces"
)

func main() {
	err := redisExample()
 	if err != nil {
 		log.Fatal(err);
 	}
}

func redisExample() error {
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	// create a new context with an "example" namespace
	ctx := namespaces.WithNamespace(context.Background(), "example")

	// pull the redis image from DockerHub
	image, err := client.Pull(ctx, "docker.io/library/redis:latest", containerd.WithPullUnpack)
	if err != nil {
		return err
	}

	// // generate an OCI runtime spec using the Args, Env, etc from the redis image that we pulled
	// spec, err := containerd.GenerateSpec(containerd.WithImageConfig(ctx, image))
	// if err != nil {
	// 	return err
	// }

	// create a container
	container, err := client.NewContainer(
		ctx,
		"redis-server",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("redis-rootfs", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)


	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	
	// create a task from the container
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	defer task.Delete(ctx)

	// make sure we wait before calling start
	time.Sleep(3 * time.Second)

	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		return err
	}

	time.Sleep(3*time.Second)

	fmt.Printf("Starting checkpoint")

	start := time.Now()

	// checkpoint the task then push it to a registry
	checkpoint, err := task.Checkpoint(ctx)

	err = client.Push(ctx, "myregistry/checkpoints/redis:master", checkpoint)

	elapsed := time.Since(start)
	log.Printf("Elapsed: %s", elapsed)
	fmt.Printf("Checkpoint completed")

	// // sleep for a lil bit to see the logs
	// time.Sleep(3 * time.Second)

	// // kill the process and get the exit status
	// if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
	// 	return err
	// }

	// wait for the process to fully exit and print out the exit status
	// status := <-exitStatusC
	// fmt.Printf("redis-server exited with status: %d\n", status)
	return nil
}
