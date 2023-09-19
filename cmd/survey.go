package cmd

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/lichens"
	"github.com/lichensio/slichens/pkg/survey"
	"github.com/spf13/cobra"
)

// surveyCmd represents the survey command
var surveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "Summarize a siretta survey",
	Long:  `Summarize a siretta survey by making basic statistics over the samples for each cellID, MNO and frequency band observed.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		primarySortColumn, _ := cmd.Flags().GetString("primarySortColumn")
		if filename != "" {
			fmt.Println(filename)
			summaryOut, _ := survey.ProcessSurvey(filename, false, false, false)
			lichens.TableConsolePrintALL("Survey", summaryOut, primarySortColumn)
			lichens.TableConsolePrintStats("Survey", false, summaryOut, "2G", primarySortColumn)
			lichens.TableConsolePrintStats("Survey", false, summaryOut, "3G", primarySortColumn)
			lichens.TableConsolePrintStats("Survey", false, summaryOut, "4G", primarySortColumn)
		} else {
			fmt.Println("survey file name requiered")
		}
	},
}

func init() {
	rootCmd.AddCommand(surveyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	surveyCmd.PersistentFlags().String("filename", "", "siretta filename Lxxxxx.csv")
	surveyCmd.PersistentFlags().String("primarySortColumn", "", "primary Sort Column: BAND, MNO. Default POWER")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// surveyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
