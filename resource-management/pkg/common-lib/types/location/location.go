package location

import "fmt"

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
type Region string

const (
	// Regions
	Beijing   Region = "Beijing"
	Shanghai  Region = "Shanghai"
	Wulan     Region = "Wulan"
	Guizhou   Region = "Guizhou"
	Reserved1 Region = "Reserved1"
	Reserved2 Region = "Reserved2"
	Reserved3 Region = "Reserved3"
	Reserved4 Region = "Reserved4"
	Reserved5 Region = "Reserved5"
)

var Regions = []Region{
	0: Beijing,
	1: Shanghai,
	2: Wulan,
	3: Guizhou,
	4: Reserved1,
	5: Reserved2,
	6: Reserved3,
	7: Reserved4,
	8: Reserved5,
}

var regionToArc map[string]arc

// ResourcePartition defines the possible resource partition of a given node
// Defined and doced by region admin
type ResourcePartition string

const (
	ResourcePartition1  ResourcePartition = "RP1"
	ResourcePartition2  ResourcePartition = "RP2"
	ResourcePartition3  ResourcePartition = "RP3"
	ResourcePartition4  ResourcePartition = "RP4"
	ResourcePartition5  ResourcePartition = "RP5"
	ResourcePartition6  ResourcePartition = "RP6"
	ResourcePartition7  ResourcePartition = "RP7"
	ResourcePartition8  ResourcePartition = "RP8"
	ResourcePartition9  ResourcePartition = "RP9"
	ResourcePartition10 ResourcePartition = "RP10"
)

var ResourcePartitions = []ResourcePartition{ResourcePartition1, ResourcePartition2, ResourcePartition3, ResourcePartition4, ResourcePartition5,
	ResourcePartition6, ResourcePartition7, ResourcePartition8, ResourcePartition9, ResourcePartition10}
var regionRPToArc map[Location]arc

func init() {
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

func (loc *Location) GetArcRangeFromLocation() (float64, float64) {
	locArc := regionRPToArc[*loc]
	return locArc.lower, locArc.upper
}

func (loc *Location) Equal(locToCompare Location) bool {
	return loc.region == loc.region && loc.partition == locToCompare.partition
}

func (loc *Location) String() string {
	return fmt.Sprintf("[Region %s, ResoucePartition %s]", loc.region, loc.partition)
}
