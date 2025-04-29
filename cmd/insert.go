package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Вставляет новые данные",
	Long:  `Эта команда вставляет новые данные в приложение Phone Book`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Println("Недействительное имя:", name)
			return
		}

		surname, _ := cmd.Flags().GetString("surname")
		if surname == "" {
			fmt.Println("Недействительная фамилия")
			return
		}

		tel, _ := cmd.Flags().GetString("telephone")
		if tel == "" {
			fmt.Println("Недействительный телефон")
			return
		}

		t := strings.ReplaceAll(tel, "-", "")
		if !matchTel(t) {
			fmt.Println("Недействительный номер телефона:", tel)
			return
		}

		temp := initS(name, surname, t)
		if temp == nil {
			fmt.Println("Недействительная запись:", temp)
			return
		}

		err := insert(temp)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringP("name", "n", "", "имя")
	insertCmd.Flags().StringP("surname", "s", "", "фамилия")
	insertCmd.Flags().StringP("telephone", "t", "", "телефон")
}

func insert(pS *Entry) error {
	_, ok := index[(*pS).Tel]
	if ok {
		return fmt.Errorf("%s уже существует", pS.Tel)
	}
	data = append(data, *pS)

	err := saveJSONFile(JSONFILE)
	if err != nil {
		return err
	}

	return nil
}

func matchTel(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`\d+$`)
	return re.Match(t)
}