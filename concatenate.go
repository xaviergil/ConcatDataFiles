// batchFile.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func randomBoolean(probability int) bool {
	if probability == 100 {
		return true
	} else {
		return rand.Intn(100) < probability
	}
}

func main() {

	fileSeparatorPtr := flag.String("separa", "", "nombre del archivo separador")
	fileDataPtr := flag.String("data", "", "fichero de datos de Exstream")
	fileFormsPtr := flag.String("forms", "", "lista de posibles formularios para modificar la selección en los datos")
	fileOutputPtr := flag.String("out", "output.txt", "fichero de salida sin separador")
	fileOutputSepPtr := flag.String("outsep", "output.txt", "fichero de salida con separador")
	numeroPtr := flag.Uint("numero", 10, "numero de clientes")
	flag.Parse()

	// check that filenames are OK
	if *fileSeparatorPtr == "" || *fileDataPtr == "" || *fileFormsPtr == "" || *fileOutputPtr == "" || *fileOutputSepPtr == ""  {
		log.Fatal("Aborting execution as some of the input files have not been defined")
	}

	// read data file
	dataLines, err := readLines(*fileDataPtr)
	if err != nil {
		log.Fatal("readLines of Data File: %s", err)
	}

	// read separator file
	separatorLines, err := readLines(*fileSeparatorPtr)
	if err != nil {
		log.Fatal("readLines of Separator File: %s", err)
	}

	// read forms file
	formLines, err := readLines(*fileFormsPtr)
	if err != nil {
		log.Fatal("readLines of Forms File: %s", err)
	}

	var linesOut []string
	var linesOutSep []string
	
	rand.Seed(time.Now().UnixNano())

	for n := uint(1); n <= *numeroPtr; n++ {
		// for s := 1; s < len(separatorLines); s++ {
		// 	linesOut = append(linesOut, separatorLines[s])
		// }
		separatorFields := separatorLines[0]
		separatorText := separatorFields + "<&#&>" + fmt.Sprintf("%05d", n)
		linesOut = append(linesOut, separatorText)

		// shuffle the formID
		formFields := strings.Split(formLines[rand.Intn(len(formLines))], ",")

		for i := 1; i < len(dataLines); i++ {
			dataFields := strings.Split(dataLines[i], "<&#&>")
			switch strings.ToUpper(dataFields[0]) {
			case "FORMID":
				linesOut = append(linesOut, dataFields[0]+"<&#&>"+formFields[0])
				linesOutSep = append(linesOutSep, dataFields[0]+"<&#&>"+formFields[0])
			case "TIPO_DOC":
				if dataFields[1] != "-" {
					linesOut = append(linesOut, dataFields[0]+"<&#&>"+formFields[1])
					linesOutSep = append(linesOutSep, dataFields[0]+"<&#&>"+formFields[1])
				}
			default:
				linesOut = append(linesOut, dataLines[i])
				linesOutSep = append(linesOutSep, dataLines[i])
			}
		}
	}
	if err := writeLines(linesOut, *fileOutputPtr); err != nil {
		log.Fatalf("writeLines on Output File: %s", err)
	}
	if err := writeLines(linesOutSep, *fileOutputSepPtr); err != nil {
		log.Fatalf("writeLines on Output Separator File: %s", err)
	}
}
