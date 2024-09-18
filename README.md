To run the application:
```
go run cmd/server/main.go
```

While the app is running, send a CURL request to it:
```
curl http://localhost:8080/getShortForecast\?lat\=34.03\&long\=-118.15
```

Sample Success Response:
```
{
    "short_forecast":"Partly Cloudy",
    "characterization":"Moderate"
}  
```