package terraform

import (
	"github.com/hashicorp/terraform/internal/addrs"
	"github.com/hashicorp/terraform/internal/configs"
	"github.com/hashicorp/terraform/internal/dag"
)

// NestedTransformer is a final check that handles the special cases devoted
// to nested data blocks.
//
// For now, we only have nested data blocks within check blocks, and we just
// need to make sure these are the last things that get executed.
type NestedTransformer struct {
	// Config for the entire module.
	Config *configs.Config
}

func (s NestedTransformer) Transform(graph *Graph) error {

	var resources []dag.Vertex
	var checkBlockDataSources []dag.Vertex

	for _, vertex := range graph.Vertices() {
		if node, isResource := vertex.(GraphNodeConfigResource); isResource {
			addr := node.ResourceAddr()

			if node.ResourceAddr().Resource.Mode == addrs.ManagedResourceMode {
				resources = append(resources, node)
				continue
			}

			config := s.Config
			if !addr.Module.IsRoot() {
				config = s.Config.Descendent(addr.Module)
			}

			resource := config.Module.ResourceByAddr(addr.Resource)
			if resource != nil && resource.Container != nil {
				if _, ok := resource.Container.(*configs.Check); ok {
					checkBlockDataSources = append(checkBlockDataSources, node)
				}
			}
		}
	}

	for _, resource := range resources {
		for _, nested := range checkBlockDataSources {
			graph.Connect(dag.BasicEdge(nested, resource))
		}
	}

	return nil
}
