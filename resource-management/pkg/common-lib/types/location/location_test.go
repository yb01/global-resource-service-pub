package location

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
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

/*
Get location range 10000 times in 1.498767ms - 10K
Get location range 100000 times in 15.064011ms - 100K
Get location range 1000000 times in 150.8075ms - 1M
Get location range 10000000 times in 1.513191225s - 10M
*/
func TestGetLocationRangeByStruct_Performance(t *testing.T) {
	count := []int{10000, 100000, 1000000, 10000000}
	regionNum := len(Regions)
	rpNum := len(ResourcePartitions)
	for i := 0; i < len(count); i++ {
		start := time.Now()
		for j := 0; j < count[i]; j++ {
			region := Regions[rand.Int()%regionNum]
			partition := ResourcePartitions[rand.Int()%rpNum]
			loc := NewLocation(region, partition)
			loc.GetArcRangeFromLocation()
		}
		duration := time.Since(start)
		fmt.Printf("Get location range %d times in %v\n", count[i], duration)
	}
}

/*
=== RUN   TestGetLocationRangeByPointer_Performance
Get location range 10000 times in 1.336781ms
Get location range 100000 times in 13.266604ms
Get location range 1000000 times in 132.641325ms
Get location range 10000000 times in 1.327625666s
--- PASS: TestGetLocationRangeByPointer_Performance (1.47s)
*/
func TestGetLocationRangeByPointer_Performance(t *testing.T) {
	locNum := len(Regions) * len(ResourcePartitions)
	locations := make([]*Location, locNum)
	for i := 0; i < len(Regions); i++ {
		for j := 0; j < len(ResourcePartitions); j++ {
			pos := i*len(ResourcePartitions) + j
			locations[pos] = NewLocation(Regions[i], ResourcePartitions[j])
		}
	}

	count := []int{10000, 100000, 1000000, 10000000}
	for i := 0; i < len(count); i++ {
		start := time.Now()
		for j := 0; j < count[i]; j++ {
			pos := rand.Int() % locNum
			loc := locations[pos]
			loc.GetArcRangeFromLocation()
		}
		duration := time.Since(start)
		fmt.Printf("Get location range %d times in %v\n", count[i], duration)
	}
}

/*
=== RUN   TestGetLocationRangeByObject_Performance
Get location range 10000 times in 1.242111ms
Get location range 100000 times in 12.574774ms
Get location range 1000000 times in 125.588971ms
Get location range 10000000 times in 1.25763665s
--- PASS: TestGetLocationRangeByObject_Performance (1.40s)
*/
func TestGetLocationRangeByObject_Performance(t *testing.T) {
	locNum := len(Regions) * len(ResourcePartitions)
	locations := make([]Location, locNum)
	for i := 0; i < len(Regions); i++ {
		for j := 0; j < len(ResourcePartitions); j++ {
			pos := i*len(ResourcePartitions) + j
			locations[pos] = Location{
				region:    Regions[i],
				partition: ResourcePartitions[j],
			}
		}
	}

	count := []int{10000, 100000, 1000000, 10000000}
	for i := 0; i < len(count); i++ {
		start := time.Now()
		for j := 0; j < count[i]; j++ {
			pos := rand.Int() % locNum
			loc := locations[pos]
			loc.GetArcRangeFromLocation()
		}
		duration := time.Since(start)
		fmt.Printf("Get location range %d times in %v\n", count[i], duration)
	}
}
