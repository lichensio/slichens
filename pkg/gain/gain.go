package gain

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/attenuation"
	"github.com/lichensio/slichens/pkg/slichens"
	"github.com/lichensio/slichens/pkg/survey"
)

func ProcessGain(filename1, filename2 string, primarySortColumn string) (slichens.SurveyDeltaStatsSummary, error) {
	if filename1 == "" || filename2 == "" {
		fmt.Println("Please provide a siretta survey file name 1 & 2, L____.CSV")
		return slichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Please provide a siretta survey file name  1 & 2, L____.CSV")
	}
	summaryindoor, _ := survey.ProcessSurvey(filename1, false, false, false)
	summarybooster, _ := survey.ProcessSurvey(filename2, false, false, false)

	slichens.TableConsolePrintALL("Survey Indoor", summaryindoor, primarySortColumn)
	slichens.TableConsolePrintALL("Survey Booster", summarybooster, primarySortColumn)

	common, uniqueToSetOutdoor, uniqueToSetIndoor, _ := attenuation.GenerateDeltaStats(summaryindoor, summarybooster, slichens.IndoorBooster)
	slichens.TableConsolePrintALL("Survey unique to Indoor", uniqueToSetOutdoor, primarySortColumn)
	slichens.TableConsolePrintALL("Survey unique to Booster", uniqueToSetIndoor, primarySortColumn)
	slichens.PrintDeltaStatsTable("Gain between Indoor and Booster", false, common, "4G", primarySortColumn)
	return common, nil
}
