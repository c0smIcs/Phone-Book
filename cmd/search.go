package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Ищет номер телефона",
	Long: `Проверьте, существует ли номер телефона в приложении Phone Book или нет`,
	Run: func(cmd *cobra.Command, args []string) {
		searchKey, _ := cmd.Flags().GetString("key")
		if searchKey == "" {
			fmt.Println("Недействительный ключ", searchKey)
			return
		}
		t := strings.ReplaceAll(searchKey, "-", "")

		if !matchTel(t) {
			fmt.Println("Недействительный номер телефона", t)
			return
		}

		temp := search(t)
		if temp == nil {
			fmt.Println("Номер не найден:", t)
			return
		}
		fmt.Println(*temp)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("key", "k", "", "Ключ для поиска")
}

func search(key string) *Entry {
	i, ok := index[key]
	if !ok {
		return nil
	}

	return &data[i]
}
