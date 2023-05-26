/*
 * Copyright © 2023 LICHENS http://www.lichens.io
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package attenuation

import (
	_ "encoding/csv"
	"fmt"
	"github.com/lichensio/slichens/pkg/slichens"
	"time"
)

func Attenuation(fileOut, fileIn string, allStat, freq, sample, level bool) {
	if fileOut == "" || fileIn == "" {
		fmt.Println("Please provide an outdoor  and an indoor siretta survey file name , L____.CSV")
	} else {

		currentTime := time.Now()
		// sorting prep
		netname := func(c1, c2 *slichens.SurveyKey) bool {
			return c1.NetName < c2.NetName
		}

		cellid := func(c1, c2 *slichens.SurveyKey) bool {
			return c1.CellID < c2.CellID
		}

		band := func(c1, c2 *slichens.SurveyKey) bool {
			return c1.Band < c2.Band
		}

		surveyOut, _ := slichens.ReadMultiCSV(fileOut)
		surveyIn, _ := slichens.ReadMultiCSV(fileIn)

		surveyso := surveyOut.Surveys
		surveysi := surveyIn.Surveys

		keyso := slichens.KeysSurvey(surveyso)
		keysi := slichens.KeysSurvey(surveysi)

		if sample {
			surveyso = slichens.SurveySampleRemove(surveyso, 2)
			surveysi = slichens.SurveySampleRemove(surveysi, 2)
		}

		if freq {
			slichens.OrderedBy(band, netname, cellid).Sort(keyso)
			slichens.OrderedBy(band, netname, cellid).Sort(keysi)
		} else {
			slichens.OrderedBy(netname, band, cellid).Sort(keyso)
			slichens.OrderedBy(netname, band, cellid).Sort(keysi)
		}

		summaryOut := slichens.SurveyStatGen(surveyOut)
		summaryIn := slichens.SurveyStatGen(surveyIn)

		key := &slichens.SurveyKey{
			Band:    0,
			CellID:  0,
			NetName: "",
		}
		if level {
			summaryOut.Stat = slichens.StatRemove(summaryOut.Stat, -139.99)
			summaryIn.Stat = slichens.StatRemove(summaryIn.Stat, -139.99)
		}
		summaryOut.Stat = slichens.Select(summaryOut.Stat, *key)
		summaryIn.Stat = slichens.Select(summaryIn.Stat, *key)

		slichens.SurveyConsolePrint("OutDoor", currentTime, allStat, keyso, summaryOut)

		slichens.SurveyConsolePrint("Indoor ", currentTime, allStat, keysi, summaryIn)

		merge, rejectiono, rejectioni := slichens.SurveyTwoSamplesMerge(summaryOut, summaryIn)
		keysmerge := slichens.GetKeys(merge.Data)
		keysrejectiono := slichens.GetKeys(rejectiono.Data)
		keysrejectioni := slichens.GetKeys(rejectioni.Data)

		if freq {
			slichens.OrderedBy(band, cellid, netname).Sort(keysmerge)
			slichens.OrderedBy(band, cellid, netname).Sort(keysrejectiono)
			slichens.OrderedBy(band, cellid, netname).Sort(keysrejectioni)
		} else {
			slichens.OrderedBy(netname, band, cellid).Sort(keysmerge)
			slichens.OrderedBy(netname, band, cellid).Sort(keysrejectiono)
			slichens.OrderedBy(netname, band, cellid).Sort(keysrejectioni)
		}

		slichens.TwoSampleConsoleIntersectPrint("Outdoor/Indoor Survey: ", currentTime, allStat, keysmerge, merge)
		slichens.TwoSampleConsoleExcluPrint("Outdoor/Indoor Survey \n Exclusion 1 \n", currentTime, allStat, keysrejectiono, rejectiono)
		slichens.TwoSampleConsoleExcluPrint("Outdoor/Indoor Survey \n Exclusion 2 \n", currentTime, allStat, keysrejectioni, rejectioni)
	}
}
