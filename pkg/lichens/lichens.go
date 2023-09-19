package lichens

import (
	"github.com/lichensio/slichens/pkg/statistics"
	"math"
)

func NewSurveyStatsSummary(surveytype string) *SurveySummary {
	return &SurveySummary{
		SurveyType: surveytype,
		Stat:       make(statistics.SurveyStatsMap),
		Min:        math.MaxFloat64,
		Max:        -math.MaxFloat64,
	}
}

func (fm *SurveySummary) Set(key SurveyKey, value statistics.SurveyStats) {
	fm.Stat[key] = value
	if value["RSSI"].Mean < fm.Min {
		fm.Min = value["RSSI"].Mean
	}
	if value["RSSI"].Mean > fm.Max {
		fm.Max = value["RSSI"].Mean
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

func (fm *SurveyDeltaStatsSummary) Set(key SurveyKey, value statistics.SurveyDeltaStats) {
	fm.DeltaStats[key] = value
	if value["RSSI"].Delta < fm.Min {
		fm.Min = value["RSSI"].Delta
	}
	if value["RSSI"].Delta > fm.Max {
		fm.Max = value["RSSI"].Delta
	}
}

func SurveyStatGen(data SurveyInfo) SurveySummary {

	result := NewSurveyStatsSummary(data.SurveyType)

	for key, slice := range data.Surveys {

		var stats map[string]statistics.Stats

		switch key.NetworkType {
		case "2G":
			calculator := &statistics.TwoGCalculator{}
			stats = calculator.Calculate(slice)
		case "3G":
			calculator := &statistics.ThreeGCalculator{}
			stats = calculator.Calculate(slice)
		case "4G":
			calculator := &statistics.FourGCalculator{}
			stats = calculator.Calculate(slice)
		default:
			// Handle default or unknown case, maybe log an error or return.
		}
		result.Set(key, stats)
	}
	return *result
}

func KeysSurvey(survey SurveyMap) []SurveyKey {
	res := make([]SurveyKey, 0, len(survey))
	for k := range survey {
		res = append(res, k)
	}
	return res
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CompareInts(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func GetKeys[K SurveyKey, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
