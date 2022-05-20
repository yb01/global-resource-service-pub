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

func (r Region) String() string {
	switch r {
	case Beijing:
		return "Beijing"
	case Shanghai:
		return "Shanghai"
	case Wulan:
		return "Wulan"
	case Guizhou:
		return "Guizhou"
	case Reserved1:
		return "Reserved1"
	case Reserved2:
		return "Reserved2"
	case Reserved3:
		return "Reserved3"
	case Reserved4:
		return "Reserved4"
	case Reserved5:
		return "Reserved5"
	}
	return "undefined"
}

var Regions = []Region{Beijing, Shanghai, Wulan, Guizhou, Reserved1, Reserved2, Reserved3, Reserved4, Reserved5}

var regionToArc map[string]arc

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

func (rp ResourcePartition) String() string {
	switch rp {
	case ResourcePartition1:
		return "RP1"
	case ResourcePartition2:
		return "RP2"
	case ResourcePartition3:
		return "RP3"
	case ResourcePartition4:
		return "RP4"
	case ResourcePartition5:
		return "RP5"
	case ResourcePartition6:
		return "RP6"
	case ResourcePartition7:
		return "RP7"
	case ResourcePartition8:
		return "RP8"
	case ResourcePartition9:
		return "RP9"
	case ResourcePartition10:
		return "RP10"
	case ResourcePartition11:
		return "RP11"
	case ResourcePartition12:
		return "RP12"
	case ResourcePartition13:
		return "RP13"
	case ResourcePartition14:
		return "RP14"
	case ResourcePartition15:
		return "RP15"
	case ResourcePartition16:
		return "RP16"
	case ResourcePartition17:
		return "RP17"
	case ResourcePartition18:
		return "RP18"
	case ResourcePartition19:
		return "RP19"
	case ResourcePartition20:
		return "RP20"
	}
	return "undefined"
}

var ResourcePartitions = []ResourcePartition{ResourcePartition1, ResourcePartition2, ResourcePartition3, ResourcePartition4, ResourcePartition5,
	ResourcePartition6, ResourcePartition7, ResourcePartition8, ResourcePartition9, ResourcePartition10, ResourcePartition11, ResourcePartition12,
	ResourcePartition13, ResourcePartition14, ResourcePartition15, ResourcePartition16, ResourcePartition17, ResourcePartition18, ResourcePartition19,
	ResourcePartition20}

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
	return loc.region == loc.region && loc.partition == locToCompare.partition
}

func (loc *Location) String() string {
	return fmt.Sprintf("[Region %s, ResoucePartition %s]", loc.region, loc.partition)
}
