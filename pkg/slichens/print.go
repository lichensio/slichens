package slichens

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"math"
	"os"
	"sort"
	"strconv"
)

func getCurrentGroup(k SurveyKey, freq bool) string {
	if freq {
		return strconv.Itoa(k.Band) // Convert int to string
	}
	return k.NetName
}

func TableConsolePrintALL(title string, surveySummary SurveySummary, primarySortColumn string) {

	keys := GetKeys(surveySummary.Stat)

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
	tableWriter.SetTitle(title + " " + surveySummary.SurveyType + " ALL BAND \n " + fmt.Sprint("RSSI Min: ", int(surveySummary.Min), " Max: ", int(surveySummary.Max)))
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "DBM"})

	for _, key := range keys {
		dbmValue := roundTo2DP(surveySummary.Stat[key]["DBM"].Mean)
		rssiValue := surveySummary.Stat[key]["RSSI"].Mean
		color := getColorCoding(int(rssiValue), int(surveySummary.Min), int(surveySummary.Max))

		row := table.Row{
			color.Sprint(key.NetworkType),
			color.Sprint(key.Band),
			color.Sprint(key.NetName),
			color.Sprint(key.CellID),
			color.Sprint(dbmValue),
		}
		tableWriter.AppendRow(row)
	}

	tableWriter.Render()
}

func TableConsolePrint(title string, surveySummary SurveySummary, networkType string, primarySortColumn string) {

	// Check the value of surveySummary.SurveyType
	validTypes := []string{"Full", networkType}
	if !contains(validTypes, surveySummary.SurveyType) {
		return
	}

	key4Select := SurveyKey{
		Band:        0,
		CellID:      0,
		NetName:     "",
		NetworkType: networkType,
	}
	newSurveyStatsMap := SelectStats(surveySummary.Stat, key4Select)

	keys := GetKeys(newSurveyStatsMap)
	sort.Slice(keys, func(i, j int) bool {
		switch primarySortColumn {
		case "MNO":
			if keys[i].NetName == keys[j].NetName {
				return newSurveyStatsMap[keys[i]]["DBM"].Mean > newSurveyStatsMap[keys[j]]["DBM"].Mean
			}
			return keys[i].NetName < keys[j].NetName
		case "BAND":
			if keys[i].Band == keys[j].Band {
				return newSurveyStatsMap[keys[i]]["DBM"].Mean > newSurveyStatsMap[keys[j]]["DBM"].Mean
			}
			return keys[i].Band < keys[j].Band
		default:
			return newSurveyStatsMap[keys[i]]["DBM"].Mean > newSurveyStatsMap[keys[j]]["DBM"].Mean
		}
	})

	tableWriter := table.NewWriter()
	tableWriter.SetTitle(title + " " + surveySummary.SurveyType + " " + networkType + "\n" + fmt.Sprint("RSSI Min: ", int(surveySummary.Min), " Max: ", int(surveySummary.Max)))
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "DBM"})

	for _, key := range keys {
		dbmValue := roundTo2DP(newSurveyStatsMap[key]["DBM"].Mean)
		rssiValue := newSurveyStatsMap[key]["RSSI"].Mean
		color := getColorCoding(int(rssiValue), int(surveySummary.Min), int(surveySummary.Max))

		row := table.Row{
			color.Sprint(key.NetworkType),
			color.Sprint(key.Band),
			color.Sprint(key.NetName),
			color.Sprint(key.CellID),
			color.Sprint(dbmValue),
		}
		tableWriter.AppendRow(row)
	}
	tableWriter.Render()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getDbmColor(dbm float64) text.Colors {
	switch {
	case dbm >= -80:
		return text.Colors{text.FgGreen}
	case dbm < -80 && dbm >= -90:
		return text.Colors{text.FgYellow}
	case dbm < -90 && dbm >= -100:
		// If FgOrange is not available, you can use FgMagenta or another color as a placeholder
		return text.Colors{text.FgMagenta}
	default:
		return text.Colors{text.FgRed}
	}
}

func getSStrengthsColor(rssi int) text.Colors {
	switch {
	case rssi >= 55:
		return text.Colors{text.FgGreen}
	case rssi < 55 && rssi >= 25:
		return text.Colors{text.FgYellow}
	case rssi < 25 && rssi >= 0:
		// If FgOrange is not available, you can use FgMagenta or another color as a placeholder
		return text.Colors{text.FgRed}
	default:
		return text.Colors{text.FgRed}
	}
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

func TableConsolePrintStats(title string, freq bool, surveySummary SurveySummary, networkType string, primarySortColumn string) {

	// Check the value of surveySummary.SurveyType
	validTypes := []string{"Full", networkType}
	if !contains(validTypes, surveySummary.SurveyType) {
		return
	}

	// keysIn := GetKeys(surveySummary.Stat)

	tableWriter := table.NewWriter()
	tableWriter.SetTitle(title + " " + surveySummary.SurveyType + " " + networkType + " Stats")
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	switch networkType {
	case "2G":
		fmt.Println("2G Case")
		key4Select := SurveyKey{
			Band:        0,
			CellID:      0,
			NetName:     "",
			NetworkType: networkType,
		}
		newSurveyStatsMap := SelectStats(surveySummary.Stat, key4Select)

		keys := GetKeys(newSurveyStatsMap)

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

		tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "#", "RSSI", "MIN", "MAX", "STD"})
		for _, key := range keys {
			Value := roundTo2DP(surveySummary.Stat[key]["RSSI"].Mean)
			count := surveySummary.Stat[key]["DBM"].Number
			Min := roundTo2DP(surveySummary.Stat[key]["RSSI"].Min)
			Max := roundTo2DP(surveySummary.Stat[key]["RSSI"].Max)
			STD := roundTo2DP(surveySummary.Stat[key]["RSSI"].StandardDeviation)
			color := getColorCoding(int(Value), int(surveySummary.Min), int(surveySummary.Max))

			row := table.Row{
				color.Sprint(key.NetworkType),
				color.Sprint(key.Band),
				color.Sprint(key.NetName),
				color.Sprint(key.CellID),
				color.Sprint(count),
				color.Sprint(Value),
				color.Sprint(Max),
				color.Sprint(Min),
				color.Sprint(STD),
			}
			tableWriter.AppendRow(row)
		}
	case "3G":
		fmt.Println("3G Case")
		key4Select := SurveyKey{
			Band:        0,
			CellID:      0,
			NetName:     "",
			NetworkType: networkType,
		}
		newSurveyStatsMap := SelectStats(surveySummary.Stat, key4Select)

		keys := GetKeys(newSurveyStatsMap)
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
		tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "#", "RSSI", "MIN", "MAX", "STD"})
		for _, key := range keys {
			// ValueDBM := roundTo2DP(surveySummary.Stat[key]["DBM"].Mean)
			count := surveySummary.Stat[key]["DBM"].Number
			Value := roundTo2DP(surveySummary.Stat[key]["RSSI"].Mean)
			Min := roundTo2DP(surveySummary.Stat[key]["RSSI"].Min)
			Max := roundTo2DP(surveySummary.Stat[key]["RSSI"].Max)
			STD := roundTo2DP(surveySummary.Stat[key]["RSSI"].StandardDeviation)

			color := getColorCoding(int(Value), int(surveySummary.Min), int(surveySummary.Max))

			row := table.Row{
				color.Sprint(key.NetworkType),
				color.Sprint(key.Band),
				color.Sprint(key.NetName),
				color.Sprint(key.CellID),
				color.Sprint(count),
				color.Sprint(Value),
				color.Sprint(Max),
				color.Sprint(Min),
				color.Sprint(STD),
			}
			tableWriter.AppendRow(row)
		}
	case "4G":

		fmt.Println("4G Case")
		key4Select := SurveyKey{
			Band:        0,
			CellID:      0,
			NetName:     "",
			NetworkType: networkType,
		}
		newSurveyStatsMap := SelectStats(surveySummary.Stat, key4Select)

		keys := GetKeys(newSurveyStatsMap)
		sort.Slice(keys, func(i, j int) bool {
			switch primarySortColumn {
			case "MNO":
				if keys[i].NetName == keys[j].NetName {
					return surveySummary.Stat[keys[i]]["RSRP"].Mean > surveySummary.Stat[keys[j]]["RSRP"].Mean
				}
				return keys[i].NetName < keys[j].NetName
			case "BAND":
				if keys[i].Band == keys[j].Band {
					return surveySummary.Stat[keys[i]]["RSRP"].Mean > surveySummary.Stat[keys[j]]["RSRP"].Mean
				}
				return keys[i].Band < keys[j].Band
			default:
				return surveySummary.Stat[keys[i]]["RSRP"].Mean > surveySummary.Stat[keys[j]]["RSRP"].Mean
			}
		})

		tableWriter.AppendHeader(table.Row{"GSMA", "BAND", "MNO", "CellID", "#", "RSRP", "MIN", "MAX", "STD", "RSRQ", "MIN", "MAX", "STD"})
		for _, key := range keys {
			// RSRP
			// ValueDBM := roundTo2DP(surveySummary.Stat[key]["DBM"].Mean)
			count := surveySummary.Stat[key]["DBM"].Number
			Value := roundTo2DP(surveySummary.Stat[key]["RSRP"].Mean)
			Min := roundTo2DP(surveySummary.Stat[key]["RSRP"].Min)
			Max := roundTo2DP(surveySummary.Stat[key]["RSRP"].Max)
			STD := roundTo2DP(surveySummary.Stat[key]["RSRP"].StandardDeviation)
			// RSRQ
			Value2 := roundTo2DP(surveySummary.Stat[key]["RSRQ"].Mean)
			Min2 := roundTo2DP(surveySummary.Stat[key]["RSRQ"].Min)
			Max2 := roundTo2DP(surveySummary.Stat[key]["RSRQ"].Max)
			STD2 := roundTo2DP(surveySummary.Stat[key]["RSRQ"].StandardDeviation)
			rssiValue := roundTo2DP(surveySummary.Stat[key]["RSSI"].Mean)
			color := getColorCoding(int(rssiValue), int(surveySummary.Min), int(surveySummary.Max))

			row := table.Row{
				color.Sprint(key.NetworkType),
				color.Sprint(key.Band),
				color.Sprint(key.NetName),
				color.Sprint(key.CellID),
				color.Sprint(count),
				color.Sprint(Value),
				color.Sprint(Max),
				color.Sprint(Min),
				color.Sprint(STD),
				color.Sprint(Value2),
				color.Sprint(Max2),
				color.Sprint(Min2),
				color.Sprint(STD2),
			}
			tableWriter.AppendRow(row)
		}
	}

	tableWriter.Render()
}

func PrintDeltaStatsTable(title string, freq bool, surveySummary SurveyDeltaStatsSummary, networkType string, primarySortColumn string) {
	// Check the value of surveySummary.SurveyType
	validTypes := []string{"Full", networkType}
	if !contains(validTypes, surveySummary.SurveyType) {
		fmt.Println(surveySummary.SurveyType, " not valid")
		return
	}

	tableWriter := table.NewWriter()
	tableWriter.SetTitle(title + " " + surveySummary.SurveyType + " " + " Stats " + fmt.Sprintf(" - Delta RSSI Min: %d Max: %d", int(surveySummary.Min), int(surveySummary.Max)))
	tableWriter.SetAutoIndex(true)
	tableWriter.SetOutputMirror(os.Stdout)

	// Common logic to reduce repetition
	key4Select := SurveyKey{
		Band:        0,
		CellID:      0,
		NetName:     "",
		NetworkType: networkType,
	}

	newSurveyStatsMap := SelectDeltaStats(surveySummary.DeltaStats, key4Select)
	keys := GetKeys(newSurveyStatsMap)

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
		fmt.Println(networkType + " Case")
		header = table.Row{"GSMA", "BAND", "MNO", "CellID", "#1", "#2", "DELTA", "DIFFERENT"}
	case "4G":
		fmt.Println("4G Case")
		header = table.Row{"GSMA", "BAND", "MNO", "CellID", "#1", "#2", "DELTA RSRP", "DIFFERENT", "DELTA RSRQ", "DIFFERENT"}
	}

	tableWriter.AppendHeader(header)

	for _, key := range keys {
		count1 := surveySummary.DeltaStats[key]["DBM"].Number1
		count2 := surveySummary.DeltaStats[key]["DBM"].Number2
		differentRsrp := surveySummary.DeltaStats[key]["RSRP"].AreSignificantlyDiff
		Value1 := roundTo2DP(surveySummary.DeltaStats[key]["RSRP"].Delta)

		rssiValue := roundTo2DP(surveySummary.DeltaStats[key]["RSSI"].Delta)
		color := getColorCoding(int(rssiValue), int(surveySummary.Min), int(surveySummary.Max))

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
}
