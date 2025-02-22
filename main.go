package main

import (
	"fmt"
	"log"
	"strconv"

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
	fmt.Println(rows)
}

func toFloat(str string) float32{       // converts a string to a float
	i, _ := strconv.ParseFloat(str, 32)
	return float32(i)
}

func findEmptyRows(rows [][]string) []int{
	var elementsToPop []int

	for index, row := range rows{
		total_pre_compre := toFloat(row[4]) + toFloat(row[5]) + toFloat(row[6]) + toFloat(row[7])
		total := total_pre_compre + toFloat(row[9])
		if len(row) < 11{
			log.Printf("Data not found for sr no: %v\n", row[0])
			elementsToPop = append(elementsToPop, index)
		}
		if toFloat(row[10]) != total && toFloat(row[10]) != total_pre_compre{
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