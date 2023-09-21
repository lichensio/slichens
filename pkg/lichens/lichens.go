package lichens

import (
	"fmt"
	"math"
)

func NewSurveyStatsSummary(surveytype string) *SurveySummary {
	return &SurveySummary{
		SurveyType: surveytype,
		Stat:       make(SurveyStatsMap),
		Min:        math.MaxFloat64,
		Max:        -math.MaxFloat64,
	}
}

func (fm *SurveySummary) Set(key SurveyKey, value SurveyStats) {
	fm.Stat[key] = value
	if value["DBM"].Mean < fm.Min {
		fm.Min = value["DBM"].Mean
	}
	if value["DBM"].Mean > fm.Max {
		fm.Max = value["DBM"].Mean
	}
}

func NewSurveyDeltaSummary(surveytype string, deltatype DeltaType) *SurveyDeltaStatsSummary {
	return &SurveyDeltaStatsSummary{
		SurveyType: surveytype,
		DeltaStats: make(SurveyDeltaMap),
		DeltaType:  deltatype,
		Min:        math.MaxFloat64,
		Max:        -math.MaxFloat64,
	}
}

func (fm *SurveyDeltaStatsSummary) Set(key SurveyKey, value SurveyDeltaStats) {
	fm.DeltaStats[key] = value
	if value["DBM"].Delta < fm.Min {
		fm.Min = value["DBM"].Delta
	}
	if value["RSSI"].Delta > fm.Max {
		fm.Max = value["DBM"].Delta
	}
}

func SurveyStatGen(data SurveyInfo) SurveySummary {

	result := NewSurveyStatsSummary(data.SurveyType)

	for key, slice := range data.Surveys {

		var stats map[string]Stats

		switch key.NetworkType {
		case "2G":
			calculator := &TwoGCalculator{}
			stats = calculator.Calculate(slice)
		case "3G":
			calculator := &ThreeGCalculator{}
			stats = calculator.Calculate(slice)
		case "4G":
			calculator := &FourGCalculator{}
			stats = calculator.Calculate(slice)
		default:
			// Handle default or unknown case, maybe log an error or return.
		}
		result.Set(key, stats)
	}
	return *result
}

func GetKeys[K SurveyKey, V any](m map[K]V) ([]K, error) {
	// Check if the map is nil
	if m == nil {
		return nil, fmt.Errorf("input map is nil")
	}

	// Check if the map is empty
	if len(m) == 0 {
		return nil, fmt.Errorf("input map is empty")
	}

	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys, nil
}
