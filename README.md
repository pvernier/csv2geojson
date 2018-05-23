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

## Alternatives

* [csv2geojson](https://github.com/mapbox/csv2geojson) (Javascript)
* [ogr2og2](http://www.gdal.org/ogr2ogr.html)

## Inspiration

 * [Ahmad-Magdy/CSV-To-JSON-Converter](https://github.com/Ahmad-Magdy/CSV-To-JSON-Converter)
