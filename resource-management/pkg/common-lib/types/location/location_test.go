package location

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRPsForRegion(t *testing.T) {
	region := Beijing
	beijingRPs := GetRPsForRegion(region)
	for i := 0; i < len(ResourcePartitions); i++ {
		assert.Equal(t, fmt.Sprintf("RP%d", i+1), string(beijingRPs[i]), "Unexpected RP name")
	}
}

func TestLocationInit(t *testing.T) {
	preLower := float64(-1)
	preUpper := float64(-1)
	for i := 0; i < len(Regions); i++ {
		region := Regions[i]
		rps := GetRPsForRegion(region)
		for j := 0; j < len(rps); j++ {
			rp := rps[j]
			loc := Location{
				region:    region,
				partition: rp,
			}
			lower, upper := loc.GetArcRangeFromLocation()
			if preLower >= lower || preUpper >= upper || lower < 0 || upper > RingRange || (preUpper > 0 && preUpper != lower) {
				assert.Fail(t, "Invalid ranges for region/resource paritions", "RP %s has unexpected hash range (%f, %f]\n\n", loc.partition, lower, upper)
				fmt.Printf("All hash range listed as follows:\n")
				printLocationRange()
				assert.Fail(t, "")
			}

			preLower = lower
			preUpper = upper
		}
	}
	fmt.Printf("All hash range listed as follows:\n")
	printLocationRange()
}

func printLocationRange() {
	for i := 0; i < len(Regions); i++ {
		region := Regions[i]
		rps := GetRPsForRegion(region)
		for j := 0; j < len(rps); j++ {
			rp := rps[j]
			loc := Location{
				region:    region,
				partition: rp,
			}
			lower, upper := loc.GetArcRangeFromLocation()
			fmt.Printf("%s, %s, [%f, %f]\n", region, rp, lower, upper)
		}
	}
}
