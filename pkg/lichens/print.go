package lichens

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"math"
	"os"
	"sort"
)

func TablePrintALL(title string, surveySummary SurveySummary, primarySortColumn string) error {

	if surveySummary.Stat == nil {
		return fmt.Errorf("invalid surveySummary provided")
	}
	keys, _ := GetKeys(surveySummary.Stat)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Manually sort the keys based on the primarySortColumn (if provided) and then by "DBM"
	sort.Slice(keys, func(i, j int) bool {
		switch primarySortColumn {
		case "MNO":
			if keys[i].NetName == keys[j].NetName {
				return surveySummary.Stat[keys[i]]["DBM"].Mean > surveySummary.Stat[keys[j]]["DBM"].Mean
			}
			return keys[i].NetName < keys[j].NetName
		case "BAND":
			if keys[i].Band == keys[j].Band {
				return surveySummary.Stat[keys[i]]["DBM"].Mean > surveySummary.Stat[keys[j]]["DBM"].Mean
			}
			return keys[i].Band < keys[j].Band
		default:
			return surveySummary.Stat[keys[i]]["DBM"].Mean > surveySummary.Stat[keys[j]]["DBM"].Mean
		}
	})

	tableWriter := table.NewWriter()
	tableWriter.SetTitle(title + " " + surveySummary.SurveyType + " ALL BAND \n " + fmt.Sprint("DBM Min: ", int(surveySummary.Min), " Max: ", int(surveySummary.Max)))
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "DBM"})

	for _, key := range keys {

		stat, ok := surveySummary.Stat[key]
		if !ok {
			return fmt.Errorf("invalid key provided")
		} else {
			dbm, ok := stat["DBM"]
			if !ok {
				return fmt.Errorf("invalid key DBM provided")
			}
			dbmValue := roundTo2DP(dbm.Mean)
			color := getColorCoding(int(dbmValue), int(surveySummary.Min), int(surveySummary.Max))

			row := table.Row{
				color.Sprint(key.NetworkType),
				color.Sprint(key.Band),
				color.Sprint(key.NetName),
				color.Sprint(key.CellID),
				color.Sprint(dbmValue),
			}
			tableWriter.AppendRow(row)
		}

	}

	tableWriter.Render()
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getColorCoding(value, min, max int) text.Colors {
	// Ensure min is less than or equal to max
	if min > max {
		min, max = max, min
	}

	threshold1 := min + (max-min)/3
	threshold2 := min + 2*(max-min)/3

	switch {
	case value <= threshold1:
		return text.Colors{text.FgRed}
	case value > threshold1 && value <= threshold2:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgGreen}
	}
}

func roundTo2DP(val float64) float64 {
	return math.Floor(val*100) / 100
}

func PrintDeltaStatsTable(title string, freq bool, surveySummary SurveyDeltaStatsSummary, networkType string, primarySortColumn string) error {
	// Check the value of surveySummary.SurveyType
	validTypes := []string{"Full", networkType}
	if !contains(validTypes, surveySummary.SurveyType) {
		return fmt.Errorf("surveySummary.SurveyType %s not valid", surveySummary.SurveyType)
	}

	tableWriter := table.NewWriter()
	tableWriter.SetTitle(title + " " + surveySummary.SurveyType + " " + " Stats " + fmt.Sprintf(" - Delta DBM Min: %d Max: %d", int(surveySummary.Min), int(surveySummary.Max)))
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	// Common logic to reduce repetition
	key4Select := SurveyKey{
		Band:        0,
		CellID:      0,
		NetName:     "",
		NetworkType: networkType,
	}

	newSurveyStatsMap, err := SelectDeltaStats(surveySummary.DeltaStats, key4Select)
	if err != nil {
		return fmt.Errorf("error selecting delta stats: %v", err)
	}

	keys, err := GetKeys(newSurveyStatsMap)
	if err != nil {
		return fmt.Errorf("error getting keys: %v", err)
	}

	sort.Slice(keys, func(i, j int) bool {
		isIndoorBooster := surveySummary.DeltaType == IndoorBooster

		rsrpDeltaI := surveySummary.DeltaStats[keys[i]]["RSRP"].Delta
		rsrpDeltaJ := surveySummary.DeltaStats[keys[j]]["RSRP"].Delta

		// Prioritize sorting by DELTA RSRP
		if rsrpDeltaI != rsrpDeltaJ {
			if isIndoorBooster {
				return rsrpDeltaI > rsrpDeltaJ
			}
			return rsrpDeltaI < rsrpDeltaJ
		}

		// If DELTA RSRP values are the same, use the existing logic
		switch primarySortColumn {
		case "MNO":
			if keys[i].NetName == keys[j].NetName {
				return rsrpDeltaI > rsrpDeltaJ
			}
			if isIndoorBooster {
				return keys[i].NetName > keys[j].NetName
			}
			return keys[i].NetName < keys[j].NetName
		case "BAND":
			if keys[i].Band == keys[j].Band {
				return rsrpDeltaI > rsrpDeltaJ
			}
			if isIndoorBooster {
				return keys[i].Band > keys[j].Band
			}
			return keys[i].Band < keys[j].Band
		default:
			return rsrpDeltaI > rsrpDeltaJ
		}
	})

	var header table.Row
	switch networkType {
	case "2G", "3G":
		header = table.Row{"GSMA", "BAND", "MNO", "CellID", "#1", "#2", "DELTA", "DIFFERENT"}
	case "4G":
		header = table.Row{"GSMA", "BAND", "MNO", "CellID", "#1", "#2", "DELTA RSRP", "DIFFERENT", "DELTA RSRQ", "DIFFERENT"}
	}

	tableWriter.AppendHeader(header)

	for _, key := range keys {
		count1 := surveySummary.DeltaStats[key]["DBM"].Number1
		count2 := surveySummary.DeltaStats[key]["DBM"].Number2
		differentRsrp := surveySummary.DeltaStats[key]["RSRP"].AreSignificantlyDiff
		Value1 := roundTo2DP(surveySummary.DeltaStats[key]["RSRP"].Delta)

		dbmValue := roundTo2DP(surveySummary.DeltaStats[key]["RSSI"].Delta)
		color := getColorCoding(int(dbmValue), int(surveySummary.Min), int(surveySummary.Max))

		var row table.Row
		switch networkType {
		case "2G", "3G":
			different := surveySummary.DeltaStats[key]["DBM"].AreSignificantlyDiff
			row = table.Row{
				color.Sprint(key.NetworkType),
				color.Sprint(key.Band),
				color.Sprint(key.NetName),
				color.Sprint(key.CellID),
				color.Sprint(count1),
				color.Sprint(count2),
				color.Sprint(Value1),
				color.Sprint(different),
			}
		case "4G":
			differentRsrq := surveySummary.DeltaStats[key]["RSRQ"].AreSignificantlyDiff
			Value2 := roundTo2DP(surveySummary.DeltaStats[key]["RSRQ"].Delta)
			row = table.Row{
				color.Sprint(key.NetworkType),
				color.Sprint(key.Band),
				color.Sprint(key.NetName),
				color.Sprint(key.CellID),
				color.Sprint(count1),
				color.Sprint(count2),
				color.Sprint(Value1),
				color.Sprint(differentRsrp),
				color.Sprint(Value2),
				color.Sprint(differentRsrq),
			}
		}
		tableWriter.AppendRow(row)
	}

	tableWriter.Render()
	return nil
}

func TablePrintStats(title string, freq bool, surveySummary SurveySummary, networkType string, primarySortColumn string) error {
	// 1. Error Handling
	validTypes := []string{"Full", networkType}
	if !contains(validTypes, surveySummary.SurveyType) {
		return fmt.Errorf("invalid surveySummary.SurveyType: %s", surveySummary.SurveyType)
	}

	tableWriter := table.NewWriter()
	tableWriter.SetTitle(title + " " + surveySummary.SurveyType + " " + networkType + " Stats")
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	key4Select := SurveyKey{
		Band:        0,
		CellID:      0,
		NetName:     "",
		NetworkType: networkType,
	}
	newSurveyStatsMap := SelectStats(surveySummary.Stat, key4Select)
	keys, _ := GetKeys(newSurveyStatsMap)

	// 2. Code Reuse
	sortKeysByColumn(keys, primarySortColumn, surveySummary)

	switch networkType {
	case "2G", "3G":
		tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "#", "DBM", "RSSI", "MIN", "MAX", "STD"})
		appendRowsToTable(tableWriter, keys, surveySummary)
	case "4G":
		tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "#", "DBM", "RSRP", "MIN", "MAX", "STD", "RSRQ", "MIN", "MAX", "STD"})
		appendRowsToTable4G(tableWriter, keys, surveySummary)
	default:
		return fmt.Errorf("unsupported networkType: %s", networkType)
	}

	tableWriter.Render()
	return nil
}

func sortKeysByColumn(keys []SurveyKey, primarySortColumn string, surveySummary SurveySummary) {
	sort.Slice(keys, func(i, j int) bool {
		switch primarySortColumn {
		case "MNO":
			if keys[i].NetName == keys[j].NetName {
				return surveySummary.Stat[keys[i]]["RSSI"].Mean > surveySummary.Stat[keys[j]]["RSSI"].Mean
			}
			return keys[i].NetName < keys[j].NetName
		case "BAND":
			if keys[i].Band == keys[j].Band {
				return surveySummary.Stat[keys[i]]["RSSI"].Mean > surveySummary.Stat[keys[j]]["RSSI"].Mean
			}
			return keys[i].Band < keys[j].Band
		default:
			return surveySummary.Stat[keys[i]]["RSSI"].Mean > surveySummary.Stat[keys[j]]["RSSI"].Mean
		}
	})
}

func appendRowsToTable(tableWriter table.Writer, keys []SurveyKey, surveySummary SurveySummary) error {
	for _, key := range keys {
		stat, ok := surveySummary.Stat[key]
		if !ok {
			return fmt.Errorf("missing data for key: %v", key)
		}

		dbm := roundTo2DP(stat["DBM"].Mean)
		Value := roundTo2DP(stat["RSSI"].Mean)
		count := stat["DBM"].Number
		Min := roundTo2DP(stat["RSSI"].Min)
		Max := roundTo2DP(stat["RSSI"].Max)
		STD := roundTo2DP(stat["RSSI"].StandardDeviation)
		color := getColorCoding(int(dbm), int(surveySummary.Min), int(surveySummary.Max))

		row := table.Row{
			color.Sprint(key.NetworkType),
			color.Sprint(key.Band),
			color.Sprint(key.NetName),
			color.Sprint(key.CellID),
			color.Sprint(count),
			color.Sprint(dbm),
			color.Sprint(Value),
			color.Sprint(Max),
			color.Sprint(Min),
			color.Sprint(STD),
		}
		tableWriter.AppendRow(row)
	}
	return nil
}

func appendRowsToTable4G(tableWriter table.Writer, keys []SurveyKey, surveySummary SurveySummary) error {
	for _, key := range keys {
		stat, ok := surveySummary.Stat[key]
		if !ok {
			return fmt.Errorf("missing data for key: %v", key)
		}

		count := stat["DBM"].Number
		// RSRP values
		ValueRSRP := roundTo2DP(stat["RSRP"].Mean)
		MinRSRP := roundTo2DP(stat["RSRP"].Min)
		MaxRSRP := roundTo2DP(stat["RSRP"].Max)
		STDRSRP := roundTo2DP(stat["RSRP"].StandardDeviation)
		// RSRQ values
		ValueRSRQ := roundTo2DP(stat["RSRQ"].Mean)
		MinRSRQ := roundTo2DP(stat["RSRQ"].Min)
		MaxRSRQ := roundTo2DP(stat["RSRQ"].Max)
		STDRSRQ := roundTo2DP(stat["RSRQ"].StandardDeviation)
		dbmValue := roundTo2DP(stat["DBM"].Mean)
		color := getColorCoding(int(dbmValue), int(surveySummary.Min), int(surveySummary.Max))

		row := table.Row{
			color.Sprint(key.NetworkType),
			color.Sprint(key.Band),
			color.Sprint(key.NetName),
			color.Sprint(key.CellID),
			color.Sprint(count),
			color.Sprint(dbmValue),
			color.Sprint(ValueRSRP),
			color.Sprint(MaxRSRP),
			color.Sprint(MinRSRP),
			color.Sprint(STDRSRP),
			color.Sprint(ValueRSRQ),
			color.Sprint(MaxRSRQ),
			color.Sprint(MinRSRQ),
			color.Sprint(STDRSRQ),
		}
		tableWriter.AppendRow(row)
	}
	return nil
}
