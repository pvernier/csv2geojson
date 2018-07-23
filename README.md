# csv2geojson

A CLI to convert CSV to GeoJSON from a file or URL.

## Download

Get the binaries (64 bits) for Windows and Linux [here](https://github.com/pvernier/csv2geojson/releases).

## Usage

```
$csv2geojson -h
Usage: csv2geojson [-options] <input> [output]

Options:
  -delimiter: Delimiter character (default ",")
  -long:      Name of the column containing the longitude coordinates. If not provided, will try to guess
  -lat:       Name of the column containing the latitude coordinates. If not provided, will try to guess
  -keep:      (y/n) If set to "y" and the input CSV is an URL, keep the input CSV file on disk (default "n")

```

## Examples

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
# Convert from a URL and keep the CSV file
$csv2geojson -keep y https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_month.csv
The GeoJSON file all_month.geojson was successfully created.
```

```
$csv2geojson -delimiter \t ..\coords_tab.txt
The GeoJSON file ..\coords_tab.geojson was successfully created.
```

csv2geojson can also convert tab separated text files. Options can be entered in any order. The input CSV file doesn't need to be in the current folder. It can be a relative or absolute path.

## Alternatives

* [csv2geojson](https://github.com/mapbox/csv2geojson) (Javascript)
* [ogr2ogr2](http://www.gdal.org/ogr2ogr.html)

## Inspiration

 * [Ahmad-Magdy/CSV-To-JSON-Converter](https://github.com/Ahmad-Magdy/CSV-To-JSON-Converter)
 * [Golang - Read CSV/JSON from URL](https://gist.github.com/stupidbodo/71f2b164744a18a18e74)
