package slichens

import (
	"github.com/lichensio/slichens/pkg/student"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"math"
	"sort"
)

func SurveyStatGen(data SurveyResult) SurveySummary {
	var result SurveySummary
	result.Stat = make(map[SurveyKey]SurveyStat)
	result.SurveyType = data.SurveyType
	for key, slice := range data.Surveys {
		var data []float64
		var statsurvey SurveyStat
		for _, j := range slice {
			data = append(data, j.RSRP)
		}
		sort.Float64s(data)
		statsurvey.RSRPMax = floats.Max(data)
		statsurvey.RSRPMin = floats.Min(data)
		statsurvey.RSRPRange = statsurvey.RSRPMax - statsurvey.RSRPMin
		statsurvey.Number = uint(len(data))
		statsurvey.RSRPMean = stat.Mean(data, nil)
		statsurvey.RSRPVariance = stat.Variance(data, nil)
		statsurvey.RSRPStandardDeviation = math.Sqrt(statsurvey.RSRPVariance)
		statsurvey.RSRPMedian = stat.Quantile(0.5, stat.Empirical, data, nil)

		result.Stat[key] = statsurvey

	}
	return result
}

func Select(data SurveyStatMap, filter SurveyKey) SurveyStatMap {
	for k, _ := range data {
		if !KeyFilter(k, filter) {
			delete(data, k)
		}
	}
	return data
}

func SurveySampleRemove(data SurveyMap, number int) SurveyMap {
	for key, item := range data {
		if len(item) < number+1 {
			delete(data, key)
		}
	}
	return data
}

func StatRemove(data SurveyStatMap, level float64) SurveyStatMap {
	for key, item := range data {
		if item.RSRPMean <= level {
			delete(data, key)
		}
	}
	return data
}

func KeyFilter(item, filter SurveyKey) bool {

	if filter.NetName == "" && filter.Band == 0 && filter.CellID == 0 {
		return true
	}
	// one key
	if filter.NetName == "" && filter.Band == 0 && filter.CellID == item.CellID {
		return true
	}
	if filter.NetName == "" && filter.Band == item.Band && filter.CellID == 0 {
		return true
	}
	if filter.NetName == item.NetName && filter.Band == 0 && filter.CellID == 0 {
		return true
	}
	// two keys
	if filter.NetName == item.NetName && filter.Band == item.Band && filter.CellID == 0 {
		return true
	}
	if filter.NetName == item.NetName && filter.Band == 0 && filter.CellID == item.CellID {
		return true
	}
	if filter.NetName == "" && filter.Band == item.Band && filter.CellID == item.CellID {
		return true
	}
	// 3 keys
	if filter == item {
		return true
	}
	return false
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

func SurveyTwoSamplesMerge(out, in SurveySummary) (SurveyTwoSamplesSummary, SurveyTwoSamplesSummary, SurveyTwoSamplesSummary) {
	var mergedData, onlyOut, onlyIn SurveyTwoSamplesSummary
	mergedData.Data = make(map[SurveyKey]SurveyTwoSamples)
	onlyOut.Data = make(map[SurveyKey]SurveyTwoSamples)
	onlyIn.Data = make(map[SurveyKey]SurveyTwoSamples)
	mergedData.SurveyType = out.SurveyType

	var avgOutIn SurveyTwoSamples

	// Mapping data from out and in
	for keyo, itemo := range out.Stat {
		if itemi, ok := in.Stat[keyo]; ok { // If keyo is present in in.Stat
			avgOutIn.RSRPavOut = itemo.RSRPMean
			avgOutIn.RSRPavIn = itemi.RSRPMean
			avgOutIn.RSRPmaxOut = itemo.RSRPMax
			avgOutIn.RSRPmaxIn = itemi.RSRPMax
			avgOutIn.RSRPminOut = itemo.RSRPMin
			avgOutIn.RSRPminIn = itemi.RSRPMin
			avgOutIn.RSRPStandardDeviationOut = itemo.RSRPStandardDeviation
			avgOutIn.RSRPStandardDeviationIn = itemi.RSRPStandardDeviation
			avgOutIn.DeltaRSRP = -itemo.RSRPMean + itemi.RSRPMean
			avgOutIn.Number = min(itemo.Number, itemi.Number)

			s2no := avgOutIn.RSRPStandardDeviationOut * avgOutIn.RSRPStandardDeviationOut / float64(itemo.Number)
			s2ni := avgOutIn.RSRPStandardDeviationIn * avgOutIn.RSRPStandardDeviationIn / float64(itemi.Number)
			avgOutIn.T = (avgOutIn.RSRPavOut - avgOutIn.RSRPavIn) / math.Sqrt(s2no+s2ni)
			dfnum := math.Pow(s2no+s2ni, 2.0)
			dfDen := math.Pow(s2no, 2.0)/(float64(itemo.Number)-1.0) + math.Pow(s2ni, 2.0)/(float64(itemi.Number)-1.0)
			avgOutIn.Df = dfnum / dfDen
			avgOutIn.P = 2 * (1 - student.StudentCDF(avgOutIn.T, avgOutIn.Df))
			mergedData.Data[keyo] = avgOutIn
		} else {
			avgOutIn.RSRPavOut = itemo.RSRPMean
			avgOutIn.Number = itemo.Number
			onlyOut.Data[keyo] = avgOutIn
		}
	}

	for keyi, itemi := range in.Stat {
		if _, ok := out.Stat[keyi]; !ok { // If keyi is not present in out.Stat
			avgOutIn.RSRPavOut = itemi.RSRPMean
			avgOutIn.Number = itemi.Number
			onlyIn.Data[keyi] = avgOutIn
		}
	}

	return mergedData, onlyOut, onlyIn
}

func GetKeys[K SurveyKey, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Max(x, y float64) float64 {
	return math.Max(x, y)
}
