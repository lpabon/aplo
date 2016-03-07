package main

import (
	"fmt"
	// "net/http"
	// "net/url"
	"os"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/remotecommand"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
)

func main() {
	fmt.Println("Hello ams")

	config := &restclient.Config{
		Host: "http://192.168.10.90:8080",
	}
	client, err := client.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	pods, err := client.Pods(api.NamespaceDefault).List(api.ListOptions{
		LabelSelector: labels.Everything(),
		FieldSelector: fields.Everything(),
	})
	fmt.Println(pods)

	// Execute command on gluser-1
	req := client.RESTClient.Post().
		Resource("pods").
		Name("gluster-1").
		Namespace("default").
		SubResource("exec")
	req.VersionedParams(&api.PodExecOptions{
		Command: []string{"ls", "/"},
		Stdout:  true,
		Stderr:  true,
	}, api.ParameterCodec)
	fmt.Println(req.URL())
	exec, err := remotecommand.NewExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Println(err)
		return
	}

	err = exec.Stream(nil, os.Stdout, os.Stderr, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Execute command on gluser-1
	req = client.RESTClient.Post().
		Resource("pods").
		Name("gluster-1").
		Namespace("default").
		SubResource("exec")
	req.VersionedParams(&api.PodExecOptions{
		Command: []string{"gluster", "volume", "info"},
		Stdout:  true,
		Stderr:  true,
	}, api.ParameterCodec)
	fmt.Println(req.URL())
	exec, err = remotecommand.NewExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Println(err)
		return
	}

	err = exec.Stream(nil, os.Stdout, os.Stderr, false)
	if err != nil {
		fmt.Println(err)
		return
	}
}
