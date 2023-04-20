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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lichensio/slichens/pkg/slichens"
	"math"
	"os"
	"time"
)

func Survey(filename string, allStat, freq, sample bool) {
	if filename == "" {
		fmt.Println("Please provide a siretta survey file name, L____.CSV")
	} else {

		// sorting prep
		netname := func(c1, c2 *slichens.SurveyData) bool {
			return c1.Keys.NetName < c2.Keys.NetName
		}

		cellid := func(c1, c2 *slichens.SurveyData) bool {
			return c1.Keys.CellID < c2.Keys.CellID
		}

		band := func(c1, c2 *slichens.SurveyData) bool {
			return c1.Keys.Band < c2.Keys.Band
		}

		// flags

		survey, _ := slichens.ReadMultiCSV(filename)

		surveys := survey.Surveys

		if freq {
			slichens.OrderedBy(band, netname, cellid).Sort(surveys)
		} else {
			slichens.OrderedBy(netname, band, cellid).Sort(surveys)
		}

		summary := slichens.SurveyAverage(surveys)

		key := &slichens.SurveyKey{
			Band:    0,
			CellID:  0,
			NetName: "",
		}
		if sample {
			summary.Avg = slichens.SampleRemove(summary.Avg, 2)
		}
		summary.Avg = slichens.Select(summary.Avg, *key)
		ts := table.NewWriter()

		if !allStat {
			ts.SetOutputMirror(os.Stdout)
			ts.AppendHeader(table.Row{"MNO", "BAND", "CellID", "RSRP Avg", "RSRQ Avg"})
			for _, item := range summary.Avg {
				ts.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID,
						math.Floor(item.RSRPav*100) / 100, math.Floor(item.RSRQav*100) / 100},
				})
			}
		} else {

			ts.SetOutputMirror(os.Stdout)
			ts.AppendHeader(table.Row{"MNO", "BAND", "CellID", "#", "RSRP min", "RSRP Avg", "RSRP max", "RSRP SD", "RSRQ Avg"})
			for _, item := range summary.Avg {
				ts.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID, item.Number,
						math.Floor(item.RSRPmin*100) / 100, math.Floor(item.RSRPav*100) / 100, math.Floor(item.RSRPmax*100) / 100, math.Floor(item.RSRPStandardDeviation*100) / 100, math.Floor(item.RSRQav*100) / 100},
				})
			}
		}
		ts.Render()
		currentTime := time.Now()

		f, err := os.Create("SR" + currentTime.Format("010220061504") + ".csv")
		slichens.Check(err)
		defer f.Close()
		ts.SetOutputMirror(f)
		ts.RenderCSV()
	}
}
