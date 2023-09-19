/*
 * Copyright © 2023 LICHENS http://www.lichens.io
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package cmd

import (
	"fmt"
	"github.com/lichensio/slichens/pkg/attenuation"
	"github.com/spf13/cobra"
)

// attenuationCmd represents the attenuation command
var attenuationCmd = &cobra.Command{
	Use:   "attenuation",
	Short: "Attenuation compare outdoor to indoor signal level",
	Long:  `Attenuation compare outdoor to indoor signal level.`,
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("outfile")
		in, _ := cmd.Flags().GetString("infile")
		primarySortColumn, _ := cmd.Flags().GetString("primarySortColumn")
		if out != "" && in != "" {
			attenuation.ProcessAttenuation(out, in, primarySortColumn)
		} else {
			fmt.Println("survey files name requiered")
		}
	},
}

func init() {
	rootCmd.AddCommand(attenuationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// attenuationCmd.PersistentFlags().String("foo", "", "A help for foo")
	attenuationCmd.PersistentFlags().String("outfile", "", "Outdoor siretta filename Lxxxxx.csv")
	attenuationCmd.PersistentFlags().String("infile", "", "Indoor siretta filename Lxxxxx.csv")
	attenuationCmd.PersistentFlags().String("primarySortColumn", "", "primary Sort Column: BAND, MNO. Default POWER")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// attenuationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
