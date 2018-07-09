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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	colLong := flag.String("long", "", "Name of the column containing the longitude coordinates. If not provided, will try to guess")
	colLat := flag.String("lat", "", "Name of the column containing the latitude coordinates. If not provided, will try to guess")
	delimiter := flag.String("delimiter", ",", "Delimiter character")
	keep := flag.String("keep", "n", "(y/n) If set to 'y' and the input CSV is an URL, keep the input CSV file on disk")

	flag.Usage = func() {
		help := "\nOptions:\n" + "  -" + flag.CommandLine.Lookup("delimiter").Name + ": " + flag.CommandLine.Lookup("delimiter").Usage + " (default \"" + flag.CommandLine.Lookup("delimiter").DefValue + "\")" + "\n"
		help += "  -" + flag.CommandLine.Lookup("long").Name + ":      " + flag.CommandLine.Lookup("long").Usage + "\n"
		help += "  -" + flag.CommandLine.Lookup("lat").Name + ":       " + flag.CommandLine.Lookup("lat").Usage + "\n"
		help += "  -" + flag.CommandLine.Lookup("keep").Name + ":      " + flag.CommandLine.Lookup("keep").Usage + " (default \"" + flag.CommandLine.Lookup("keep").DefValue + "\")" + "\n"
		fmt.Fprintf(os.Stderr, "Usage: %s [-options] <input> [output]\n%s", os.Args[0], help)
	}

	flag.Parse()

	var csvFile, jsonFile string

	if len(flag.Args()) == 0 {
		fmt.Println("Error: You need to specify a CSV file. To consult the help, use '-h'.")
		os.Exit(1)
	} else if len(flag.Args()) > 2 {
		fmt.Println("Error: You can only specify 2 arguments. To consult the help, use '-h'.")
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

	var r io.ReadCloser

	// If the input CSV is a URL
	if isValidURL(csvFile) {
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

	} else { // If is a file
		if strings.ToLower(*keep) == "y" || strings.ToLower(*keep) == "yes" {
			fmt.Println("Info: The option '-keep' is only considered when the input file is an URL.")
		}
		r = readFile(csvFile)
		defer r.Close()
	}

	convert(r, csvFile, *colLong, *colLat, jsonFile, newDelimiter)
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
		os.Exit(1)
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
			fmt.Println("Couldn't determine the column containing the longitude. Please specify it using the '-long' option.")
			os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Couldn't find column: %s.\n", colLongitude)
			os.Exit(1)
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
			fmt.Println("Couldn't determine the column containing the latitude. Please specify it using the '-lat' option.")
			os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Couldn't find column: %s.\n", colLatitude)
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
		os.Exit(1)
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
		if isValidURL(inputFile) {
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
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "The GeoJSON file %s was successfully created.\n", output)
}

// isValidURL checks is a string is a valid URL
func isValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return true
}
