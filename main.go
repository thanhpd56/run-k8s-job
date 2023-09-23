package main

import (
	"context"
	"flag"
	"github.com/sethvargo/go-githubactions"
	"k8s.io/client-go/tools/clientcmd"
	"time"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	action := githubactions.New()

	input := ActionInput{
		kubeconfigFile: action.GetInput("kubeconfig-file"),
		namespace:      action.GetInput("namespace"),
		image:          action.GetInput("image"),
		jobName:        action.GetInput("job-name"),
	}

	flag.Parse()
	action.Debugf("kubeconfig input %s\n", input.kubeconfigFile)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", input.kubeconfigFile)
	if err != nil {
		action.Fatalf("%v", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		action.Fatalf("%v", err)
	}

	action = action.WithFieldsMap(map[string]string{
		"job": input.jobName,
	})

	namespace := input.namespace
	jobs := clientset.BatchV1().Jobs(namespace)
	pods := clientset.CoreV1().Pods(namespace)
	runner := NewJobRunner(jobs, pods, 5*time.Second, action)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	logs, err := runner.RunJob(ctx, input.jobName, namespace, input.image)
	defer cancel()

	if err != nil {
		if len(logs) == 0 {
			action.Fatalf("%v", err)
		} else {
			action.Fatalf("job failed\njob logs:\n%s", logs)
		}
	}

	action.Debugf("job completed successfully\njob logs:\n%s", logs)
}
