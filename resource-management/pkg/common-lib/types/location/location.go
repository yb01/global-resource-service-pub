package location

import (
	"encoding/json"
	"fmt"
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

func NewLocationFromName(regionName, partitionName string) *Location {
	region := GetRegionFromRegionName(regionName)
	partition := GetPartitionFromPartitionName(partitionName)
	if region >= 0 && partition >= 0 {
		return NewLocation(region, partition)
	}
	return nil
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
)

// later this map will be construction from config
var partitionToPartitionName = map[ResourcePartition]string{
	ResourcePartition1:  "RP1",
	ResourcePartition2:  "RP2",
	ResourcePartition3:  "RP3",
	ResourcePartition4:  "RP4",
	ResourcePartition5:  "RP5",
	ResourcePartition6:  "RP6",
	ResourcePartition7:  "RP7",
	ResourcePartition8:  "RP8",
	ResourcePartition9:  "RP9",
	ResourcePartition10: "RP10",
	ResourcePartition11: "RP11",
	ResourcePartition12: "RP12",
	ResourcePartition13: "RP13",
	ResourcePartition14: "RP14",
	ResourcePartition15: "RP15",
	ResourcePartition16: "RP16",
	ResourcePartition17: "RP17",
	ResourcePartition18: "RP18",
	ResourcePartition19: "RP19",
	ResourcePartition20: "RP20",
}

// later this map will be constructed from config
var partitionNameToPartition = map[string]ResourcePartition{
	"RP1":  ResourcePartition1,
	"RP2":  ResourcePartition2,
	"RP3":  ResourcePartition3,
	"RP4":  ResourcePartition4,
	"RP5":  ResourcePartition5,
	"RP6":  ResourcePartition6,
	"RP7":  ResourcePartition7,
	"RP8":  ResourcePartition8,
	"RP9":  ResourcePartition9,
	"RP10": ResourcePartition10,
	"RP11": ResourcePartition11,
	"RP12": ResourcePartition12,
	"RP13": ResourcePartition13,
	"RP14": ResourcePartition14,
	"RP15": ResourcePartition15,
	"RP16": ResourcePartition16,
	"RP17": ResourcePartition17,
	"RP18": ResourcePartition18,
	"RP19": ResourcePartition19,
	"RP20": ResourcePartition20,
}

func (rp ResourcePartition) String() string {
	return partitionToPartitionName[rp]
}

func GetPartitionFromPartitionName(partitionName string) ResourcePartition {
	if partition, isOK := partitionNameToPartition[partitionName]; isOK {
		return partition
	} else {
		return -1 //undefined
	}
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
	for i := 0; i < len(partitionToPartitionName); i++ {
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

// TODO - read resource parition from configuration or metadata server
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

func (loc Location) MarshalText() (text []byte, err error) {
	type l Location
	return json.Marshal(l(loc))
}

func (loc *Location) UnmarshalText(text []byte) error {
	type l Location
	return json.Unmarshal(text, (*l)(loc))
}
