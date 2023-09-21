package attenuation

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/lichens"
	"github.com/lichensio/slichens/pkg/survey"
)

func GenerateDeltaStats(set1, set2 lichens.SurveySummary, DeltaType lichens.DeltaType) (lichens.SurveyDeltaStatsSummary, lichens.SurveySummary, lichens.SurveySummary, error) {
	set1StatsMap := set1.Stat
	set2StatsMap := set2.Stat
	set1KeysSet := lichens.CreateKeySet(set1StatsMap)
	set2KeysSet := lichens.CreateKeySet(set2StatsMap)
	common, uniqueToSet1, uniqueToSet2 := lichens.CompareKeySets(set1KeysSet, set2KeysSet)
	// deltas := make(map[SurveyKey]SurveyDeltaStats)
	survey := lichens.NewSurveyDeltaSummary(set1.SurveyType, DeltaType)
	// survey.DeltaType = DeltaType
	// survey.SurveyType = set1.SurveyType
	// survey := SurveyDeltaStatsSummary{set1.SurveyType, deltas, DeltaType, 0.0, 0.0}

	for key, _ := range common {
		deltaStatsForThisKey, exists := survey.DeltaStats[key]
		if !exists {
			deltaStatsForThisKey = make(lichens.SurveyDeltaStats)
		} // get the SurveyDeltaStats value (it's a copy)
		deltaStatsForThisKey.CalculateDelta(set1.Stat[key], set2.Stat[key]) // modify the copy
		survey.Set(key, deltaStatsForThisKey)
		// survey.DeltaStats[key] = deltaStatsForThisKey                       // put the modified copy back into the map
	}

	surveySet1 := lichens.NewSurveyStatsSummary(set1.SurveyType)
	for key, _ := range uniqueToSet1 {
		surveySet1.Set(key, set1StatsMap[key]) // put the modified copy back into the map
	}

	surveySet2 := lichens.NewSurveyStatsSummary(set2.SurveyType)
	for key, _ := range uniqueToSet2 {
		surveySet2.Set(key, set2StatsMap[key]) // put the modified copy back into the map
	}
	return *survey, *surveySet1, *surveySet2, nil

}

func ProcessAttenuation(filename1, filename2 string, primarySortColumn string) (lichens.SurveyDeltaStatsSummary, error) {
	if filename1 == "" || filename2 == "" {
		return lichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Please provide a siretta survey file name 1 & 2, L____.CSV")
	}

	summaryOutdoor, errOutdoor := survey.ProcessSurvey(filename1, false, false, false)
	if errOutdoor != nil {
		return lichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Error processing outdoor survey: %v", errOutdoor)
	}

	summaryIndoor, errIndoor := survey.ProcessSurvey(filename2, false, false, false)
	if errIndoor != nil {
		return lichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Error processing indoor survey: %v", errIndoor)
	}

	// Assuming the following functions return errors, handle them accordingly
	if err := lichens.TablePrintALL("Survey Outdoor", summaryOutdoor, primarySortColumn); err != nil {
		return lichens.SurveyDeltaStatsSummary{}, err
	}

	if err := lichens.TablePrintALL("Survey Indoor", summaryIndoor, primarySortColumn); err != nil {
		return lichens.SurveyDeltaStatsSummary{}, err
	}

	common, uniqueToSetOutdoor, uniqueToSetIndoor, errDelta := GenerateDeltaStats(summaryOutdoor, summaryIndoor, lichens.IndoorOutdoor)
	if errDelta != nil {
		return lichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Error generating delta stats: %v", errDelta)
	}

	if err := lichens.PrintDeltaStatsTable("Attenuation between Outdoor and Indoor", false, common, "4G", primarySortColumn); err != nil {
		return lichens.SurveyDeltaStatsSummary{}, err
	}

	if err := lichens.TablePrintALL("Survey unique to Outdoor", uniqueToSetOutdoor, primarySortColumn); err != nil {
		return lichens.SurveyDeltaStatsSummary{}, err
	}

	if err := lichens.TablePrintALL("Survey unique to Indoor", uniqueToSetIndoor, primarySortColumn); err != nil {
		return lichens.SurveyDeltaStatsSummary{}, err
	}

	return common, nil
}
