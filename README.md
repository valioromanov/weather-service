# weather-service

A serverless weather service that fetches forecast data from Open-Meteo, caches it in DynamoDB, and exposes it via a REST API (API Gateway). Built in Go, deployed via Terraform.

## API Endpoint
### `GET /weather?lat={latitude}&lon={longitude}&date={date}`

### ✅ Query Parameters

| Parameter | Type     | Required | Description                                                  |
|-----------|----------|----------|--------------------------------------------------------------|
| `lat`     | `float`  | Yes      | Latitude of the location (e.g., `42.6975`)                  |
| `lon`     | `float`  | Yes      | Longitude of the location (e.g., `23.3241`)                 |
| `date`    | `string` | No       | Date in `YYYY-MM-DD` format (defaults to today)             |

---

### 📤 Example Response

```json
{
    "date": "2025-07-11",
    "latitude": "43.6875",
    "longitude": "23.3125",
    "temperature": 26.7,
    "uvIndex": 7.05,
    "rainProbability": 0
}
```

## Api Logic
1. Cache Check: The Lambda function first checks DynamoDB for a cached forecast using lat+lon+date as the key.
2. API Fallback: If not cached or expired, it fetches fresh data from Open-Meteo.
3. Cache Store: The new forecast is stored in DynamoDB with a TTL (Time-To-Live).
4. Response: Returns the weather data to the user.

## Error Responses

| HTTP Status | Message                                   |
| ----------- | ----------------------------------------- |
| 400         | Missing or invalid query parameters       |
| 404         | Weather data for the given date not found |
| 500         | Internal server or external API error     |

## Build and deploy
Before deploying, you should have AWS CLI configured
- make build       # Build Go binary and zip
- make deploy      # Deploy via Terraform

## Testing
- make tests             # Make only tests
- make testsWithCoverage # Make tests with coverage, generate an HTML file that will be opened, where you can check which lines are not covered

## Requirements
- Go 1.21+
- AWS CLI (configured)
- Terraform >= 1.3
- Ginkgo / GoMock for testing

## Running

Service now runs at [10vfl9na9l.execute-api.eu-west-1.amazonaws.com](10vfl9na9l.execute-api.eu-west-1.amazonaws.com) </br>
Example request: https://10vfl9na9l.execute-api.eu-west-1.amazonaws.com/weather?lat=43.6875&lon=23.3125&date=2025-07-11

