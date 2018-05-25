package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	csvFile := flag.String("in", "", "Input CSV file")
	colLong := flag.String("long", "", "Name of the column containing the longitude coordinate. If not provided I will try to guess")
	colLat := flag.String("lat", "", "Name of the column containing the latitude coordinate. If not provided I will try to guess")
	delimiter := flag.String("delimiter", ",", "Delimiter character")
	jsonFile := flag.String("out", "", "Output GeoJSON file (extension will be added if omitted)")

	flag.Parse()

	if *csvFile == "" {
		fmt.Println("Error: You need to specify a CSV file. Use the '-in' option. To consult the help, use '-h'.")
		os.Exit(1)
	}
	inFile, err := os.Open(*csvFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Couldn't find the input CSV file: %s.\n", *csvFile)
		os.Exit(1)
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)

	*delimiter = strings.Trim(*delimiter, "'")
	if strings.Contains(*delimiter, "\\t") {
		reader.Comma = '\t'
	} else {
		reader.Comma = []rune(*delimiter)[0]
	}

	content, err := reader.ReadAll()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Couldn't read the input CSV file: %s.\n", *csvFile)
		os.Exit(1)
	}

	if len(content) <= 1 {
		fmt.Fprintf(os.Stderr, "The input CSV file %s is empty. Nothing to convert.\n", *csvFile)
		os.Exit(1)
	}

	header := make([]string, 0)
	for _, headE := range content[0] {
		header = append(header, headE)
	}

	var indexX, indexY int
	if *colLong == "" {
		found := false
		for i, v := range header {
			if strings.ToLower(v) == "x" || strings.ToLower(v) == "longitude" || strings.ToLower(v) == "long" || strings.ToLower(v) == "lon" {
				indexX = i
				found = true
			}
		}
		if !found {
			fmt.Println("Couldn't determine the column containing the longitude. Please specify it using the '-long' option.")
			os.Exit(1)
		}
	} else {
		for i, v := range header {
			if strings.ToLower(v) == strings.ToLower(*colLong) {
				indexX = i
			}
		}
		if indexX == 0 {
			fmt.Fprintf(os.Stderr, "Couldn't find column: %s.\n", *colLong)
			os.Exit(1)
		}
	}

	if *colLat == "" {
		found := false
		for i, v := range header {
			if strings.ToLower(v) == "y" || strings.ToLower(v) == "latitude" || strings.ToLower(v) == "lat" {
				indexY = i
				found = true
			}
		}
		if !found {
			fmt.Println("Couldn't determine the column containing the latitude. Please specify it using the '-lat' option.")
			os.Exit(1)
		}
	} else {
		for i, v := range header {
			if strings.ToLower(v) == strings.ToLower(*colLat) {
				indexY = i
			}
		}
		if indexY == 0 {
			fmt.Fprintf(os.Stderr, "Couldn't find column: %s.\n", *colLat)
			os.Exit(1)
		}
	}

	if indexX < indexY {
		header = append(header[:indexX], header[indexX+1:]...)
		header = append(header[:indexY-1], header[indexY:]...)
	} else {
		header = append(header[:indexY], header[indexY+1:]...)
		header = append(header[:indexX-1], header[indexX:]...)
	}
	//Remove the header row
	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString(`{
"type": "FeatureCollection",
"crs": { "type": "name", "properties": { "name": "urn:ogc:def:crs:OGC:1.3:CRS84" } },                                                                  
"features": [
`)
	for i, d := range content {
		buffer.WriteString(`{ "type": "Feature", "properties": {`)
		coordX := d[indexX]
		coordY := d[indexY]
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
	buffer.WriteString(`]
}`)
	rawMessage := json.RawMessage(buffer.String())
	var output string
	ext := ".geojson"
	if *jsonFile == "" {
		output = strings.TrimSuffix(*csvFile, filepath.Ext(*csvFile)) + ext
	} else if *jsonFile == strings.TrimSuffix(*jsonFile, ext) { // If no extension provided
		output = *jsonFile + ext
	} else {
		output = *jsonFile
	}
	if err := ioutil.WriteFile(output, rawMessage, os.FileMode(0644)); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create the GeoJSON file: %s.\n", output)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "The GeoJSON file %s was successfully created.\n", output)
}
