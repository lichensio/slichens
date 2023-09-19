package attenuation

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/lichens"
	"github.com/lichensio/slichens/pkg/statistics"
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
			deltaStatsForThisKey = make(statistics.SurveyDeltaStats)
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
		fmt.Println("Please provide a siretta survey file name 1 & 2, L____.CSV")
		return lichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Please provide a siretta survey file name  1 & 2, L____.CSV")
	}
	summaryOutdoor, _ := survey.ProcessSurvey(filename1, false, false, false)
	summaryIndoor, _ := survey.ProcessSurvey(filename2, false, false, false)

	lichens.TableConsolePrintALL("Survey Outdoor", summaryOutdoor, primarySortColumn)
	lichens.TableConsolePrintALL("Survey Indoor", summaryIndoor, primarySortColumn)

	common, uniqueToSetOutdoor, uniqueToSetIndoor, _ := GenerateDeltaStats(summaryOutdoor, summaryIndoor, lichens.IndoorOutdoor)
	lichens.PrintDeltaStatsTable("Attenuation between Outdoor and Indoor", false, common, "4G", primarySortColumn)
	lichens.TableConsolePrintALL("Survey unique to Outdoor", uniqueToSetOutdoor, primarySortColumn)
	lichens.TableConsolePrintALL("Survey unique to Indoor", uniqueToSetIndoor, primarySortColumn)
	return common, nil
}
