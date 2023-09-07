package survey

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/slichens"
)

/*
 * Copyright © 2023 LICHENS http://www.lichens.io
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

func ProcessSurvey(filename string, allStat, freq, sample bool) (slichens.KeySlice, slichens.SurveySummary, error) {
	if filename == "" {
		fmt.Println("Please provide a siretta survey file name, L____.CSV")
		return nil, slichens.SurveySummary{}, fmt.Errorf("Please provide a siretta survey file name, L____.CSV")
	}

	// Get the survey data from the file
	survey, err := slichens.ReadMultiCSV(filename)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, slichens.SurveySummary{}, fmt.Errorf("Error reading CSV:", err)
	}

	surveys := survey.Surveys
	if sample {
		surveys = slichens.SurveySampleRemove(surveys, slichens.MinimumSampleCount) // Removed magic number, used '2' directly here as specified in your original code
	}

	keys := slichens.KeysSurvey(surveys)
	sortFunc := slichens.GetSortFunctions(freq)
	sorter := slichens.NewMultiSorter(sortFunc...)
	sorter.Sort(keys)

	summary := slichens.SurveyStatGen(survey)
	key := &slichens.SurveyKey{
		Band:    0,
		CellID:  0,
		NetName: "",
	}
	summary.Stat = slichens.Select(summary.Stat, *key)

	// currentTime := time.Now()
	// SurveyConsolePrint("Survey", currentTime, allStat, freq, keys, summary)
	return keys, summary, nil
}
