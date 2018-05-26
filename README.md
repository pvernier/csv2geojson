# csv2geojson

A CLI to convert CSV to GeoJSON from a file or URL.

## Download

Get the binaries (64 bits) for Windows and Linux [here](https://github.com/pvernier/csv2geojson/releases).

## Usage

```
$csv2geojson -h
Usage of csv2geojson:
  -delimiter string
        Delimiter character (default ",")
  -in string
        Input CSV file
  -lat string
        Name of the column containing the latitude coordinate. If not provided I will try to guess
  -long string
        Name of the column containing the longitude coordinate. If not provided I will try to guess
  -out string
        Output GeoJSON file (extension will be added if omitted)

```

## Examples

```
# Simplest way
$csv2geojson -in data.csv
The GeoJSON file data.geojson was successfully created.
```
The *-in* option is the only mandatory option. In this case csv2geojson will try to guess which fields contain the longitude and latitude coordinates. Also, if the *-out* option is omitted, the output GeoJSON file gets the same name as the input CSV file.

```
# Complete way
$csv2geojson -in pois_09.csv -delimiter ; -long field4 -lat field5 -out pois
The GeoJSON file pois.geojson was successfully created.
```
If the fields of the input CSV file are not separated by commas, use the *-delimiter* option. If the fields containing the longitude and latitude don't have a explicit name, use the *-long* and *-lat* options. Explicit names are: longitude', 'long', 'lon' and 'x' for the longitude. 'latitude', 'lat' and 'y' for the latitude (case insensitive).

```
# Convert from a URL
$csv2geojson -in https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_month.csv
The GeoJSON file all_month.geojson was successfully created.
```

```
$csv2geojson -delimiter \t -in ..\coords_tab.txt
The GeoJSON file ..\coords_tab.geojson was successfully created.
```
csv2geojson can also convert tab separated text files. Options can be entered in any order. The input CSV file doesn't need to be in the current folder. It can be a relative or absolute path.

## Alternatives

* [csv2geojson](https://github.com/mapbox/csv2geojson) (Javascript)
* [ogr2ogr2](http://www.gdal.org/ogr2ogr.html)

## Inspiration

 * [Ahmad-Magdy/CSV-To-JSON-Converter](https://github.com/Ahmad-Magdy/CSV-To-JSON-Converter)
 * [Golang - Read CSV/JSON from URL](https://gist.github.com/stupidbodo/71f2b164744a18a18e74)
 * [How to check if a string is a URL](https://golangcode.com/how-to-check-if-a-string-is-a-url/)
