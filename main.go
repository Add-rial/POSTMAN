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
		log.Fatalf("Error encountered while opening the required file\nERROR: %v\n", err)
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
	elementsToPop := findEmptyRows(rows)
	rows = removeEmptyRows(rows, elementsToPop)
}

func findEmptyRows(rows [][]string) []int{
	var elementsToPop []int

	for index, row := range rows{
		if len(row) < 11{
			log.Printf("Data not found for sr no: %v\n", row[0])
			elementsToPop = append(elementsToPop, index)
		}
		if row[10] != row[4] + row[5] + row[6] + row[7] + row[8] + row[9]{
			log.Printf("Data mismatch for sr no: %v\n", row[0])
			elementsToPop = append(elementsToPop, index)
		}
	}

	log.Printf("The following sr. no. will be removed to continue%v\n", elementsToPop)

	return elementsToPop
}

func removeEmptyRows(rows [][]string, elementsToPop []int) [][]string{
	for i := len(elementsToPop) - 1; i >= 0; i--{
		rows = append(rows[:elementsToPop[i]], rows[elementsToPop[i] + 1:]...)    //...unpacks the results of the slice in the 2nd argument
	}

	return rows
}