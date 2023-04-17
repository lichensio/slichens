/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/survey"
	"github.com/spf13/cobra"
)

var Verbose bool
var Freq bool
var Sample bool

// surveyCmd represents the survey command
var surveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "Summarize a siretta survey",
	Long:  `Summarize a siretta survey by making basic statistics over the samples for each cellID, MNO and frequency band observed.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		if filename != "" {
			survey.Survey(filename, Verbose, Freq, Sample)
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
	surveyCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose statistic output")
	surveyCmd.PersistentFlags().BoolVarP(&Freq, "band", "b", false, "Sorted by frequency band")
	surveyCmd.PersistentFlags().BoolVarP(&Sample, "exclude", "s", false, "Remove CellID having less than 3 samples")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// surveyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
