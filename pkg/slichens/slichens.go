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
	"github.com/axiomhq/variance"
	"golang.org/x/exp/constraints"
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
	Surveys            SurveyDataSlice
}

// Survey data structure

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
	Keys       SurveyKey
}

// Slice of survey data

type SurveyDataSlice []SurveyData

// Survey basic statistic data, average over the samples

type SurveyAvg struct {
	Keys                  SurveyKey
	Number                uint
	RSRPav                float64
	RSRQav                float64
	RSRPmax               float64
	RSRQmax               float64
	RSRPmin               float64
	RSRQmin               float64
	RSRPStandardDeviation float64
}

type SurveyAvgData []SurveyAvg

type SurveyAvgOutIn struct {
	Keys                     SurveyKey
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
}

type SurveyOutInAvgData []SurveyAvgOutIn

type SurveyOutInSummary struct {
	SurveyType string
	Data       SurveyOutInAvgData
}

type SurveySummary struct {
	SurveyType string
	Avg        SurveyAvgData
}

type SurveyKey struct {
	Band    int    // 0 all
	CellID  int    // 0 all
	NetName string // all all
}

// reading from csv siretta detailled file: LXXXXX.CSV

func ReadMultiCSV(filename string) (SurveyResult, error) {
	var survey SurveyResult
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
			surveyData.Keys.CellID, _ = strconv.Atoi(r[10])
			surveyData.LACTAC, _ = strconv.Atoi(r[11])
			surveyData.BandNum, _ = strconv.Atoi(r[12])
			frequencyParts := strings.Split(r[13], " ")
			surveyData.Keys.Band, _ = strconv.Atoi(frequencyParts[0])
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
			surveyData.Keys.NetName = r[24]
			survey.Surveys = append(survey.Surveys, surveyData)
		}

		reader.FieldsPerRecord = 0
	}
	// fmt.Println(survey)
	return survey, nil
}

// Sorting algo

type lessFunc func(p1, p2 *SurveyData) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	changes []SurveyData
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(changes []SurveyData) {
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

type lessFuncSAOI func(p1, p2 *SurveyAvgOutIn) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorterSAOI struct {
	changes []SurveyAvgOutIn
	less    []lessFuncSAOI
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorterSAOI) SortSAOI(changes []SurveyAvgOutIn) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBySAOI(less ...lessFuncSAOI) *multiSorterSAOI {
	return &multiSorterSAOI{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorterSAOI) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *multiSorterSAOI) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *multiSorterSAOI) Less(i, j int) bool {
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

func SurveyAverage(data SurveyDataSlice) SurveySummary {
	var result SurveySummary
	var avgsurvey SurveyAvg
	var avgRSRP, avgRSRQ float64
	var key SurveyKey
	var i int
	var min, max float64
	max = math.Inf(-1)
	min = math.Inf(1)
	result.SurveyType = data[0].Network
	key = data[0].Keys
	statsRSRP := variance.New()
	statsRSRQ := variance.New()
	for _, item := range data {
		if item.Keys == key {
			statsRSRP.Add(item.RSRP)
			statsRSRQ.Add(item.RSRQ)
			avgRSRP += item.RSRP
			avgRSRQ += item.RSRQ
			i += 1
			if item.RSRP > max {
				max = item.RSRP
			}
			if item.RSRP < min {
				min = item.RSRP
			}
		} else {
			avgsurvey.Keys = key
			avgsurvey.RSRPav = statsRSRP.Mean()
			avgsurvey.RSRQav = statsRSRQ.Mean()
			avgsurvey.RSRPStandardDeviation = statsRSRP.StandardDeviation()
			avgsurvey.Number = statsRSRP.NumDataValues()
			avgsurvey.RSRPmax = max
			avgsurvey.RSRPmin = min
			result.Avg = append(result.Avg, avgsurvey)
			statsRSRP.Clear()
			statsRSRQ.Clear()
			i = 1
			statsRSRP.Add(item.RSRP)
			statsRSRQ.Add(item.RSRQ)
			key = item.Keys
			min = item.RSRP
			max = item.RSRP
		}

	}
	avgsurvey.Keys = key
	avgsurvey.RSRPav = statsRSRP.Mean()
	avgsurvey.RSRQav = statsRSRQ.Mean()
	avgsurvey.RSRPStandardDeviation = statsRSRP.StandardDeviation()
	if data[len(data)-1].RSRP > max {
		max = data[len(data)-1].RSRP
	}
	if data[len(data)-1].RSRP < min {
		min = data[len(data)-1].RSRP
	}
	avgsurvey.RSRPmax = max
	avgsurvey.RSRPmin = min
	avgsurvey.Number = statsRSRP.NumDataValues()
	result.Avg = append(result.Avg, avgsurvey)
	return result
}

func SampleRemove(data SurveyAvgData, number uint) SurveyAvgData {
	var result SurveyAvgData
	for _, item := range data {
		if number < item.Number {
			result = append(result, item)
		}
	}
	return result
}

func Select(data SurveyAvgData, filter SurveyKey) SurveyAvgData {
	var result SurveyAvgData
	for _, item := range data {
		if KeyFilter(item.Keys, filter) {
			result = append(result, item)
		}
	}
	return result
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

func SurveyMergeOutIn(out, in SurveyAvgData) (SurveyOutInSummary, SurveyOutInSummary, SurveyOutInSummary) {

	var res, rejo, reji SurveyOutInSummary
	var avgOutIn SurveyAvgOutIn
	// res.SurveyType =
	var tej bool
	for _, itemo := range out {
		tej = false
		for _, itemi := range in {
			if itemo.Keys == itemi.Keys {
				avgOutIn.Keys = itemo.Keys
				avgOutIn.RSRPavOut = itemo.RSRPav
				avgOutIn.RSRPavIn = itemi.RSRPav
				avgOutIn.RSRPmaxOut = itemo.RSRPmax
				avgOutIn.RSRPmaxIn = itemi.RSRPmax
				avgOutIn.RSRPminOut = itemo.RSRPmin
				avgOutIn.RSRPminIn = itemi.RSRPmin
				avgOutIn.RSRPStandardDeviationOut = itemo.RSRPStandardDeviation
				avgOutIn.RSRPStandardDeviationIn = itemi.RSRPStandardDeviation
				avgOutIn.DeltaRSRP = -itemo.RSRPav + itemi.RSRPav
				avgOutIn.DeltaRSRQ = -itemo.RSRQav + itemi.RSRQav
				avgOutIn.Number = min(itemo.Number, itemi.Number)
				res.Data = append(res.Data, avgOutIn)
				tej = true
				break
			}
		}
		if !tej {
			avgOutIn.Keys = itemo.Keys
			rejo.Data = append(rejo.Data, avgOutIn)
		}
	}
	for _, itemi := range in {
		tej = false
		for _, itemo := range out {
			if itemo.Keys == itemi.Keys {
				tej = true
				break
			}
		}
		if !tej {
			avgOutIn.Keys = itemi.Keys
			reji.Data = append(reji.Data, avgOutIn)
		}
	}

	return res, rejo, reji
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

func SurveyMergeOutIn2(out, in SurveyAvgData) (SurveyOutInSummary, SurveyOutInSummary, SurveyOutInSummary) {

	var res, rejo, reji SurveyOutInSummary

	m := make(map[SurveyKey]uint8)
	itemo := make(map[SurveyKey]SurveyAvg)
	itemi := make(map[SurveyKey]SurveyAvg)
	for _, k := range out {
		itemo[k.Keys] = k
	}
	for _, k := range in {
		itemi[k.Keys] = k
	}
	for _, k := range out {
		m[k.Keys] |= (1 << 0)
	}
	for _, k := range in {
		m[k.Keys] |= (1 << 1)
	}
	for k, v := range m {
		var avgOutIn SurveyAvgOutIn
		a := v&(1<<0) != 0
		b := v&(1<<1) != 0
		switch {
		case a && b:
			avgOutIn.Keys = k
			avgOutIn.RSRPavOut = itemo[k].RSRPav
			avgOutIn.RSRPavIn = itemi[k].RSRPav
			avgOutIn.RSRPmaxOut = itemo[k].RSRPmax
			avgOutIn.RSRPmaxIn = itemi[k].RSRPmax
			avgOutIn.RSRPminOut = itemo[k].RSRPmin
			avgOutIn.RSRPminIn = itemi[k].RSRPmin
			avgOutIn.RSRPStandardDeviationOut = itemo[k].RSRPStandardDeviation
			avgOutIn.RSRPStandardDeviationIn = itemi[k].RSRPStandardDeviation
			avgOutIn.DeltaRSRP = -itemo[k].RSRPav + itemi[k].RSRPav
			avgOutIn.DeltaRSRQ = -itemo[k].RSRQav + itemi[k].RSRQav
			avgOutIn.Number = min(itemo[k].Number, itemi[k].Number)
			res.Data = append(res.Data, avgOutIn)
		case a && !b:
			avgOutIn.Keys = k
			avgOutIn.RSRPavOut = itemo[k].RSRPav
			avgOutIn.RSRPmaxOut = itemo[k].RSRPmax
			avgOutIn.RSRPminOut = itemo[k].RSRPmin
			avgOutIn.RSRPStandardDeviationOut = itemo[k].RSRPStandardDeviation
			avgOutIn.Number = itemo[k].Number
			rejo.Data = append(rejo.Data, avgOutIn)
		case !a && b:
			avgOutIn.Keys = k
			avgOutIn.RSRPavIn = itemi[k].RSRPav
			avgOutIn.RSRPmaxIn = itemi[k].RSRPmax
			avgOutIn.RSRPminIn = itemi[k].RSRPmin
			avgOutIn.RSRPStandardDeviationIn = itemi[k].RSRPStandardDeviation
			avgOutIn.Number = itemi[k].Number
			reji.Data = append(reji.Data, avgOutIn)
		}
	}
	return res, rejo, reji
}
