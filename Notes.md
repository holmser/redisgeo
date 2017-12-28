# Athlinks API Notes
- Need coordinates for race to make geo search useful
- For sign ups only need races with 0 finishers

- [REST Api Design Notes](http://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api#caching)
 - caching is of particular interest

- add `expand` parameter?

- [use ISO 3166 Country Codes] Races](http://www.nationsonline.org/oneworld/country_code_list.htm)

Current Response:
```json
"City": "Miami",
"CountryName": "US",
"CountryID": "US",
"StateProvID": "US_FL",
"StateProvAbbrev": "FL",
"StateProvName": null,
```

More desirable response:
- Coordinates should be the actual location of the race, starting line if possible, or geographic center of course map?
- This provides all information needed to determine location from ISO 3166 lookup.  Right now we are providing both.

```json
{
  "Location" : {
    "City" : "Miami",
    "StateProvID" : "US_FL",
    "Coordinates":{
      "Latitude": 25.7617,
      "Longitude": -80.1918
    }
  }
}
```
- Get rid of "DistUnit", all distances in meters.  Conversion from there.
