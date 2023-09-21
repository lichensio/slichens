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
		filename, errorFN := cmd.Flags().GetString("filename")
		primarySortColumn, _ := cmd.Flags().GetString("primarySortColumn")

		if errorFN != nil {
			fmt.Println("Error retrieving filename:", errorFN)
			return
		}

		if filename == "" {
			fmt.Println("survey file name required")
			return
		}

		fmt.Println(filename)
		summaryOut, errorPS := survey.ProcessSurvey(filename, false, false, false)
		if errorPS != nil {
			fmt.Println("survey.ProcessSurvey error:", errorPS)
			return
		}

		err1 := lichens.TablePrintALL("Survey", summaryOut, primarySortColumn)
		err2 := lichens.TablePrintStats("Survey", false, summaryOut, "2G", primarySortColumn)
		err3 := lichens.TablePrintStats("Survey", false, summaryOut, "3G", primarySortColumn)
		err4 := lichens.TablePrintStats("Survey", false, summaryOut, "4G", primarySortColumn)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			if err1 != nil {
				fmt.Println("TablePrintALL error:", err1)
			}
			if err2 != nil {
				fmt.Println("TablePrintStats (2G) error:", err2)
			}
			if err3 != nil {
				fmt.Println("TablePrintStats (3G) error:", err3)
			}
			if err4 != nil {
				fmt.Println("TablePrintStats (4G) error:", err4)
			}
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
