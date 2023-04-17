/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slichens/pkg/gain"

	"github.com/spf13/cobra"
)

// gainCmd represents the gain command
var gainCmd = &cobra.Command{
	Use:   "gain",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("indoor")
		in, _ := cmd.Flags().GetString("mbooster")
		if out != "" && in != "" {
			gain.Gain(out, in, Verbose, Freq, true)
		} else {
			fmt.Println("survey files name requiered")
		}
	},
}

func init() {
	rootCmd.AddCommand(gainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gainCmd.PersistentFlags().String("foo", "", "A help for foo")
	gainCmd.PersistentFlags().String("indoor", "", "Indoor siretta filename Lxxxxx.csv")
	gainCmd.PersistentFlags().String("mbooster", "", "Improved Indoor siretta filename Lxxxxx.csv")
	gainCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose statistic output")
	gainCmd.PersistentFlags().BoolVarP(&Freq, "band", "b", false, "Sorted by frequency band")
	gainCmd.PersistentFlags().BoolVarP(&Sample, "exclude", "s", false, "Sorted by frequency band")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
