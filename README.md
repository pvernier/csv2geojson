[![Build Status](https://travis-ci.com/pvernier/csv2geojson.svg?branch=master)](https://travis-ci.com/pvernier/csv2geojson)

# csv2geojson

A CLI to convert CSV to GeoJSON from a file, folder or URL.

## Download

Get the binaries (64 bits) for Linux, Windows and macOS [here](https://github.com/pvernier/csv2geojson/releases).

## Usage

```
$csv2geojson -h
Usage: csv2geojson [-options] <input> [output]

Options:
  -delimiter: Delimiter character (default ",")
  -long:      Name of the column containing the longitude coordinates. If not provided, will try to guess
  -lat:       Name of the column containing the latitude coordinates. If not provided, will try to guess
  -keep:      (y/n) If set to "y" and the input CSV is an URL, keep the input CSV file on disk (default "n")
  -threads:   Number of threads (used when converting more than one file) (default "1")
  -suffix:    Suffix to add to the name of output GeoJSON file(s)

```

## Examples

### Convert a single file

```
# Simplest way
$csv2geojson data.csv
The GeoJSON file data.geojson was successfully created.
```

In this case csv2geojson will try to guess which fields contain the longitude and latitude coordinates. Also, if **[output]** is omitted, the output GeoJSON file gets the same name as the input CSV file.

```
# Complete way
$csv2geojson -delimiter ; -long field4 -lat field3 data_fr.csv pois
The GeoJSON file pois.geojson was successfully created.
```

If the fields of the input CSV file are not separated by commas, use the *-delimiter* option. If the fields containing the longitude and latitude don't have a explicit name, use the *-long* and *-lat* options. Explicit names are: 'longitude', 'long', 'lon', 'lng' and 'x' for the longitude. 'latitude', 'lat' and 'y' for the latitude (case insensitive).

```
$csv2geojson -delimiter \t ..\coords_tab.txt
The GeoJSON file ..\coords_tab.geojson was successfully created.
```

csv2geojson can also convert tab separated text files. Options can be entered in any order. The input CSV file doesn't need to be in the current folder. It can be a relative or absolute path.

### Convert from a URL

```
$csv2geojson -keep y https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_month.csv
The GeoJSON file all_month.geojson was successfully created.
```

*-keep y* saves the downloaded CSV. By default, it is not saved.

### Convert all the CSV files inside a folder

```
# The folder 'data' contains 5 CSV files: zone_A1, zone_A2, zone_A3, zone_B1 and zone_B2.

$csv2geojson -threads 3 data
The GeoJSON file data\zone_A3.geojson was successfully created.
The GeoJSON file data\zone_A2.geojson was successfully created.
The GeoJSON file data\zone_A1.geojson was successfully created.
The GeoJSON file data\zone_B1.geojson was successfully created.
The GeoJSON file data\zone_B2.geojson was successfully created.
```
When converting multiple files, using the option *-threads* can make it faster. In this case, the order in which the files are converted is not
guaranteed to always be the same.

### Convert a subset of CSV files inside a folder

```
$csv2geojson data\*B*
The GeoJSON file data\zone_B1.geojson was successfully created.
The GeoJSON file data\zone_B2.geojson was successfully created.
```


## Alternatives

* [csv2geojson](https://github.com/mapbox/csv2geojson) (Javascript)
* [ogr2ogr2](http://www.gdal.org/ogr2ogr.html)

## Inspiration

 * [Ahmad-Magdy/CSV-To-JSON-Converter](https://github.com/Ahmad-Magdy/CSV-To-JSON-Converter)
 * [Golang - Read CSV/JSON from URL](https://gist.github.com/stupidbodo/71f2b164744a18a18e74)
