package gain

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/attenuation"
	"github.com/lichensio/slichens/pkg/lichens"
	"github.com/lichensio/slichens/pkg/survey"
)

func ProcessGain(filename1, filename2 string, primarySortColumn string) (lichens.SurveyDeltaStatsSummary, error) {
	if filename1 == "" || filename2 == "" {
		fmt.Println("Please provide a siretta survey file name 1 & 2, L____.CSV")
		return lichens.SurveyDeltaStatsSummary{}, fmt.Errorf("Please provide a siretta survey file name  1 & 2, L____.CSV")
	}
	summaryindoor, _ := survey.ProcessSurvey(filename1, false, false, false)
	summarybooster, _ := survey.ProcessSurvey(filename2, false, false, false)

	lichens.TableConsolePrintALL("Survey Indoor", summaryindoor, primarySortColumn)
	lichens.TableConsolePrintALL("Survey Booster", summarybooster, primarySortColumn)

	common, uniqueToSetOutdoor, uniqueToSetIndoor, _ := attenuation.GenerateDeltaStats(summaryindoor, summarybooster, lichens.IndoorBooster)
	lichens.TableConsolePrintALL("Survey unique to Indoor", uniqueToSetOutdoor, primarySortColumn)
	lichens.TableConsolePrintALL("Survey unique to Booster", uniqueToSetIndoor, primarySortColumn)
	lichens.PrintDeltaStatsTable("Gain between Indoor and Booster", false, common, "4G", primarySortColumn)
	return common, nil
}
