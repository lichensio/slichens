package gain

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/slichens"
	"github.com/lichensio/slichens/pkg/survey"
	"time"
)

func Gain(fileOut, fileIn string, allStat, freq, sample, level bool) {
	if fileOut == "" || fileIn == "" {
		fmt.Println("Please provide both an outdoor and an indoor with films Siretta survey file name in the format, L____.CSV.")
		return
	}

	keysOut, summaryOut, errOut := survey.ProcessSurvey(fileOut, allStat, freq, sample)
	keysIn, summaryIn, errIn := survey.ProcessSurvey(fileIn, allStat, freq, sample)

	if errOut != nil || errIn != nil {
		// Depending on the nature of errors, you can print specific error messages here.
		return
	}

	if level {
		threshold := slichens.MinimumSignalLevel // Assuming a predefined MinimumSignalLevel constant.
		summaryOut.Stat = slichens.StatRemove(summaryOut.Stat, threshold)
		summaryIn.Stat = slichens.StatRemove(summaryIn.Stat, threshold)
	}

	currentTime := time.Now()
	slichens.SurveyConsolePrint("Indoor", currentTime, allStat, freq, keysOut, summaryOut)
	slichens.SurveyConsolePrint("Indoor with booster", currentTime, allStat, freq, keysIn, summaryIn)

	// Merge and process additional data.
	merge, rejectiono, rejectioni := slichens.SurveyTwoSamplesMerge(summaryOut, summaryIn)
	keysmerge := slichens.GetKeys(merge.Data)
	keysrejectiono := slichens.GetKeys(rejectiono.Data)
	keysrejectioni := slichens.GetKeys(rejectioni.Data)

	sortFunctions := slichens.GetSortFunctions(freq)
	orderedBy := slichens.OrderedBy(sortFunctions...)
	orderedBy.Sort(keysmerge)
	orderedBy.Sort(keysrejectiono)
	orderedBy.Sort(keysrejectioni)

	slichens.TwoSampleConsoleIntersectPrint(false, "Gain Survey: ", currentTime, allStat, freq, keysmerge, merge)
	slichens.TwoSampleConsoleExcluPrint("A", "Gain Survey \n Cell ID lost \n after applying the booster", currentTime, allStat, freq, keysrejectiono, rejectiono)
	slichens.TwoSampleConsoleExcluPrint("B", "Gain Survey \n New Cell ID \n after applying the booster", currentTime, allStat, freq, keysrejectioni, rejectioni)
}
