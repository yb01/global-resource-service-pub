/*
Copyright The Kubernetes Authors.
Copyright 2022 Authors of Global Resource Service - file modified.
Copyright 2020 Authors of Arktos - file modified.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aggregrator

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"k8s.io/klog/v2"

	"global-resource-service/resource-management/pkg/clientSdk/rest"
	"global-resource-service/resource-management/pkg/clientSdk/watch"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	ep "global-resource-service/resource-management/pkg/service-api/endpoints"
	apiTypes "global-resource-service/resource-management/pkg/service-api/types"
	simulatorTypes "global-resource-service/resource-management/test/resourceRegionMgrSimulator/types"
)

const ResourceName = "resources"

// Config contains client config that the client can setup for
type Config struct {
	ServiceUrl     string
	RequestTimeout time.Duration
}

// ListOptions contains optional settings for List nodes
type ListOptions struct {
	// Limit is equivalent to URL query parameter ?limit=500
	Limit int
}

type SimInterface interface {
	List(ListOptions) ([][]*event.NodeEvent, types.TransitResourceVersionMap, uint64, error)
	Watch(types.TransitResourceVersionMap) (watch.Interface, error)
}

// simClient implements SimInterface
type simClient struct {
	config Config
	// REST client to region manager/simulator service
	restClient rest.Interface
}

// NewSimClient returns a reference to the simClient object
func NewSimClient(cfg Config) *simClient {
	httpclient := http.Client{Timeout: cfg.RequestTimeout}
	url, err := rest.DefaultServerURL(cfg.ServiceUrl, "", false)

	if err != nil {
		klog.Errorf("failed to get the default URL. error %v", err)
		return nil
	}

	c, err := rest.NewRESTClient(url, rest.ClientContentConfig{}, nil, &httpclient)
	if err != nil {
		klog.Errorf("failed to get the RESTClient. error %v", err)
		return nil
	}

	return &simClient{
		config:     cfg,
		restClient: c,
	}
}

// List takes label and field selectors, and returns the list of Nodes that match those selectors.
func (c *simClient) List(opts ListOptions) ([][]*event.NodeEvent, types.TransitResourceVersionMap, uint64, error) {
	req := c.restClient.Get()
	req = req.Resource(ResourceName)
	req = req.Timeout(c.config.RequestTimeout)
	req = req.Param("limit", strconv.Itoa(opts.Limit))

	respRet, err := req.DoRaw()
	if err != nil {
		return nil, nil, 0, err
	}

	resp := simulatorTypes.ResponseFromRRM{}

	err = json.Unmarshal(respRet, &resp)

	actualCrv := resp.RvMap

	return resp.RegionNodeEvents, actualCrv, resp.Length, nil

}

// Watch returns a watch.Interface that watches the requested simClient.
func (c *simClient) Watch(versionMap types.TransitResourceVersionMap) (watch.Interface, error) {
	req := c.restClient.Post()
	req = req.Resource(ResourceName)
	req = req.Timeout(c.config.RequestTimeout)
	req = req.Param(ep.WatchParameter, ep.WatchParameterTrue)

	crv := apiTypes.WatchRequest{ResourceVersions: versionMap}

	body, err := json.Marshal(crv)
	if err != nil {
		return nil, err
	}
	req = req.Body(body)

	watcher, err := req.Watch()
	if err != nil {
		return nil, err
	}

	return watcher, nil
}
