package compose

import (
	"fmt"
	"io"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/yamlutil"
)

type Loader struct {
	// ProjectName is used as a prefix for the resources created by this importer.
	ProjectName string
}

func (i *Loader) Load(r io.Reader) manifest.LoadResult {
	composeProject, err := Parse(r)
	if err != nil {
		return manifest.LoadResult{Err: fmt.Errorf("parsing: %w", err)}
	}
	return i.convert(composeProject)
}

func (i *Loader) convert(project *Project) manifest.LoadResult {
	res := manifest.LoadResult{
		Manifest: &manifest.Manifest{},
	}

	// Since containers reference networks and volumes by their docker-compose name, but the
	// Docker components will have a namespaced name, so we need to keep track of which
	// volumes/components a service references.
	networksByComposeName := map[string]string{}
	volumesByComposeName := map[string]string{}

	for originalName, volume := range project.Volumes {
		name := manifest.MangleName(originalName)
		if originalName != name {
			res = res.AddRenameWarning(originalName, name)
		}

		v := volume.ToMap()
		if _, ok := v["name"]; !ok {
			v["name"] = i.prefixedName(name)
		}
		volumesByComposeName[originalName] = v["name"].(string)

		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "volume",
			Spec: yamlutil.MustMarshalString(v),
		})
	}

	// Set up networks.
	hasDefaultNetwork := false
	for originalName, network := range project.Networks {
		if originalName == "default" {
			hasDefaultNetwork = true
		}
		name := manifest.MangleName(originalName)
		if originalName != name {
			res = res.AddRenameWarning(originalName, name)
		}
		n := network.ToMap()

		// If a `name` key is specified in the network configuration (usually used in conjunction with `external: true`),
		// then we should honor that as the docker network name. Otherwise, we should set the name as
		// `<project_name>_<network_key>`.
		if _, ok := n["name"]; !ok {
			n["name"] = i.prefixedName(name)
		}
		networksByComposeName[originalName] = n["name"].(string)

		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(n),
		})
	}
	if !hasDefaultNetwork {
		componentName := "default"
		name := i.prefixedName(componentName)
		networksByComposeName[componentName] = name

		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: componentName,
			Type: "network",
			Spec: yamlutil.MustMarshalString(map[string]string{
				"name":   name,
				"driver": "bridge",
			}),
		})
	}

	for originalName, service := range project.Services {
		name := manifest.MangleName(originalName)
		if originalName != name {
			res = res.AddRenameWarning(originalName, name)
		}
		s := service.ToMap()
		if _, ok := s["container_name"]; !ok {
			s["container_name"] = i.prefixedName(name)
		}

		// Map the docker-compose network name to the name of the docker network that is created.
		defaultNetworkName := networksByComposeName["default"]
		if networks, ok := s["networks"]; ok {
			switch typedNetworks := networks.(type) {
			case map[string]interface{}:
				mappedNetworks := make(map[string]interface{})
				for network, config := range typedNetworks {
					networkName, ok := networksByComposeName[network]
					if !ok {
						res.Err = fmt.Errorf("unknown network: %q", network)
						continue
					}
					mappedNetworks[networkName] = config
				}
				mappedNetworks[defaultNetworkName] = map[interface{}]interface{}{}
				s["networks"] = mappedNetworks

			case []string:
				mappedNetworks := make([]string, len(typedNetworks)+1)
				for i, network := range typedNetworks {
					networkName, ok := networksByComposeName[network]
					if !ok {
						res.Err = fmt.Errorf("unknown network: %q", network)
						continue
					}
					mappedNetworks[i] = networkName
				}
				mappedNetworks[len(mappedNetworks)-1] = defaultNetworkName
				s["networks"] = mappedNetworks
			}
		} else {
			s["networks"] = []string{defaultNetworkName}
		}

		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "container",
			Spec: yamlutil.MustMarshalString(s),
		})
	}

	return res
}

func (i *Loader) prefixedName(name string) string {
	return fmt.Sprintf("%s_%s", i.ProjectName, name)
}
