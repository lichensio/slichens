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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lichensio/slichens/pkg/slichens"
	"math"
	"os"
	"time"
)

func Attenuation(fileOut, fileIn string, allStat, freq, sample bool) {
	if fileOut == "" || fileIn == "" {
		fmt.Println("Please provide an outdoor  and an indoor siretta survey file name , L____.CSV")
	} else {

		currentTime := time.Now()
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

		netnameSAOI := func(c1, c2 *slichens.SurveyAvgOutIn) bool {
			return c1.Keys.NetName < c2.Keys.NetName
		}

		cellidSAOI := func(c1, c2 *slichens.SurveyAvgOutIn) bool {
			return c1.Keys.CellID < c2.Keys.CellID
		}

		bandSAOI := func(c1, c2 *slichens.SurveyAvgOutIn) bool {
			return c1.Keys.Band < c2.Keys.Band
		}

		survey_out, _ := slichens.ReadMultiCSV(fileOut)
		survey_in, _ := slichens.ReadMultiCSV(fileIn)

		surveyso := survey_out.Surveys
		surveysi := survey_in.Surveys

		if freq {
			slichens.OrderedBy(band, netname, cellid).Sort(surveyso)
			slichens.OrderedBy(band, netname, cellid).Sort(surveysi)
		} else {
			slichens.OrderedBy(netname, band, cellid).Sort(surveyso)
			slichens.OrderedBy(netname, band, cellid).Sort(surveysi)
		}

		summary_out := slichens.SurveyAverage(surveyso)
		summary_in := slichens.SurveyAverage(surveysi)

		key := &slichens.SurveyKey{
			Band:    0,
			CellID:  0,
			NetName: "",
		}
		if sample {
			summary_out.Avg = slichens.SampleRemove(summary_out.Avg, 2)
			summary_in.Avg = slichens.SampleRemove(summary_in.Avg, 2)
		}
		summary_out.Avg = slichens.Select(summary_out.Avg, *key)
		summary_in.Avg = slichens.Select(summary_in.Avg, *key)

		// table formatting
		ts := table.NewWriter()
		ts.SetTitle("OutDoor Survey: " + summary_out.SurveyType)
		ts.SetAutoIndex(true)
		if !allStat {
			ts.SetOutputMirror(os.Stdout)
			ts.AppendHeader(table.Row{"MNO", "BAND", "CellID", "RSRP Avg", "RSRQ Avg"})
			for _, item := range summary_out.Avg {
				ts.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID,
						math.Floor(item.RSRPav*100) / 100, math.Floor(item.RSRQav*100) / 100},
				})
			}
		} else {

			ts.SetOutputMirror(os.Stdout)
			ts.AppendHeader(table.Row{"MNO", "BAND", "CellID", "#", "RSRP min", "RSRP Avg", "RSRP max", "RSRP SD", "RSRQ Avg"})
			for _, item := range summary_out.Avg {
				ts.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID, item.Number,
						math.Floor(item.RSRPmin*100) / 100, math.Floor(item.RSRPav*100) / 100, math.Floor(item.RSRPmax*100) / 100, math.Floor(item.RSRPStandardDeviation*100) / 100, math.Floor(item.RSRQav*100) / 100},
				})
			}
		}
		ts.Render()
		// ts.RenderCSV()

		ts.ResetRows()

		tsin := ts //table.NewWriter()
		tsin.SetAutoIndex(true)
		tsin.SetTitle("Indoor Survey: " + summary_in.SurveyType)
		if !allStat {
			tsin.SetOutputMirror(os.Stdout)
			// tsin.AppendHeader(table.Row{"MNO", "BAND", "CellID", "RSRP Avg", "RSRQ Avg"})
			for _, item := range summary_in.Avg {
				tsin.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID,
						math.Floor(item.RSRPav*100) / 100, math.Floor(item.RSRQav*100) / 100},
				})
			}
		} else {

			tsin.SetOutputMirror(os.Stdout)
			// tsin.AppendHeader(table.Row{"MNO", "BAND", "CellID", "#", "RSRP min", "RSRP Avg", "RSRP max", "RSRP SD", "RSRQ Avg"})
			for _, item := range summary_in.Avg {
				tsin.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID, item.Number,
						math.Floor(item.RSRPmin*100) / 100, math.Floor(item.RSRPav*100) / 100, math.Floor(item.RSRPmax*100) / 100, math.Floor(item.RSRPStandardDeviation*100) / 100, math.Floor(item.RSRQav*100) / 100},
				})
			}
		}
		tsin.Render()
		// tsin.RenderCSV()
		merge, rejectiono, rejectioni := slichens.SurveyMergeOutIn2(summary_out.Avg, summary_in.Avg)

		if freq {
			slichens.OrderedBySAOI(bandSAOI, netnameSAOI, cellidSAOI).SortSAOI(merge.Data)
			slichens.OrderedBySAOI(bandSAOI, netnameSAOI, cellidSAOI).SortSAOI(rejectioni.Data)
			slichens.OrderedBySAOI(bandSAOI, netnameSAOI, cellidSAOI).SortSAOI(rejectiono.Data)
		} else {
			slichens.OrderedBySAOI(netnameSAOI, bandSAOI, cellidSAOI).SortSAOI(merge.Data)
			slichens.OrderedBySAOI(netnameSAOI, bandSAOI, cellidSAOI).SortSAOI(rejectioni.Data)
			slichens.OrderedBySAOI(bandSAOI, netnameSAOI, cellidSAOI).SortSAOI(rejectiono.Data)
		}
		tsm := table.NewWriter()
		tsm.ResetRows()
		tsm.SetAutoIndex(true)
		tsm.SetTitle("Outdoor/Indoor Survey: " + summary_in.SurveyType)
		if !allStat {
			tsm.SetOutputMirror(os.Stdout)
			tsm.AppendHeader(table.Row{"MNO", "BAND", "CellID", "Outdoor RSRP Avg", "Indoor RSRP Avg", "Delta RSRP"})
			for _, item := range merge.Data {
				tsm.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID,
						math.Floor(item.RSRPavOut*100) / 100, math.Floor(item.RSRPavIn*100) / 100, math.Floor(item.DeltaRSRP*100) / 100},
				})
			}
		} else {

			tsm.SetOutputMirror(os.Stdout)
			tsm.AppendHeader(table.Row{"MNO", "BAND", "CellID", "# Min", "Delta Outdoor/Indoor", "Outdoor RSRP min", "Outdoor RSRP Avg", "Outdoor RSRP max", "Outdoor RSRP SD", "Indoor RSRP min", "Indoor RSRP Avg", "Indoor RSRP max", "Indoor RSRP SD"})
			for _, item := range merge.Data {
				tsm.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID, item.Number, math.Floor(item.DeltaRSRP*100) / 100,
						math.Floor(item.RSRPminOut*100) / 100, math.Floor(item.RSRPavOut*100) / 100, math.Floor(item.RSRPmaxOut*100) / 100, math.Floor(item.RSRPStandardDeviationOut*100) / 100, math.Floor(item.RSRPminIn*100) / 100, math.Floor(item.RSRPavIn*100) / 100, math.Floor(item.RSRPmaxIn*100) / 100, math.Floor(item.RSRPStandardDeviationIn*100) / 100},
				})
			}
		}
		tsm.Render()
		// tsm.RenderCSV()

		// rejection 1
		tej := table.NewWriter()
		tej.ResetRows()
		tej.SetAutoIndex(true)
		tej.SetAllowedRowLength(100)
		tej.SetTitle("Outdoor/Indoor Survey : " + summary_in.SurveyType + "\n Rejection 1")
		if !allStat {
			tej.SetOutputMirror(os.Stdout)
			tej.AppendHeader(table.Row{"MNO", "BAND", "CellID"})
			for _, item := range rejectiono.Data {
				tej.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID},
				})
			}
		} else {

			tej.SetOutputMirror(os.Stdout)
			tej.AppendHeader(table.Row{"MNO", "BAND", "CellID", "# Min", "Indoor RSRP Avg"})
			for _, item := range rejectiono.Data {
				tej.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID, item.Number, math.Floor(item.RSRPavOut*100) / 100},
				})
			}
		}
		tej.Render()
		r, rerr := os.Create("RJ" + currentTime.Format("010220061504") + ".csv")
		slichens.Check(rerr)
		defer r.Close()
		tej.SetOutputMirror(r)
		tej.RenderCSV()
		// tej.RenderCSV()
		// rejection 2
		teji := tej //table.NewWriter()
		teji.ResetRows()
		teji.SetAutoIndex(true)
		teji.SetAllowedRowLength(100)
		teji.SetTitle("Outdoor/Indoor Survey : " + summary_in.SurveyType + "\n Rejection 2")
		if !allStat {
			teji.SetOutputMirror(os.Stdout)
			// teji.AppendHeader(table.Row{"MNO", "BAND", "CellID"})
			for _, item := range rejectioni.Data {
				teji.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID},
				})
			}
		} else {

			teji.SetOutputMirror(os.Stdout)
			// teji.AppendHeader(table.Row{"MNO", "BAND", "CellID"})
			for _, item := range rejectioni.Data {
				teji.AppendRows([]table.Row{
					{item.Keys.NetName, item.Keys.Band, item.Keys.CellID, item.Number, math.Floor(item.RSRPavIn*100) / 100},
				})
			}
		}
		teji.Render()
		// teji.RenderCSV()

		f, err := os.Create("AT" + currentTime.Format("010220061504") + ".csv")
		slichens.Check(err)
		defer f.Close()
		tsm.SetOutputMirror(f)
		tsm.RenderCSV()

	}
}
