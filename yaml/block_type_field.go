package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

func BlockTypeField(blockType types.BlockType) string {
	switch blockType {
	case types.BlockTypeApplication:
		return "apps"
	case types.BlockTypeDomain:
		return "domains"
	case types.BlockTypeSubdomain:
		return "subdomains"
	case types.BlockTypeIngress:
		return "ingresses"
	case types.BlockTypeDatastore:
		return "datastores"
	case types.BlockTypeClusterNamespace:
		return "cluster_namespaces"
	case types.BlockTypeCluster:
		return "clusters"
	case types.BlockTypeNetwork:
		return "networks"
	default:
		return "blocks"
	}
}
