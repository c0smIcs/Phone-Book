package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var reverseCmd = &cobra.Command{
	Use:   "reverse",
	Short: "команда reverse",
	Run: func(cmd *cobra.Command, args []string) {
		reverse()
	},
}

func reverse() {
	sort.Sort(sort.Reverse(PhoneBook(data)))
	for _, v := range data {
		fmt.Println(v)
	}
}

func init() {
	rootCmd.AddCommand(reverseCmd)
}
