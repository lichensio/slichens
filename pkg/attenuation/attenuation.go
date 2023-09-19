package attenuation

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/slichens"
	"github.com/lichensio/slichens/pkg/stats"
	"github.com/lichensio/slichens/pkg/survey"
)

func GenerateDeltaStats(set1, set2 slichens.SurveySummary, DeltaType slichens.DeltaType) (slichens.SurveyDeltaStatsSummary, slichens.SurveySummary, slichens.SurveySummary, error) {
	set1StatsMap := set1.Stat
	set2StatsMap := set2.Stat
	set1KeysSet := slichens.CreateKeySet(set1StatsMap)
	set2KeysSet := slichens.CreateKeySet(set2StatsMap)
	common, uniqueToSet1, uniqueToSet2 := slichens.CompareKeySets(set1KeysSet, set2KeysSet)
	// deltas := make(map[SurveyKey]SurveyDeltaStats)
	survey := slichens.NewSurveyDeltaSummary(set1.SurveyType, DeltaType)
	// survey.DeltaType = DeltaType
	// survey.SurveyType = set1.SurveyType
	// survey := SurveyDeltaStatsSummary{set1.SurveyType, deltas, DeltaType, 0.0, 0.0}

	for key, _ := range common {
		deltaStatsForThisKey, exists := survey.DeltaStats[key]
		if !exists {
			deltaStatsForThisKey = make(stats.SurveyDeltaStats)
		} // get the SurveyDeltaStats value (it's a copy)
		deltaStatsForThisKey.CalculateDelta(set1.Stat[key], set2.Stat[key]) // modify the copy
		survey.Set(key, deltaStatsForThisKey)
		// survey.DeltaStats[key] = deltaStatsForThisKey                       // put the modified copy back into the map
	}

	surveySet1 := slichens.NewSurveyStatsSummary(set1.SurveyType)
	for key, _ := range uniqueToSet1 {
		surveySet1.Set(key, set1StatsMap[key]) // put the modified copy back into the map
	}

	surveySet2 := slichens.NewSurveyStatsSummary(set2.SurveyType)
	for key, _ := range uniqueToSet2 {
		surveySet2.Set(key, set2StatsMap[key]) // put the modified copy back into the map
	}
	return *survey, *surveySet1, *surveySet2, nil

}

func ProcessAttenuation(filename1, filename2 string, primarySortColumn string) (slichens.SurveyDeltaStatsSummary, error) {
	if filename1 == "" || filename2 == "" {
		fmt.Println("Please provide a siretta survey file name 1 & 2, L____.CSV")
		return slichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Please provide a siretta survey file name  1 & 2, L____.CSV")
	}
	summaryOutdoor, _ := survey.ProcessSurvey(filename1, false, false, false)
	summaryIndoor, _ := survey.ProcessSurvey(filename2, false, false, false)

	slichens.TableConsolePrintALL("Survey Outdoor", summaryOutdoor, primarySortColumn)
	slichens.TableConsolePrintALL("Survey Indoor", summaryIndoor, primarySortColumn)

	common, uniqueToSetOutdoor, uniqueToSetIndoor, _ := GenerateDeltaStats(summaryOutdoor, summaryIndoor, slichens.IndoorOutdoor)
	slichens.PrintDeltaStatsTable("Attenuation between Outdoor and Indoor", false, common, "4G", primarySortColumn)
	slichens.TableConsolePrintALL("Survey unique to Outdoor", uniqueToSetOutdoor, primarySortColumn)
	slichens.TableConsolePrintALL("Survey unique to Indoor", uniqueToSetIndoor, primarySortColumn)
	return common, nil
}
