package main

import (
	"github.com/pkg/errors"
)

const (
	kubeconfigPath = "run-k8s-job-kubeconfig"
)

var (
	errNoAuth = errors.New("you must provide either 'kubeconfig-file' or both 'cluster-url' and 'cluster-token'")
)

type ActionInput struct {
	kubeconfigFile  string
	image           string
	jobName         string
	namespace       string
	dsn             string
	migrationSource string
}
