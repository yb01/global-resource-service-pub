package aggregrator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	interfaces "global-resource-service/resource-management/pkg/common-lib/interfaces/distributor"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

type Aggregator struct {
	urls           []string
	EventProcessor interfaces.InterfacesOfDistributor
}

// To be client of Resource Region Manager
// RRM: Resource Region Manager
//
type ClientOfRRM struct {
	BaseURL    string
	HTTPClient *http.Client
}

// RRM: Resource Region Manager
//
type ResponseFromRRM struct {
	MinRecordNodeEvents []*event.NodeEvent
	RvMap               types.ResourceVersionMap
	Length              int
}

// RRM: Resource Region Manager
//
type PullDataFromRRM struct {
	BatchLength int
	CRV         types.ResourceVersionMap
}

const (
	DefaultBatchLength = 1000
)

// Initialize aggregator
//
func NewAggregator(urls []string, EventProcessor interfaces.InterfacesOfDistributor) *Aggregator {
	return &Aggregator{
		urls:           urls,
		EventProcessor: EventProcessor,
	}
}

// Main loop to get resources from resource region managers and send to distributor
// This is only initial code structure for aggregator method
// TODO:
//    Based on the speed of process events from resource distributor, dynamic descision of batch length
//    will be made by aggregator and this batch length will be used to pull resources from resource region manager
//
func (a *Aggregator) Run() {
	numberOfURLs := len(a.urls)

	fmt.Println("Running for loop to connect to to resource region manager...")

	for i := 0; i < numberOfURLs; i++ {
		go func(i int) {
			var crv types.ResourceVersionMap
			var minRecordNodeEvents []*event.NodeEvent

			// Connect to resource region manager
			c := a.createClient(a.urls[i])

			for {
				time.Sleep(100 * time.Millisecond)

				// Call the Pull methods
				// when composite RV is nil, the method initPull is called;
				// otherwise the method subsequentPull is called.
				// To simplify the codes, we use one method initPullOrSubsequentPull instead
				minRecordNodeEvents, _ = a.initPullOrSubsequentPull(c, DefaultBatchLength, crv)

				if minRecordNodeEvents != nil {
					// Call ProcessEvents() and get the CRV from distributor as default success
					// TODO:
					//    1. wrap up the logical Node record per Distributor needs ( the interfaces need to be updated as well )
					//    2. call the ProcessEvents Per RP to unload some cost from the Distributor
					eventProcess, crv := a.EventProcessor.ProcessEvents(minRecordNodeEvents)

					// Call resource region manager, POST CRV to release old node events when ProcessEvents is successful
					if eventProcess {
						a.postCRV(c, crv)
					}
				}
			}
		}(i)
	}

	fmt.Println("Finished for loop to connect to to resource region manager...")
}

// Connect to resource region manager
//
func (a *Aggregator) createClient(url string) *ClientOfRRM {
	return &ClientOfRRM{
		BaseURL: url,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// Call resource region manager's InitPull method {url}/resources/initpull when crv is nil
// or
// Call the resource region manager's SubsequentPull method {url}/resources/subsequentpull when crv is not nil
//
func (a *Aggregator) initPullOrSubsequentPull(c *ClientOfRRM, batchLength int, crv types.ResourceVersionMap) ([]*event.NodeEvent, int) {
	var path string

	if crv == nil {
		path = "c.baseURL/resources/initPull"
	} else {
		path = "c.baseURL/resources/subsequentPull"
	}

	bytes, _ := json.Marshal(PullDataFromRRM{BatchLength: batchLength, CRV: crv.Copy()})
	req, err := http.NewRequest(http.MethodGet, path, strings.NewReader((string(bytes))))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var ResponseObject ResponseFromRRM
	err = json.Unmarshal(bodyBytes, &ResponseObject)
	if err != nil {
		fmt.Println("Error from JSON Unmarshal:", err)
	}

	return ResponseObject.MinRecordNodeEvents, ResponseObject.Length
}

// Call resource region manager's POST method {url}/resources/crv to update the CRV
// error indicate failed POST, CRV means Composite Resource Version
//
func (a *Aggregator) postCRV(c *ClientOfRRM, crv types.ResourceVersionMap) error {
	path := "c.baseURL/resources/crv"
	bytes, _ := json.Marshal(crv.Copy())
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader((string(bytes))))

	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	_, err = c.HTTPClient.Do(req)

	return err
}
