package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"slices"

	"github.com/xuri/excelize/v2"
)	

var keys = []string{"quiz", "midsem", "labtest", "weeklylab", "precompre", "compre", "total"}
var ele int = 7 
var generalAverages = map[string]float32{
	"quiz":			0,
	"midsem":		0,
	"labtest":		0,
	"weeklylab":	0,
	"precompre":	0,
	"compre":		0,
	"total":		0,
}
var branchAverages = map[string][]float32{}
var top3 = map[string][][]string{
	"quiz":			{{}},
	"midsem":		{{}},
	"labtest":		{{}},
	"weeklylab":	{{}},
	"precompre":	{{}},
	"compre":		{{}},
	"total":		{{}},
}

/*
	While calculating the branch wise averages, i am calculating the averages of all tests instead of
	only for total because by the time i read the statement properly i had already written the code 
*/

func main()  {
	fmt.Println("Hello world")

	var class int = -1

	//Reading the command line arguments
	/* if len(os.Args) < 2 {
		log.Fatalln("ENTER PATH TO THE EXCEL FILE TO BE PARSED")
	}else {

	}

	fmt.Println(os.Args[1:]) */

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
	elementsToPop := findEmptyRows(rows, class)
	rows = removeEmptyRows(rows, elementsToPop)

	calaculateSum(rows)                     //Calculating averages
	numberOfRows := float32(len(rows))
	calculateAverages(numberOfRows)

	getTop3(rows)                           //Gets the top 3 across all categories and stores it into the map
	//printResults()
	toJSON()
}

func toFloat(str string) float32{       // converts a string to a float
	i, _ := strconv.ParseFloat(str, 32)
	return float32(i)
}

func findEmptyRows(rows [][]string, classToFilter int) []int{
	var elementsToPop []int
	var isClassFilter bool = classToFilter != -1     //didn't put in the else if condition to reduce computations of the same task
	for index, row := range rows{
		total_pre_compre := toFloat(row[4]) + toFloat(row[5]) + toFloat(row[6]) + toFloat(row[7])
		total := total_pre_compre + toFloat(row[9])
		if len(row) < 11{
			log.Printf("Data not found for sr no: %v\n", row[0])
			elementsToPop = append(elementsToPop, index)
		}else if toFloat(row[10]) != total && toFloat(row[10]) != total_pre_compre{
			log.Printf("Data mismatch for sr no: %v\n", row[0])
			elementsToPop = append(elementsToPop, index)
		}else if isClassFilter && int(toFloat(row[1])) != classToFilter {
			elementsToPop = append(elementsToPop, index)        //Removees cases where class does not match classFilter
		}
	}


	return elementsToPop
}

func removeEmptyRows(rows [][]string, elementsToPop []int) [][]string{
	for i := len(elementsToPop) - 1; i >= 0; i--{
		rows = slices.Delete(rows, elementsToPop[i], elementsToPop[i] + 1)
	}

	return rows
}

func calaculateSum(rows [][]string) {         //calaculates only the sum of all the required elements

	for _, row := range rows{
		data := row[4:]
		for i, key := range keys{
			generalAverages[key] += toFloat(data[i])

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
	for i := range generalAverages{
		generalAverages[i] /= n
	}
	for _, value := range branchAverages{
		for i := range ele{
			value[i] /= value[ele] / float32(ele)          //To take care of extra times the last element was 
		} 												   //incremented in calculateSum, we divide by ele
	}
}

func getTop3(rows [][]string){
	n := 4
	for _, key := range keys{
		rowscopy := slices.Clone(rows)
        customSort(rowscopy, n)

        // Assign the top 3 from the sorted copy to the map
        top3[key] = rowscopy[:3]
		rowscopy = nil
        n++
	}
}

func customSort(rows [][]string, n int) {
    for i := 0; i < len(rows)-1; i++ {
        for j := 0; j < len(rows)-i-1; j++ {
            if toFloat(rows[j][n]) < toFloat(rows[j+1][n]) { 
                rows[j], rows[j+1] = rows[j+1], rows[j]
            }
        }
    }
}


func printResults(){
	fmt.Println("\n\n\n\n---------------------------------------------------------------------------\nGeneral Averages: ")
	for _, key := range keys{
		fmt.Printf("Avg %v: %v\n", key, generalAverages[key])
	}
	fmt.Println("\n\n---------------------------------------------------------------------------\nBranch-wise Total Averages: ")
	for i := range branchAverages{
		fmt.Printf("Branch: %v--->%v\n", i, branchAverages[i][ele - 1])
	}
	fmt.Print("\n\n---------------------------------------------------------------------------\nTop 3 Rankings: ")
	for l, key := range keys{
		fmt.Printf("\n%v", key)
		for i, j := range top3[key]{
			fmt.Printf("\n\t\tRank: %v--->Emplid: %v......Marks: %v", i + 1, j[2], j[l + 4])
		}
	}
}

func toJSON(){
	top3Map := make(map[string]map[int]map[string]string)
	for l, key := range keys{
		top3Map[key] = make(map[int]map[string]string)
		for i, j := range top3[key]{
			top3Map[key][i +1 ] = make(map[string]string)
			top3Map[key][i + 1]["emplid"] = j[2]
			top3Map[key][i + 1]["marks"] = j[l + 4]
			top3Map[key][i + 1]["rank"] = strconv.Itoa(i + 1)
		}
	}

	branchAveragesMap := make(map[string]float32)
	for key, i := range branchAverages{
		branchAveragesMap[key] = i[ele - 1]
	}

	superMap := make(map[string]any)
	superMap["General Average"] = generalAverages
	superMap["Branch averages"] = branchAveragesMap
	superMap["Top 3"] = top3Map

	j, err := json.MarshalIndent(superMap, "", "	")
	if err != nil{
		log.Printf("Unable to convert to json\nERROR: %v\n", err)
	}
	file, err := os.Create("data.json")
	if err != nil {
		log.Printf("Error creating file\nERROR: %v\n", err)
	}
	defer file.Close()

	_, err = file.Write(j)
	if err != nil {
		log.Printf("Error writing to file\nERROR: %v\n", err)
	}

	fmt.Println("Data successfully written to data.json")
}