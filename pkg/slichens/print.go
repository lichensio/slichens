package slichens

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func SurveyConsolePrint(title string, currentTime time.Time, all bool, freq bool, sortedKeys []SurveyKey, surveySummary SurveySummary) {

	tableWriter := table.NewWriter()
	tableWriter.SetTitle(title + surveySummary.SurveyType + " " + surveySummary.Filename)
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	// Round a float64 value to 2 decimal places
	roundTo2DP := func(val float64) float64 {
		return math.Floor(val*100) / 100
	}

	var prevGroup string
	headersSet := false

	// Conditionally set headers and rows based on `all` and `freq` flags
	for _, key := range sortedKeys {
		currentGroup := getCurrentGroup(key, freq)

		// Insert separator if group has changed
		if currentGroup != prevGroup && prevGroup != "" {
			tableWriter.AppendSeparator()
		}
		prevGroup = currentGroup

		stats := surveySummary.Stat[key]
		if all {
			if freq {
				if !headersSet {
					tableWriter.AppendHeader(table.Row{"BAND", "MNO", "CellID", "#", "RSRP min", "RSRP Avg", "RSRP max", "RSRP SD"})
					headersSet = true
				}
				tableWriter.AppendRow(table.Row{key.Band, key.NetName, key.CellID, stats.Number, roundTo2DP(stats.RSRPMin), roundTo2DP(stats.RSRPMean), roundTo2DP(stats.RSRPMax), roundTo2DP(stats.RSRPStandardDeviation)})
			} else {
				if !headersSet {
					tableWriter.AppendHeader(table.Row{"MNO", "BAND", "CellID", "#", "RSRP min", "RSRP Avg", "RSRP max", "RSRP SD"})
					headersSet = true
				}
				tableWriter.AppendRow(table.Row{key.NetName, key.Band, key.CellID, stats.Number, roundTo2DP(stats.RSRPMin), roundTo2DP(stats.RSRPMean), roundTo2DP(stats.RSRPMax), roundTo2DP(stats.RSRPStandardDeviation)})
			}
		} else {
			if freq {
				if !headersSet {
					tableWriter.AppendHeader(table.Row{"BAND", "MNO", "CellID", "RSRP Avg"})
					headersSet = true
				}
				tableWriter.AppendRow(table.Row{key.Band, key.NetName, key.CellID, roundTo2DP(stats.RSRPMean)})
			} else {
				if !headersSet {
					tableWriter.AppendHeader(table.Row{"MNO", "BAND", "CellID", "RSRP Avg"})
					headersSet = true
				}
				tableWriter.AppendRow(table.Row{key.NetName, key.Band, key.CellID, roundTo2DP(stats.RSRPMean)})
			}
		}
	}
	tableWriter.Render()

	// Create a CSV file and render the table in CSV format
	fileName := "S" + currentTime.Format("010220061504") + ".csv"
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", fileName, err)
	}
	defer f.Close()

	tableWriter.SetOutputMirror(f)
	tableWriter.RenderCSV()
}

// Helper function to generate table header
func GetTwoSampleHeader(all bool, title1, title2 string, freq bool) table.Row {
	if !all {
		if freq {
			return table.Row{"BAND", "MNO", "CellID", title1 + " RSRP Avg", title2 + " RSRP Avg", "Delta RSRP", "t", "p"}
		}
		return table.Row{"MNO", "BAND", "CellID", title1 + " RSRP Avg", title2 + " RSRP Avg", "Delta RSRP", "t", "p"}
	}

	if freq {
		return table.Row{"BAND", "MNO", "CellID", "#1", "#2", "Delta " + title1 + "/" + title2, title1 + " RSRP min", title1 + " RSRP Avg", title1 + " RSRP max", title1 + " RSRP SD", title2 + " RSRP min", title2 + " RSRP Avg", title2 + " RSRP max", title2 + " RSRP SD", "t", "p"}
	}
	return table.Row{"MNO", "BAND", "CellID", "#1", "#2", "Delta " + title1 + "/" + title2, title1 + " RSRP min", title1 + " RSRP Avg", title1 + " RSRP max", title1 + " RSRP SD", title2 + " RSRP min", title2 + " RSRP Avg", title2 + " RSRP max", title2 + " RSRP SD", "t", "p"}
}

func TwoSampleConsoleIntersectPrint(flag bool, title string, currentTime time.Time, all bool, freq bool, sortedKeys []SurveyKey, surveyStatSummary SurveyTwoSamplesSummary) {
	tsm := table.NewWriter()
	tsm.SetAutoIndex(true)
	tsm.SetTitle(title + surveyStatSummary.SurveyType + " " + surveyStatSummary.Filename1 + " " + surveyStatSummary.Filename2)

	// Round a float64 value to 2 decimal places
	roundTo2DP := func(val float64) float64 {
		return math.Floor(val*100) / 100
	}

	var title1, title2, prevGroup string
	if flag == true {
		title1 = "Outdoor"
		title2 = "Indoor"
	} else {
		title1 = "Indoor"
		title2 = "Booster"
	}

	header := GetTwoSampleHeader(all, title1, title2, freq)
	tsm.AppendHeader(header)

	for _, k := range sortedKeys {

		currentGroup := getCurrentGroup(k, freq)

		if currentGroup != prevGroup && prevGroup != "" {
			tsm.AppendSeparator()
		}
		prevGroup = currentGroup
		if !all {
			row := []interface{}{k.NetName, k.Band, k.CellID, roundTo2DP(surveyStatSummary.Data[k].RSRPavOut), roundTo2DP(surveyStatSummary.Data[k].RSRPavIn), roundTo2DP(surveyStatSummary.Data[k].DeltaRSRP), roundTo2DP(surveyStatSummary.Data[k].T), roundTo2DP(surveyStatSummary.Data[k].P)}
			if freq {
				row = []interface{}{k.Band, k.NetName, k.CellID, roundTo2DP(surveyStatSummary.Data[k].RSRPavOut), roundTo2DP(surveyStatSummary.Data[k].RSRPavIn), roundTo2DP(surveyStatSummary.Data[k].DeltaRSRP), roundTo2DP(surveyStatSummary.Data[k].T), roundTo2DP(surveyStatSummary.Data[k].P)}
			}
			tsm.AppendRow(row)
		} else {
			row := []interface{}{k.NetName, k.Band, k.CellID, surveyStatSummary.Data[k].Number1, surveyStatSummary.Data[k].Number2, math.Floor(surveyStatSummary.Data[k].DeltaRSRP*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPminOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPavOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPmaxOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPStandardDeviationOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPminIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPavIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPmaxIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPStandardDeviationIn*100) / 100, math.Floor(surveyStatSummary.Data[k].T*100) / 100, math.Floor(surveyStatSummary.Data[k].P*100) / 100}
			if freq {
				row = []interface{}{k.Band, k.NetName, k.CellID, surveyStatSummary.Data[k].Number1, surveyStatSummary.Data[k].Number2, math.Floor(surveyStatSummary.Data[k].DeltaRSRP*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPminOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPavOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPmaxOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPStandardDeviationOut*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPminIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPavIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPmaxIn*100) / 100, math.Floor(surveyStatSummary.Data[k].RSRPStandardDeviationIn*100) / 100, math.Floor(surveyStatSummary.Data[k].T*100) / 100, math.Floor(surveyStatSummary.Data[k].P*100) / 100}
			}
			tsm.AppendRow(row)
		}
	}

	tsm.SetOutputMirror(os.Stdout)
	tsm.Render()

	f, err := os.Create("INT" + currentTime.Format("010220061504") + ".csv")
	Check(err)
	defer f.Close()
	tsm.SetOutputMirror(f)
	tsm.RenderCSV()
}

func TwoSampleConsoleExcluPrint(version, title string, currentTime time.Time, all bool, freq bool, sortedKeys []SurveyKey, surveyStatSummary SurveyTwoSamplesSummary) {
	tableWriter := table.NewWriter()
	tableWriter.SetAutoIndex(true)
	tableWriter.SetTitle(title + surveyStatSummary.SurveyType + " " + surveyStatSummary.Filename1 + " " + surveyStatSummary.Filename2)
	tableWriter.SetOutputMirror(os.Stdout)

	// Round a float64 value to 2 decimal places
	roundTo2DP := func(val float64) float64 {
		return math.Floor(val*100) / 100
	}

	var prevGroup string
	var headersSet bool
	for _, key := range sortedKeys {
		currentGroup := getCurrentGroup(key, freq)

		// Insert separator if group has changed
		if currentGroup != prevGroup && prevGroup != "" {
			tableWriter.AppendSeparator()
		}
		prevGroup = currentGroup

		if all {
			if !headersSet {
				if freq {
					tableWriter.AppendHeader(table.Row{"BAND", "MNO", "CellID", "RSRP", "#"})
				} else {
					tableWriter.AppendHeader(table.Row{"MNO", "BAND", "CellID", "RSRP", "#"})
				}
				headersSet = true
			}
			if freq {
				tableWriter.AppendRow(table.Row{key.Band, key.NetName, key.CellID, roundTo2DP(surveyStatSummary.Data[key].RSRPavOut), maxUint(surveyStatSummary.Data[key].Number1, surveyStatSummary.Data[key].Number2)})
			} else {
				tableWriter.AppendRow(table.Row{key.NetName, key.Band, key.CellID, roundTo2DP(surveyStatSummary.Data[key].RSRPavOut), maxUint(surveyStatSummary.Data[key].Number1, surveyStatSummary.Data[key].Number2)})
			}
		} else {
			if !headersSet {
				if freq {
					tableWriter.AppendHeader(table.Row{"BAND", "MNO", "CellID", "RSRP"})
				} else {
					tableWriter.AppendHeader(table.Row{"MNO", "BAND", "CellID", "RSRP"})
				}
				headersSet = true
			}
			if freq {
				tableWriter.AppendRow(table.Row{key.Band, key.NetName, key.CellID, roundTo2DP(surveyStatSummary.Data[key].RSRPavOut)})
			} else {
				tableWriter.AppendRow(table.Row{key.NetName, key.Band, key.CellID, roundTo2DP(surveyStatSummary.Data[key].RSRPavOut)})
			}
		}
	}

	// Render the table to the console
	tableWriter.Render()

	// Create a CSV file and render the table in CSV format
	fileName := "RJ2" + version + currentTime.Format("010220061504") + ".csv"
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", fileName, err)
	}
	defer f.Close()

	tableWriter.SetOutputMirror(f)
	tableWriter.RenderCSV()
}

func getCurrentGroup(k SurveyKey, freq bool) string {
	if freq {
		return strconv.Itoa(k.Band) // Convert int to string
	}
	return k.NetName
}
