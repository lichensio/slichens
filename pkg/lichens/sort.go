package lichens

import (
	"github.com/lichensio/slichens/pkg/statistics"
	"sort"
)

type LessFunc func(p1, p2 *SurveyKey) int

type MultiSorter struct {
	changes []SurveyKey
	less    []LessFunc
}

func (ms *MultiSorter) Sort(changes []SurveyKey) {
	ms.changes = changes
	sort.Sort(ms)
	defer func() { ms.changes = nil }() // Clear the slice after sorting
}

func OrderedBy(less ...LessFunc) *MultiSorter {
	return &MultiSorter{
		less: less,
	}
}

func (ms *MultiSorter) Len() int {
	return len(ms.changes)
}

func (ms *MultiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

func (ms *MultiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	for _, less := range ms.less {
		val := less(p, q)
		if val < 0 {
			return true
		} else if val > 0 {
			return false
		}
		// If val == 0, try the next comparison
	}
	return false // Default in case all comparisons indicate equality
}

// NewMultiSorter creates a new multiSorter with the provided less functions.
func NewMultiSorter(less ...LessFunc) *MultiSorter {
	return &MultiSorter{
		less: less,
	}
}

func GetSortFunctions(freq bool) []LessFunc {
	netname := func(c1, c2 *SurveyKey) int {
		if c1.NetName < c2.NetName {
			return -1
		} else if c1.NetName > c2.NetName {
			return 1
		}
		return 0
	}

	cellid := func(c1, c2 *SurveyKey) int {
		if c1.CellID < c2.CellID {
			return -1
		} else if c1.CellID > c2.CellID {
			return 1
		}
		return 0
	}

	band := func(c1, c2 *SurveyKey) int {
		if c1.Band < c2.Band {
			return -1
		} else if c1.Band > c2.Band {
			return 1
		}
		return 0
	}
	networkType := func(c1, c2 *SurveyKey) int {
		if c1.NetworkType < c2.NetworkType {
			return -1
		} else if c1.NetworkType > c2.NetworkType {
			return 1
		}
		return 0
	}
	if freq {
		return []LessFunc{networkType, band, netname, cellid}
	}
	return []LessFunc{networkType, netname, band, cellid}
}

func SurveySampleRemove(data SurveyMap, number int) SurveyMap {
	for key, item := range data {
		if len(item) < number+1 {
			delete(data, key)
		}
	}
	return data
}

func StatRemove(data statistics.SurveyStatsMap, level float64) statistics.SurveyStatsMap {
	for key, item := range data {
		if item["DBM"].Mean <= level {
			delete(data, key)
		}
	}
	return data
}

func SelectStats(data statistics.SurveyStatsMap, filter SurveyKey) statistics.SurveyStatsMap {
	// Create a new empty map
	newData := make(statistics.SurveyStatsMap)

	for k, v := range data {
		if KeyFilter(k, filter) {
			newData[k] = v
		}
	}

	return newData
}

func SelectDeltaStats(data SurveyDeltaMap, filter SurveyKey) SurveyDeltaMap {
	// Create a new empty map
	newData := make(SurveyDeltaMap)

	for k, v := range data {
		if KeyFilter(k, filter) {
			newData[k] = v
		}
	}

	return newData
}

func KeyFilter(item, filter SurveyKey) bool {

	// Check NetName
	if filter.NetName != "" && filter.NetName != item.NetName {
		return false
	}

	// Check Band
	if filter.Band != 0 && filter.Band != item.Band {
		return false
	}

	// Check CellID
	if filter.CellID != 0 && filter.CellID != item.CellID {
		return false
	}

	// Check NetworkType
	if filter.NetworkType != "" && filter.NetworkType != item.NetworkType {
		return false
	}

	return true
}