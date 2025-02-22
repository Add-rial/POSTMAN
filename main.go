package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)	
var ele int = 7 
var generalAverages = make([]float32, ele)
var branchAverages = map[string][]float32{}

/*
	While calculating the branch wise averages, i am calculating the averages of all tests instead of
	only for total because by the time i read the statement properly i had already written the code 
*/

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
	
	sheet := f.GetSheetName(0)				//removing unwanted rows from the data
	rows, err := f.GetRows(sheet)
	if err != nil{
		log.Fatal(err)
	}
	rows = rows[1:]
	elementsToPop := findEmptyRows(rows)
	rows = removeEmptyRows(rows, elementsToPop)
	
	calaculateSum(rows)
	numberOfRows := float32(len(rows))
	calculateAverages(numberOfRows)

	printMap()
	fmt.Println(generalAverages)
}

func toFloat(str string) float32{       // converts a string to a float
	i, _ := strconv.ParseFloat(str, 32)
	return float32(i)
}

func findEmptyRows(rows [][]string) []int{
	var elementsToPop []int

	for index, row := range rows{
		if len(row) < 11{
			log.Printf("Data not found for sr no: %v\n", row[0])
			elementsToPop = append(elementsToPop, index)
			continue
		}
		total_pre_compre := toFloat(row[4]) + toFloat(row[5]) + toFloat(row[6]) + toFloat(row[7])
		total := total_pre_compre + toFloat(row[9])
		if toFloat(row[10]) != total && toFloat(row[10]) != total_pre_compre{
			log.Printf("Data mismatch for sr no: %v\n", row[0])
			elementsToPop = append(elementsToPop, index)
		}
	}

	log.Printf("The following indexex will be removed to continue%v\n", elementsToPop)

	return elementsToPop
}

func removeEmptyRows(rows [][]string, elementsToPop []int) [][]string{
	for i := len(elementsToPop) - 1; i >= 0; i--{
		rows = append(rows[:elementsToPop[i]], rows[elementsToPop[i] + 1:]...)    //...unpacks the results of the slice in the 2nd argument
	}

	return rows
}

func calaculateSum(rows [][]string) {         //calaculates only the sum of all the required elements
	for _, row := range rows{
		data := row[4:]
		for i := range ele{
			generalAverages[i] += toFloat(data[i])

			if toFloat(row[3][:4]) == 2024.0{
				branchCode := row[3][4:6]
				if branchAverages[branchCode] == nil{
					branchAverages[branchCode] = make([]float32, ele + 1)    //The last element of the slice stores the number of elements in it
					branchAverages[branchCode][ele] = 0
				}
				branchAverages[branchCode][i] += toFloat(data[i])
				branchAverages[branchCode][ele]++
			}
		}
	}
}

func calculateAverages(n float32){
	for i := range ele{
		generalAverages[i] /= n
	}
	for _, value := range branchAverages{
		for i := range ele{
			value[i] /= value[ele] / float32(ele)          //To take care of extra times the last element was 
		} 												   //incremented in calculateSum, we divide by ele
	}
}

func printMap() {
    for key, value := range branchAverages {
        fmt.Printf("%s: %v\n", key, value)
    }
}