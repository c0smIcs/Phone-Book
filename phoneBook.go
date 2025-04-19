package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

type Entry2 struct {
	Name       string
	Surname    string
	Areacode   string
	Tel        string
	LastAccess string
}

type PhoneBook []Entry
type PhoneBook2 []Entry2

var CSVFILE = "/home/alan/developer/projects/phoneBook/tmp/csv.data"

var data = PhoneBook{}
var data2 = PhoneBook2{}

var index map[string]int

func readCSVFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Фaйл CSV прочитал все сразу
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	var firstLine bool = true
	var format1 = true
	for _, line := range lines {
		if firstLine {
			if len(line) == 4 {
				format1 = true
			} else if len(line) == 5 {
				format1 = false
			} else {
				return errors.New("Неизвестный формат файла!")
			}
			firstLine = false
		}

		if format1 {
			if len(line) == 4 {
				temp := Entry{
					Name:       line[0],
					Surname:    line[1],
					Tel:        line[2],
					LastAccess: line[3],
				}
				data = append(data, temp)
			}
		} else {
			if len(line) == 5 {
				temp := Entry2{
					Name:       line[0],
					Surname:    line[1],
					Areacode:   line[2],
					Tel:        line[3],
					LastAccess: line[4],
				}
				data2 = append(data2, temp)
			}
		}
	}
	return nil
}

func saveCSVFile(filepath string) error {
	csvfile, err := os.Create(filepath)
	if err != nil {
		return nil
	}
	defer csvfile.Close()

	csvwriter := csv.NewWriter(csvfile)
	for _, row := range data {
		temp := []string{row.Name, row.Surname, row.Tel, row.LastAccess}
		_ = csvwriter.Write(temp)
	}

	csvwriter.Flush()
	return nil
}

func createIndex() error {
	index = make(map[string]int)
	for i, k := range data {
		key := k.Tel
		index[key] = i
	}
	return nil
}

// Инициализируется пользователем - возвращает указатель
// Если он возвращает ноль, была ошибка
func initS(N, S, T string) *Entry {
	if T == "" || S == "" {
		return nil
	}

	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Entry{
		Name:       N,
		Surname:    S,
		Tel:        T,
		LastAccess: LastAccess}
}

func insert(pS *Entry) error {
	_, ok := index[(*pS).Tel]
	if ok {
		return fmt.Errorf("%s уже существует", pS.Tel)
	}
	data = append(data, *pS)

	_ = createIndex()

	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}

	return nil
}

func deleteEntry(key string) error {
	i, ok := index[key]
	if !ok {
		return fmt.Errorf("%s не может быть найден!", key)
	}
	data = append(data[:i], data[i+1:]...)
	// Обновление индекса - ключа больше не существует
	delete(index, key)

	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}

	return nil
}

func search(key string) *Entry {
	for i, v := range data {
		if v.Tel == key {
			return &data[i]
		}
	}
	return nil
}

func reverse() {
	sort.Sort(sort.Reverse(PhoneBook(data)))
	for _, v := range data {
		fmt.Println(v)
	}
}

func list() {
	sort.Sort(PhoneBook(data))
	for _, v := range data {
		fmt.Println(v)
	}
}

func matchTel(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`\d+$`)
	return re.Match(t)
}

func setCSVFILE() error {
	filepath := os.Getenv("PHONEBOOK")
	if filepath != "" {
		CSVFILE = filepath
	}

	_, err := os.Stat(CSVFILE)
	if err != nil {
		fmt.Println("Creating", CSVFILE)
		f, err := os.Create(CSVFILE)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	fileInfo, err := os.Stat(CSVFILE)
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		return fmt.Errorf("%s not a regular file", CSVFILE)
	}

	return nil
}

func (a PhoneBook) Len() int {
	return len(a)
}

func (a PhoneBook) Less(i, j int) bool {
	if a[i].Surname == a[j].Surname {
		return a[i].Name < a[j].Name
	}
	return a[i].Surname < a[j].Surname
}

func (a PhoneBook) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: insert | delete | search | reverse | list <arguments>")
		return
	}

	// если CSVFILE не существует, создаем пустой
	_, err := os.Stat(CSVFILE)
	if err != nil {
		fmt.Println("Creating", CSVFILE)
		f, err := os.Create(CSVFILE)
		if err != nil {
			f.Close()
			fmt.Println(err)
			return
		}
		f.Close()
	}

	fileInfo, err := os.Stat(CSVFILE)
	// Это обычный файл?
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		fmt.Println(CSVFILE, "не обычный файл!")
		return
	}

	err = readCSVFile(CSVFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = createIndex()
	if err != nil {
		fmt.Println("Невозможно создать индекс")
		return
	}

	switch arguments[1] {
	case "insert":
		if len(arguments) != 5 {
			fmt.Println("Usage: insert Name Surname Telephone")
			return
		}
		t := strings.ReplaceAll(arguments[4], "-", "")
		if !matchTel(t) {
			fmt.Println("Недействительный номер телефона:", t)
			return
		}

		temp := initS(arguments[2], arguments[3], t)
		if temp != nil {
			err := insert(temp)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	case "delete":
		if len(arguments) != 3 {
			fmt.Println("Usage: delete Number")
			return
		}

		t := strings.ReplaceAll(arguments[2], "-", "")
		if !matchTel(t) {
			fmt.Println("Недействительный номер телефона:", t)
			return
		}
		err := deleteEntry(t)
		if err != nil {
			fmt.Println(err)
		}
	case "search":
		if len(arguments) != 3 {
			fmt.Println("Usage: search Number")
			return
		}
		t := strings.ReplaceAll(arguments[2], "-", "")
		if !matchTel(t) {
			fmt.Println("Недействительный номер телефона:", t)
			return
		}
		temp := search(t)
		if temp == nil {
			fmt.Println("Number не найден:", t)
			return
		}
		fmt.Println(*temp)
	case "reverse":
		reverse()
	case "list":
		list()

	default:
		fmt.Println("Недопустимый вариант")
	}
}
