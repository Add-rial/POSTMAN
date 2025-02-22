package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)	

func main()  {
	fmt.Println("Hello world")

	//opening the required file

	f, err := excelize.OpenFile("CSF111_202425_01_GradeBook_stripped.xlsx")
	if err != nil{
		log.Fatalf("Error encountered while opening the required file\nERROR: %v", err)
	}

	defer func() {							//anonymous function which closes the file after all required
		if err := f.Close(); err != nil{    //operations are done
			fmt.Println(err)
		}
	}()
	
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil{
		log.Fatal(err)
	}
	rows = removeEmptyRows(rows)
	
}

func removeEmptyRows(rows [][]string) [][]string{
	var elementsToPop []int

	for index, row := range rows{
		if len(row) < 11{
			elementsToPop = append(elementsToPop, index)
		}
	}

	for i := len(elementsToPop) - 1; i >= 0; i--{
		rows = append(rows[:elementsToPop[i]], rows[elementsToPop[i] + 1:]...)    //...unpacks the results of the slice in the 2nd argument
	}

	return rows
}