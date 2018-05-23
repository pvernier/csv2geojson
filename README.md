# csv2geojson

A CLI to convert CSV to GeoJSON

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
The *-in* option is the only mandataory option. In this case csv2geojson will try to guess which fields contain the longitude and latitude coordinates. Also, if the *-out* option is omitted, the output GeoJSON file gets the same name as the input CSV file.

```
# Complete way
$csv2geojson -in plane_trips_coords2.csv -delimiter ; -long field4 -lat field5 -out pois
The GeoJSON file pois.geojson was successfully created.
```
If the CSV fields are not separated by commas, use the *-delimiter* option. If the fields containing the longitude and latitude don't have a explicit name, use the *-long* and *-lat* options. Explicit names are: longitude', 'long', 'lon' and 'x' for the longitude. 'latitude', 'lat' and 'y' for the latitude (case insensitive).

## Alternatives

* [csv2geojson](https://github.com/mapbox/csv2geojson) (Javascript)
* [ogr2og2](http://www.gdal.org/ogr2ogr.html)

## Inspiration

 * [Ahmad-Magdy/CSV-To-JSON-Converter](https://github.com/Ahmad-Magdy/CSV-To-JSON-Converter)
