/*
 * Copyright © 2023 LICHENS http://www.lichens.io
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package slichens

import (
	"encoding/csv"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lichensio/slichens/pkg/student"
	"golang.org/x/exp/constraints"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Front structure and information about the survey
// and back data of the survey in Survey

type SurveyResult struct {
	SurveyType         string
	FileCreated        time.Time
	IMEINumber         string
	HardwareVersion    string
	ApplicationVersion string
	FirmwareVersion    string
	Filename           string
	Timestamp          int
	Surveys            SurveyMap
}

type SurveyKey struct {
	Band    int    // 0 all
	CellID  int    // 0 all
	NetName string // all all
}

type KeySlice []SurveyKey

type SurveyData struct {
	Survey     int
	Timestamp  time.Time
	Network    string
	Index      int
	XRFCN      int
	DBM        float64
	Percentage float64
	RSSI       float64
	MCC        int
	MNC        int
	CellID     int
	LACTAC     int
	BandNum    int
	Band       int
	BSIC       string
	SCR        string
	ECIO       string
	RSCP       float64
	PCI        int
	RSRP       float64
	RSRQ       float64
	BW         int
	DL         float64
	UL         float64
	NetName    string
	Signal     string
}

type SurveyDataSlice []SurveyData

type SurveyMap map[SurveyKey]SurveyDataSlice

type SurveyStat struct {
	Number                uint
	RSRPMean              float64
	RSRPMedian            float64
	RSRPMode              float64
	RSRPRange             float64
	RSRPQuartiles         float64
	RSRPMin               float64
	RSRPMax               float64
	RSRPVariance          float64
	RSRPSkewness          float64
	RSRPKurtosis          float64
	RSRPStandardDeviation float64
}

type SurveyStatMap map[SurveyKey]SurveyStat

type SurveySummary struct {
	SurveyType string
	Stat       SurveyStatMap
}

type SurveyTwoSamples struct {
	Number                   uint
	RSRPavOut                float64
	RSRQavOut                float64
	RSRPmaxOut               float64
	RSRQmaxOut               float64
	RSRPminOut               float64
	RSRQminOut               float64
	RSRPStandardDeviationOut float64
	RSRPavIn                 float64
	RSRQavIn                 float64
	RSRPmaxIn                float64
	RSRQmaxIn                float64
	RSRPminIn                float64
	RSRQminIn                float64
	RSRPStandardDeviationIn  float64
	DeltaRSRP                float64
	DeltaRSRQ                float64
	T                        float64
	P                        float64
	Df                       float64
}

type SurveyTwoSamplesMap map[SurveyKey]SurveyTwoSamples

type SurveyTwoSamplesSummary struct {
	SurveyType string
	Data       SurveyTwoSamplesMap
}

// reading from csv siretta detailled file: LXXXXX.CSV

func ReadMultiCSV(filename string) (SurveyResult, error) {
	var survey SurveyResult
	survey.Surveys = make(map[SurveyKey]SurveyDataSlice)
	input, err := os.Open(filename)
	if err != nil {
		return survey, err
	}
	defer input.Close()

	reader := csv.NewReader(input)

	for i := 0; i < 14; i++ {
		r, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return survey, err
		}

		switch i {

		case 6:
			survey.SurveyType = r[1]
			// fmt.Println(r[1])
		case 7:
			survey.FileCreated, _ = time.Parse("01/02/06 15:04:05", r[1])
		case 8:
			survey.IMEINumber = r[1]
		case 9:
			survey.HardwareVersion = r[1]
		case 10:
			survey.ApplicationVersion = r[1]
		case 11:
			survey.FirmwareVersion = r[1]
		case 12:
			survey.Filename = r[1]
		}

		reader.FieldsPerRecord = 0
	}
	// fmt.Println(survey)
	var surveyData SurveyData
	var key SurveyKey
	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return survey, err
		}

		switch r[0] {

		case "Survey:":
			// fmt.Println(r)

		default:
			surveyData.Survey, _ = strconv.Atoi(r[0])
			surveyData.Timestamp, _ = time.Parse("01/02/06 15:04:05", r[1])
			surveyData.Network = r[2]
			surveyData.Index, _ = strconv.Atoi(r[3])
			surveyData.XRFCN, _ = strconv.Atoi(r[4])
			surveyData.DBM, _ = strconv.ParseFloat(r[5], 64)
			surveyData.Percentage, _ = strconv.ParseFloat(r[6], 64)
			surveyData.RSSI, _ = strconv.ParseFloat(r[7], 64)
			surveyData.MCC, _ = strconv.Atoi(r[8])
			surveyData.MNC, _ = strconv.Atoi(r[9])
			key.CellID, _ = strconv.Atoi(r[10])
			surveyData.LACTAC, _ = strconv.Atoi(r[11])
			surveyData.BandNum, _ = strconv.Atoi(r[12])
			frequencyParts := strings.Split(r[13], " ")
			key.Band, _ = strconv.Atoi(frequencyParts[0])
			surveyData.BSIC = r[14]
			surveyData.SCR = r[15]
			surveyData.ECIO = r[16]
			surveyData.RSCP, _ = strconv.ParseFloat(r[17], 64)
			surveyData.PCI, _ = strconv.Atoi(r[18])
			surveyData.RSRP, _ = strconv.ParseFloat(r[19], 64)
			surveyData.RSRQ, _ = strconv.ParseFloat(r[20], 64)
			surveyData.BW, _ = strconv.Atoi(r[21])
			surveyData.DL, _ = strconv.ParseFloat(r[22], 64)
			surveyData.UL, _ = strconv.ParseFloat(r[23], 64)
			key.NetName = r[24]
			survey.Surveys[key] = append(survey.Surveys[key], surveyData)
		}

		reader.FieldsPerRecord = 0
	}
	// fmt.Println(survey)
	return survey, nil
}

// Sorting algo

type lessFunc func(p1, p2 *SurveyKey) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	changes []SurveyKey
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(changes []SurveyKey) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

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

func SurveyTwoSamplesMerge(out, in SurveySummary) (SurveyTwoSamplesSummary, SurveyTwoSamplesSummary, SurveyTwoSamplesSummary) {
	var res, rejo, reji SurveyTwoSamplesSummary
	res.Data = make(map[SurveyKey]SurveyTwoSamples)
	rejo.Data = make(map[SurveyKey]SurveyTwoSamples)
	reji.Data = make(map[SurveyKey]SurveyTwoSamples)
	res.SurveyType = out.SurveyType
	var avgOutIn SurveyTwoSamples
	// res.SurveyType =
	var tej bool
	for keyo, itemo := range out.Stat {
		tej = false
		for keyi, itemi := range in.Stat {
			if keyo == keyi {

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
				avgOutIn.T = (avgOutIn.RSRPavOut*avgOutIn.RSRPavOut - avgOutIn.RSRPavIn*avgOutIn.RSRPavIn) / math.Sqrt(s2no+s2ni)
				avgOutIn.Df = math.Pow(s2no+s2ni, 2.0) / (math.Pow(s2no, 2.0)/(float64(itemo.Number)-1.0) + math.Pow(s2ni, 2.0)/(float64(itemi.Number)-1.0))
				avgOutIn.P = 2 * (1 - student.StudentCDF(avgOutIn.T, avgOutIn.Df))
				res.Data[keyo] = avgOutIn
				tej = true
				break
			}
		}
		if !tej {
			// fmt.Println(itemo.Number)
			avgOutIn.RSRPavOut = itemo.RSRPMean
			avgOutIn.Number = itemo.Number
			rejo.Data[keyo] = avgOutIn

		}
	}
	for keyi, itemi := range in.Stat {
		tej = false
		for keyo, _ := range out.Stat {
			if keyo == keyi {
				tej = true
			}
		}
		if !tej {
			// fmt.Println(itemi.Number)
			avgOutIn.Number = itemi.Number
			reji.Data[keyi] = avgOutIn
		}
	}

	return res, rejo, reji
}

func KeysSurvey(survey SurveyMap) []SurveyKey {
	res := make([]SurveyKey, 0, len(survey))
	for k := range survey {
		res = append(res, k)
	}
	return res
}

func GetKeys[K SurveyKey, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func SurveyConsolePrint(title string, currentTime time.Time, all bool, sortedKeys []SurveyKey, surveySummary SurveySummary) {
	// table formatting
	ts := table.NewWriter()
	ts.SetTitle(title + surveySummary.SurveyType)
	ts.SetAutoIndex(true)
	if !all {
		ts.SetOutputMirror(os.Stdout)
		ts.AppendHeader(table.Row{"MNO", "BAND", "CellID", "RSRP Avg"})
		for _, k := range sortedKeys {
			ts.AppendRows([]table.Row{
				{k.NetName, k.Band, k.CellID,
					math.Floor(surveySummary.Stat[k].RSRPMean*100) / 100},
			})
		}
	} else {

		ts.SetOutputMirror(os.Stdout)
		ts.AppendHeader(table.Row{"MNO", "BAND", "CellID", "#", "RSRP min", "RSRP Avg", "RSRP max", "RSRP SD"})
		for _, k := range sortedKeys {
			ts.AppendRows([]table.Row{
				{k.NetName, k.Band, k.CellID, surveySummary.Stat[k].Number,
					math.Floor(surveySummary.Stat[k].RSRPMin*100) / 100, math.Floor(surveySummary.Stat[k].RSRPMean*100) / 100, math.Floor(surveySummary.Stat[k].RSRPMax*100) / 100, math.Floor(surveySummary.Stat[k].RSRPStandardDeviation*100) / 100},
			})
		}
	}
	ts.Render()
	f, err := os.Create(title + currentTime.Format("010220061504") + ".csv")
	Check(err)
	defer f.Close()
	ts.SetOutputMirror(f)
	ts.RenderCSV()

}

func TwoSampleConsoleIntersectPrint(title string, currentTime time.Time, all bool, sortedKeys []SurveyKey, surveyStatSummary SurveyTwoSamplesSummary) {
	tsm := table.NewWriter()
	tsm.SetAutoIndex(true)
	tsm.SetTitle(title + surveyStatSummary.SurveyType)
	if !all {
		tsm.SetOutputMirror(os.Stdout)
		tsm.AppendHeader(table.Row{"MNO", "BAND", "CellID", "Outdoor RSRP Avg", "Indoor RSRP Avg", "Delta RSRP", "t", "p"})
		for _, k := range sortedKeys {
			tsm.AppendRows([]table.Row{
				{k.NetName, k.Band, k.CellID,
					math.Floor(surveyStatSummary.Data[k].RSRPavOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPavIn*100) / 100, math.Floor(surveyStatSummary.Data[k].DeltaRSRP*100) / 100, math.Floor(surveyStatSummary.Data[k].T*100) / 100, math.Floor(surveyStatSummary.Data[k].P*100) / 100},
			})
		}
	} else {

		tsm.SetOutputMirror(os.Stdout)
		tsm.AppendHeader(table.Row{"MNO", "BAND", "CellID", "#", "Delta Outdoor/Indoor", "Outdoor RSRP min", "Outdoor RSRP Avg", "Outdoor RSRP max", "Outdoor RSRP SD", "Indoor RSRP min", "Indoor RSRP Avg", "Indoor RSRP max", "Indoor RSRP SD", "t", "p"})
		for _, k := range sortedKeys {
			tsm.AppendRows([]table.Row{
				{k.NetName, k.Band, k.CellID, surveyStatSummary.Data[k].Number, math.Floor(surveyStatSummary.Data[k].DeltaRSRP*100) / 100,
					math.Floor(surveyStatSummary.Data[k].RSRPminOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPavOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPmaxOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPStandardDeviationOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPminIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPavIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPmaxIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPStandardDeviationIn*100) / 100, math.Floor(surveyStatSummary.Data[k].T*100) / 100, math.Floor(surveyStatSummary.Data[k].P*100) / 100},
			})
		}
	}
	tsm.Render()
	f, err := os.Create("INT" + currentTime.Format("010220061504") + ".csv")
	Check(err)
	defer f.Close()
	tsm.SetOutputMirror(f)
	tsm.RenderCSV()
}

func TwoSampleConsoleExcluPrint(title string, currentTime time.Time, all bool, sortedKeys []SurveyKey, surveyStatSummary SurveyTwoSamplesSummary) {
	tsm := table.NewWriter()
	tsm.SetAutoIndex(true)
	tsm.SetTitle(title + surveyStatSummary.SurveyType)
	if !all {
		tsm.SetOutputMirror(os.Stdout)
		tsm.AppendHeader(table.Row{"MNO", "BAND", "CellID"})
		for _, k := range sortedKeys {
			tsm.AppendRows([]table.Row{
				{k.NetName, k.Band, k.CellID},
			})
		}
	} else {

		tsm.SetOutputMirror(os.Stdout)
		tsm.AppendHeader(table.Row{"MNO", "BAND", "CellID", "#"})
		for _, k := range sortedKeys {
			tsm.AppendRows([]table.Row{
				{k.NetName, k.Band, k.CellID, surveyStatSummary.Data[k].Number},
			})
		}
	}
	tsm.Render()
	f, err := os.Create("RJ2" + currentTime.Format("010220061504") + ".csv")
	Check(err)
	defer f.Close()
	tsm.SetOutputMirror(f)
	tsm.RenderCSV()
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
