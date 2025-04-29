package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Удалить запись",
	Long:  `Удалить запись из приложения Phone Book`,
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			fmt.Println("Недействительный ключ:", key)
			return
		}

		err := deleteEntry(key)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().String("key", "", "Ключ для удаления")
}

func deleteEntry(key string) error {
	i, ok := index[key]
	if !ok {
		return fmt.Errorf("%s не может быть найден ", key)
	}
	data = append(data[:i], data[i+1:]...)

	err := saveJSONFile(JSONFILE)
	if err != nil {
		return err
	}
	return nil
}
