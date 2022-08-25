/*
Copyright 2022 Authors of Global Resource Service.

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

package location

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const RingRange = float64(360)

type Location struct {
	region    Region
	partition ResourcePartition
}

func NewLocation(region Region, partition ResourcePartition) *Location {
	return &Location{
		region:    region,
		partition: partition,
	}
}

type arc struct {
	lower float64
	upper float64
}

// Region defines the possible region location of a given node
// Defined and doced by region admin
type Region int

const (
	// Regions
	Beijing   Region = 0
	Shanghai  Region = 1
	Wulan     Region = 2
	Guizhou   Region = 3
	Reserved1 Region = 4
	Reserved2 Region = 5
	Reserved3 Region = 6
	Reserved4 Region = 7
	Reserved5 Region = 8
)

// later this map will be construction from config
var regionToRegionName = map[Region]string{
	Beijing:   "Beijing",
	Shanghai:  "Shanghai",
	Wulan:     "Wulan",
	Guizhou:   "Guizhou",
	Reserved1: "Reserved1",
	Reserved2: "Reserved2",
	Reserved3: "Reserved3",
	Reserved4: "Reserved4",
	Reserved5: "Reserved5",
}

// later this map will be construction from config
var regionNameToRegion = map[string]Region{
	"Beijing":   Beijing,
	"Shanghai":  Shanghai,
	"Wulan":     Wulan,
	"Guizhou":   Guizhou,
	"Reserved1": Reserved1,
	"Reserved2": Reserved2,
	"Reserved3": Reserved3,
	"Reserved4": Reserved4,
	"Reserved5": Reserved5,
}

func (r Region) String() string {
	return regionToRegionName[r]
}

func GetRegionFromRegionName(regionName string) Region {
	if r, isOK := regionNameToRegion[regionName]; isOK {
		return r
	} else {
		return -1 // undefined
	}
}

var Regions = []Region{}

// ResourcePartition defines the possible resource partition of a given node
// Defined and doced by region admin
type ResourcePartition int

// later this const list will be constructed from config
const (
	ResourcePartition1  ResourcePartition = 0
	ResourcePartition2  ResourcePartition = 1
	ResourcePartition3  ResourcePartition = 2
	ResourcePartition4  ResourcePartition = 3
	ResourcePartition5  ResourcePartition = 4
	ResourcePartition6  ResourcePartition = 5
	ResourcePartition7  ResourcePartition = 6
	ResourcePartition8  ResourcePartition = 7
	ResourcePartition9  ResourcePartition = 8
	ResourcePartition10 ResourcePartition = 9
	ResourcePartition11 ResourcePartition = 10
	ResourcePartition12 ResourcePartition = 11
	ResourcePartition13 ResourcePartition = 12
	ResourcePartition14 ResourcePartition = 13
	ResourcePartition15 ResourcePartition = 14
	ResourcePartition16 ResourcePartition = 15
	ResourcePartition17 ResourcePartition = 16
	ResourcePartition18 ResourcePartition = 17
	ResourcePartition19 ResourcePartition = 18
	ResourcePartition20 ResourcePartition = 19
	ResourcePartition21 ResourcePartition = 20
	ResourcePartition22 ResourcePartition = 21
	ResourcePartition23 ResourcePartition = 22
	ResourcePartition24 ResourcePartition = 23
	ResourcePartition25 ResourcePartition = 24
	ResourcePartition26 ResourcePartition = 25
	ResourcePartition27 ResourcePartition = 26
	ResourcePartition28 ResourcePartition = 27
	ResourcePartition29 ResourcePartition = 28
	ResourcePartition30 ResourcePartition = 29
	ResourcePartition31 ResourcePartition = 30
	ResourcePartition32 ResourcePartition = 31
	ResourcePartition33 ResourcePartition = 32
	ResourcePartition34 ResourcePartition = 33
	ResourcePartition35 ResourcePartition = 34
	ResourcePartition36 ResourcePartition = 35
	ResourcePartition37 ResourcePartition = 36
	ResourcePartition38 ResourcePartition = 37
	ResourcePartition39 ResourcePartition = 38
	ResourcePartition40 ResourcePartition = 39

	ResourcePartitionMax = 40
)

func (rp ResourcePartition) String() string {
	return rp.GetPartitionName()
}

func (rp ResourcePartition) GetPartitionName() string {
	return fmt.Sprintf("RP%d", int(rp)+1)
}

func GetPartitionFromPartitionName(partitionName string) (ResourcePartition, error) {
	values := strings.Split(partitionName, "RP")
	if len(values) != 2 || values[0] != "" {
		return -1, errors.New("Invalid resource partition") //undefined
	}

	num, err := strconv.Atoi(values[1])
	if err != nil || num > ResourcePartitionMax || num < 1 {
		return -1, errors.New("Invalid resource partition") //undefined
	}
	return ResourcePartition(num - 1), nil
}

var ResourcePartitions = []ResourcePartition{}

var regionToArc map[string]arc
var regionRPToArc map[Location]arc

func init() {
	// initialize Regions
	for i := 0; i < len(regionToRegionName); i++ {
		Regions = append(Regions, Region(i))
	}
	// initialize ResourcePartitions
	for i := 0; i < ResourcePartitionMax; i++ {
		ResourcePartitions = append(ResourcePartitions, ResourcePartition(i))
	}

	regionRPToArc = make(map[Location]arc)
	regionGrain := RingRange / float64(len(Regions))

	regionLower := float64(0)
	regionUpper := regionGrain

	for i := 0; i < len(Regions); i++ {
		region := Regions[i]
		rps := GetRPsForRegion(region)

		rpLower := regionLower
		rpGrain := regionGrain / float64(len(rps))
		rpUpper := regionLower + rpGrain
		for j := 0; j < len(rps); j++ {
			rp := rps[j]
			loc := Location{
				region:    region,
				partition: rp,
			}
			if j == len(rps)-1 {
				if i == len(Regions)-1 {
					regionRPToArc[loc] = arc{lower: rpLower, upper: RingRange}
				} else {
					regionRPToArc[loc] = arc{lower: rpLower, upper: regionUpper}
				}
			} else {
				regionRPToArc[loc] = arc{lower: rpLower, upper: rpUpper}
				rpLower = rpUpper
				rpUpper += rpGrain
			}
		}

		regionLower = regionUpper
		regionUpper += regionGrain
	}
}

func GetRegionNum() int {
	return len(Regions)
}

func GetRPNum() int {
	return len(ResourcePartitions)
}

// TODO - read resource partition from configuration or metadata server
func GetRPsForRegion(region Region) []ResourcePartition {
	rpsForRegion := make([]ResourcePartition, len(ResourcePartitions))
	for i := 0; i < len(ResourcePartitions); i++ {
		rpsForRegion[i] = ResourcePartitions[i]
	}
	return rpsForRegion
}

func (loc *Location) GetRegion() Region {
	return loc.region
}

func (loc *Location) GetResourcePartition() ResourcePartition {
	return loc.partition
}

func (loc *Location) GetArcRangeFromLocation() (float64, float64) {
	locArc := regionRPToArc[*loc]
	return locArc.lower, locArc.upper
}

func (loc *Location) Equal(locToCompare Location) bool {
	return loc.region == locToCompare.region && loc.partition == locToCompare.partition
}

func (loc *Location) String() string {
	return fmt.Sprintf("[Region %s, ResoucePartition %s]", loc.region, loc.partition)
}
