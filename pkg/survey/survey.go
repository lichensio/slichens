/*
 * Copyright © 2023 LICHENS http://www.lichens.io
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package survey

import (
	_ "encoding/csv"
	"fmt"
	"github.com/lichensio/slichens/pkg/slichens"
	"time"
)

func Survey(filename string, allStat, freq, sample bool) {
	if filename == "" {
		fmt.Println("Please provide a siretta survey file name, L____.CSV")
	} else {

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

		// flags

		survey, _ := slichens.ReadMultiCSV(filename)

		surveys := survey.Surveys

		if sample {
			surveys = slichens.SurveySampleRemove(surveys, 2)
		}

		keys := slichens.KeysSurvey(surveys)

		if freq {
			slichens.OrderedBy(band, netname, cellid).Sort(keys)
		} else {
			slichens.OrderedBy(netname, band, cellid).Sort(keys)
		}

		summary := slichens.SurveyStatGen(survey)

		key := &slichens.SurveyKey{
			Band:    0,
			CellID:  0,
			NetName: "",
		}

		summary.Stat = slichens.Select(summary.Stat, *key)
		currentTime := time.Now()

		slichens.SurveyConsolePrint("Survey", currentTime, allStat, keys, summary)
	}
}
