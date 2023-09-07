package attenuation

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/slichens"
	"github.com/lichensio/slichens/pkg/survey"
	"time"
)

func Attenuation(fileOut, fileIn string, allStat, freq, sample, level bool) {
	if fileOut == "" || fileIn == "" {
		fmt.Println("Please provide both an outdoor and an indoor Siretta survey file name in the format, L____.CSV.")
		return
	}

	var keysOut, summaryOut, errOut = survey.ProcessSurvey(fileOut, allStat, freq, sample)
	var keysIn, summaryIn, errIn = survey.ProcessSurvey(fileIn, allStat, freq, sample)

	if errOut != nil || errIn != nil {
		return
	}

	if level {
		threshold := slichens.MinimumSignalLevel // Could be moved to a constant for clarity
		summaryOut.Stat = slichens.StatRemove(summaryOut.Stat, threshold)
		summaryIn.Stat = slichens.StatRemove(summaryIn.Stat, threshold)
	}

	currentTime := time.Now()
	slichens.SurveyConsolePrint("OutDoor", currentTime, allStat, freq, keysOut, summaryOut)
	slichens.SurveyConsolePrint("Indoor", currentTime, allStat, freq, keysIn, summaryIn)

	// Merge and process additional data as before...
	// ... (remaining logic for merged data and exclusions)
	merge, rejectiono, rejectioni := slichens.SurveyTwoSamplesMerge(summaryOut, summaryIn)
	keysmerge := slichens.GetKeys(merge.Data)
	keysrejectiono := slichens.GetKeys(rejectiono.Data)
	keysrejectioni := slichens.GetKeys(rejectioni.Data)

	sortFunctions := slichens.GetSortFunctions(freq)
	orderedBy := slichens.OrderedBy(sortFunctions...)
	orderedBy.Sort(keysmerge)
	orderedBy.Sort(keysrejectiono)
	orderedBy.Sort(keysrejectioni)

	slichens.TwoSampleConsoleIntersectPrint(true, "Outdoor/Indoor Survey: ", currentTime, allStat, freq, keysmerge, merge)
	slichens.TwoSampleConsoleExcluPrint("A", "Outdoor/Indoor Survey \n Signals outdoor \n not found indoor", currentTime, allStat, freq, keysrejectiono, rejectiono)
	slichens.TwoSampleConsoleExcluPrint("B", "Outdoor/Indoor Survey \n Signals indoor \n not found outdoor", currentTime, allStat, freq, keysrejectioni, rejectioni)
}
