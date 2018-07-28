package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var isURL = false

func main() {
	colLong := flag.String("long", "", "Name of the column containing the longitude coordinates. If not provided, will try to guess")
	colLat := flag.String("lat", "", "Name of the column containing the latitude coordinates. If not provided, will try to guess")
	delimiter := flag.String("delimiter", ",", "Delimiter character")
	keep := flag.String("keep", "n", "(y/n) If set to \"y\" and the input CSV is an URL, keep the input CSV file on disk")
	threads := flag.Int("threads", 1, "Number of threads (used when converting more than one file)")

	flag.Usage = func() {
		help := "\nOptions:\n" + "  -" + flag.CommandLine.Lookup("delimiter").Name + ": " + flag.CommandLine.Lookup("delimiter").Usage + " (default \"" + flag.CommandLine.Lookup("delimiter").DefValue + "\")" + "\n"
		help += "  -" + flag.CommandLine.Lookup("long").Name + ":      " + flag.CommandLine.Lookup("long").Usage + "\n"
		help += "  -" + flag.CommandLine.Lookup("lat").Name + ":       " + flag.CommandLine.Lookup("lat").Usage + "\n"
		help += "  -" + flag.CommandLine.Lookup("keep").Name + ":      " + flag.CommandLine.Lookup("keep").Usage + " (default \"" + flag.CommandLine.Lookup("keep").DefValue + "\")" + "\n"
		help += "  -" + flag.CommandLine.Lookup("threads").Name + ":   " + flag.CommandLine.Lookup("threads").Usage + " (default \"" + flag.CommandLine.Lookup("threads").DefValue + "\")" + "\n"
		fmt.Fprintf(os.Stderr, "Usage: %s [-options] <input> [output]\n%s", os.Args[0], help)
	}

	flag.Parse()

	var csvFile, jsonFile string

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Error: You need to specify a CSV file. To consult the help, use '-h'.\n")
		os.Exit(1)
	} else if len(flag.Args()) > 2 {
		fmt.Fprintf(os.Stderr, "Error: You can specify a maximum of 2 arguments. To consult the help, use '-h'.\n")
		os.Exit(1)
	} else {
		csvFile = flag.Args()[0]
		if len(flag.Args()) == 2 {
			jsonFile = flag.Args()[1]
		}
	}

	*delimiter = strings.Trim(*delimiter, "'")
	var newDelimiter rune
	if strings.Contains(*delimiter, "\\t") {
		newDelimiter = '\t'
	} else {
		newDelimiter = []rune(*delimiter)[0]
	}

	numGoRoutines := *threads

	filesList := []string{}

	fi, err := os.Stat(csvFile)
	if err != nil {
		if strings.HasPrefix(csvFile, "https://") || strings.HasPrefix(csvFile, "http://") { // case: URL
			isURL = true
			var r io.ReadCloser
			resp, err := http.Get(csvFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Couldn't access the URL: %s.\n", csvFile)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if strings.ToLower(*keep) == "y" || strings.ToLower(*keep) == "yes" {
				parts := strings.Split(csvFile, "/")
				newFile, err := os.Create(parts[len(parts)-1])
				_, err = io.Copy(newFile, resp.Body)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: Couldn't save the CSV file: %s to disk.\n", csvFile)
				}
				csvFile = parts[len(parts)-1]
				r = readFile(csvFile)
				defer r.Close()
			} else {
				r = resp.Body
			}
			convert(r, csvFile, *colLong, *colLat, jsonFile, newDelimiter)
		} else { // case: Wild card
			if !strings.HasSuffix(csvFile, ".csv") {
				csvFile = csvFile + ".csv"
			}
			files, err := filepath.Glob(csvFile)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				filesList = append(filesList, file)
			}
		}

	} else {
		if fi.IsDir() { // case: Directory
			files, err := ioutil.ReadDir(csvFile)
			if err != nil {
				fmt.Println(err)
			}

			for _, f := range files {
				if filepath.Ext(f.Name()) == ".csv" {
					filesList = append(filesList, filepath.Join(csvFile, f.Name()))
				}
			}
		} else { // case: File
			filesList = append(filesList, csvFile)
		}
	}

	if !isURL {
		if strings.ToLower(*keep) == "y" || strings.ToLower(*keep) == "yes" {
			fmt.Println("Info: The option '-keep' is only considered when the input file is an URL.")
		}

		if len(filesList) > 1 {
			if jsonFile != "" {
				fmt.Println("Info: The output file name is not considered when there are multiple files to convert.")
				jsonFile = ""
			}
		} else if len(filesList) == 1 {
			numGoRoutines = 1
		} else {
			fmt.Println("There is no file to convert")
			os.Exit(1)
		}

		var rounds, rest int
		if len(filesList) <= numGoRoutines {
			rounds = 1
			numGoRoutines = len(filesList)
		} else {
			rounds = len(filesList) / numGoRoutines
			rest = len(filesList) % numGoRoutines
		}

		var wg sync.WaitGroup
		for i := 0; i < rounds; i++ {
			for _, f := range filesList[i*numGoRoutines : (i+1)*numGoRoutines] {
				wg.Add(1)
				go func(f string) {
					r := readFile(f)
					defer r.Close()
					convert(r, f, *colLong, *colLat, jsonFile, newDelimiter)
					wg.Done()
				}(f)
			}
			wg.Wait()
		}
		if rest > 0 {
			for _, f := range filesList[rounds*numGoRoutines:] {
				wg.Add(1)
				go func(f string) {
					r := readFile(f)
					defer r.Close()
					convert(r, f, *colLong, *colLat, jsonFile, newDelimiter)
					wg.Done()
				}(f)
			}
			wg.Wait()
		}
	}
}

// readFile opens a file and returns a *File object
func readFile(file string) *os.File {
	f, err := os.Open(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Couldn't find the input CSV file: %s.\n", file)
		os.Exit(1)
	}
	return f
}

// convert converts the data 'r' rom the input CSV file 'inputFile' to an output GeoJSON file 'outputFile'
func convert(r io.Reader, inputFile, colLongitude, colLatitude, outputFile string, delimiter rune) {
	reader := csv.NewReader(r)
	reader.Comma = delimiter

	header, err := reader.Read()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Couldn't read the input CSV file: %s. Cause: %s\n", inputFile, err)
		return
	}

	var indexX, indexY int
	if colLongitude == "" {
		found := false
		for i, v := range header {
			if strings.ToLower(v) == "x" || strings.ToLower(v) == "longitude" || strings.ToLower(v) == "long" || strings.ToLower(v) == "lon" || strings.ToLower(v) == "lng" {
				indexX = i
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "%s: Couldn't determine the column containing the longitude. Specify it using the '-long' option.\n", inputFile)
			return
		}
	} else {
		found := false
		for i, v := range header {
			if strings.ToLower(v) == strings.ToLower(colLongitude) {
				indexX = i
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "%s: Couldn't find column: %s.\n", inputFile, colLongitude)
			return
		}
	}

	if colLatitude == "" {
		found := false
		for i, v := range header {
			if strings.ToLower(v) == "y" || strings.ToLower(v) == "latitude" || strings.ToLower(v) == "lat" {
				indexY = i
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "%s: Couldn't determine the column containing the latitude. Specify it using the '-lat' option.\n", inputFile)
			return
		}
	} else {
		found := false
		for i, v := range header {
			if strings.ToLower(v) == strings.ToLower(colLatitude) {
				indexY = i
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "%s: Couldn't find column: %s.\n", inputFile, colLatitude)
			return
		}
	}

	if indexX < indexY {
		header = append(header[:indexX], header[indexX+1:]...)
		header = append(header[:indexY-1], header[indexY:]...)
	} else {
		header = append(header[:indexY], header[indexY+1:]...)
		header = append(header[:indexX-1], header[indexX:]...)
	}

	var buffer bytes.Buffer

	buffer.WriteString(`{
		"type": "FeatureCollection",
		"crs": { "type": "name", "properties": { "name": "urn:ogc:def:crs:OGC:1.3:CRS84" } },                                                                  
		"features": [
	`)

	// Read the rest of the file
	content, err := reader.ReadAll()

	if len(content) == 0 {
		fmt.Fprintf(os.Stderr, "The input CSV file %s is empty. Nothing to convert.\n", inputFile)
		return
	}

	for i, d := range content {
		coordX := d[indexX]
		coordY := d[indexY]
		// Only convert the row if both coordinates are available
		if coordX != "" && coordY != "" {
			buffer.WriteString(`{ "type": "Feature", "properties": {`)

			if indexX < indexY {
				d = append(d[:indexX], d[indexX+1:]...)
				d = append(d[:indexY-1], d[indexY:]...)
			} else {
				d = append(d[:indexY], d[indexY+1:]...)
				d = append(d[:indexX-1], d[indexX:]...)
			}
			for j, y := range d {

				buffer.WriteString(`"` + header[j] + `":`)
				_, fErr := strconv.ParseFloat(y, 32)
				_, bErr := strconv.ParseBool(y)
				if fErr == nil {
					buffer.WriteString(y)
				} else if bErr == nil {
					buffer.WriteString(strings.ToLower(y))
				} else {
					buffer.WriteString((`"` + y + `"`))
				}
				//end of property
				if j < len(d)-1 {
					buffer.WriteString(",")
				}
			}
			//end of object of the array
			buffer.WriteString(`}, "geometry": { "type": "Point", "coordinates": [` + coordX + `, ` + coordY + `]} }`)
			if i < len(content)-1 {
				buffer.WriteString(",\n")
			}
		}
	}
	buffer.WriteString(`]
}`)
	rawMessage := json.RawMessage(buffer.String())
	var output string
	ext := ".geojson"
	if outputFile == "" {
		if isURL {
			parts := strings.Split(inputFile, "/")
			output = strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(inputFile)) + ext
		} else {
			output = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ext
		}
	} else if outputFile == strings.TrimSuffix(outputFile, ext) { // If no extension provided
		output = outputFile + ext
	} else {
		output = outputFile
	}
	if err := ioutil.WriteFile(output, rawMessage, os.FileMode(0644)); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create the GeoJSON file: %s.\n", output)
		return
	}
	fmt.Fprintf(os.Stderr, "The GeoJSON file %s was successfully created.\n", output)
}
