package Utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type SpellbookPing struct {
	Version string
}

func Error(msg string, err error) {
	fmt.Printf("%s, %v", msg, err)
	// os.Exit(1)
}

func GetBleveIndex() {

}

func InsertJsonToDB(objs []interface{}) {

}

func HandleError(e error) {

	if e != nil {
		fmt.Println(e)
	}
}
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func GenerateTableHeader() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
