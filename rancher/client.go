package rancher

import (
	gorancher_client "github.com/rancher/go-rancher/client"
)

/**
 * API handlers for using rancher to provide operations
 */

// Create a new RancherClient from ClientOpts (small wrapper to simplify imports)
func MakeRancherClient(opts *gorancher_client.ClientOpts) *gorancher_client.RancherClient {
	return gorancher_client.NewRancherClient(opts)
}
