package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"time"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/core/client"
	"github.com/deref/util-go/jsonutil"
)

var binPath = "/usr/local/bin/mitmweb"
var scriptPath = "/home/mitmproxy/routing.py"

type HostAndPort struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type HostMap map[string]HostAndPort

type HostMapper interface {
	GetHostMap() HostMap
}

type StaticHostMapper HostMap

func (shm StaticHostMapper) GetHostMap() HostMap {
	return HostMap(shm)
}

func writeScript(hostMap HostMap) error {
	script := fmt.Sprintf("mapping = %s\n", jsonutil.MustMarshalString(hostMap)) + `
import urllib.parse
class Rerouter:
    def request(self, flow):
        originalHostHeader = flow.request.host_header
        host = urllib.parse.urlsplit('//' + originalHostHeader).hostname
        if host == "exo.localhost":
            flow.request.host = "localhost"
            flow.request.port = 8081
            return

        dest = mapping.get(host, None)
        if dest:
            print(dest)
            flow.request.host = dest["host"]
            flow.request.port = int(dest["port"])
            flow.request.host_header = originalHostHeader

    def response(self, flow):
        flow.response.headers["x-frame-options"] = "ALLOW"

addons = [Rerouter()]
`
	fmt.Printf("script: %+v\n", script)
	return ioutil.WriteFile(scriptPath, []byte(script), 0600)
}

func run(hostMapper HostMapper) {
	hostMap := hostMapper.GetHostMap()
	if err := writeScript(hostMap); err != nil {
		panic(err)
	}
	cmd := exec.Command(binPath, "--web-host=0.0.0.0", "--set=keep_host_header", fmt.Sprintf("-s=%q", scriptPath))
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}()
	for {
		newHostMap := hostMapper.GetHostMap()
		fmt.Printf("newHostMap: %+v\n", newHostMap)
		if !reflect.DeepEqual(hostMap, newHostMap) {
			hostMap = newHostMap
			if err := writeScript(hostMap); err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second)
	}
}

func mustGetEnv(varName string) string {
	val, isSet := os.LookupEnv(varName)
	if !isSet {
		panic(varName + " is not set")
	}
	return val
}

type ExoHostMapper struct {
	mu        *sync.Mutex
	exoClient *client.Root
	hostMap   HostMap
}

func (ehm *ExoHostMapper) GetHostMap() HostMap {
	ehm.mu.Lock()
	defer ehm.mu.Unlock()
	return ehm.hostMap
}

func (ehm *ExoHostMapper) RefreshHostMap() {
	ws := ehm.exoClient.GetWorkspace(mustGetEnv("EXO_WORKSPACE_ID"))
	endpoints, err := ws.GetServiceEndpoints(context.TODO(), &api.GetServiceEndpointsInput{})
	if err != nil {
		fmt.Println("WARN: could not get service endpoint: %w", err)
		return
	}
	hostMap := HostMap{}
	for _, endpoint := range endpoints.ServiceEndpoints {
		hostMap[fmt.Sprintf("%s.exo.localhost", endpoint.Service)] = HostAndPort{
			Host: "host.docker.internal",
			Port: endpoint.Port,
		}
	}
	ehm.mu.Lock()
	defer ehm.mu.Unlock()
	ehm.hostMap = hostMap
}

func main() {
	exoHostMapper := &ExoHostMapper{
		mu:      &sync.Mutex{},
		hostMap: HostMap{},
		exoClient: &client.Root{
			HTTP:  http.DefaultClient,
			URL:   mustGetEnv("EXO_URL"),
			Token: mustGetEnv("EXO_TOKEN"),
		},
	}

	go func() {
		for {
			exoHostMapper.RefreshHostMap()
			time.Sleep(time.Second)
		}
	}()

	run(exoHostMapper)
}
