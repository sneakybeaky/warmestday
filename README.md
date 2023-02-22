A simple application that returns the warmest day over the next 7 for a given location in Europe.

The project layout is mostly based on Alex Edward's [Let's Go](https://lets-go.alexedwards.net) structure.

To run the server from the root of the project

    go run  warmestday/cmd/web --appid <open weather map appid>

## Endpoints

### Ping
To check the server is reachable

    curl http://localhost:4000/ping
    OK

### Get the warmest day for a given location

    http://localhost:4000/summary?lat={lat}&lon={lon}

For example, for London

     curl "http://localhost:4000/summary?lat=50.8229&lon=0.1363"
    {"WarmestDay":"2023-02-22"}

If the location is outside Europe you will get a bad request response code (400).

For example, for New York

    curl http://localhost:4000/summary?lat=40.7128&lon=74.0060
    Location must be inside Europe

Latitude and Longitude must have at most 6 decimal places. Any more is treated as a bad request

For example

    curl -v http://localhost:4000/summary?lat=40.7128&lon=74.0060888888888
    Longitude has too many decimal places - no more than 6 allowed