# Unit Converter API - Go
Unit converter is an ideal code example with simple business, light weight code that help us to compare programming languages, technologies and platform. It is a stateless Microservices with no dependencies such as database, messaging broker and cache. With less focus on development of features, we could focus on the full lifecycle of the Microservice from initiation, development to deployment and E2E tests.
Following
- [Quarkus - Java](https://github.com/nikfarjam/quarkus-temperature-converter)
- [Spring Boot](https://github.com/nikfarjam/UnitConverterServer)
- [JavaScript](https://github.com/nikfarjam/celsius-converter)
- [Python](https://github.com/nikfarjam/UnitConvertor-python)

This project is written in Go.

#### Description
An RESTful API that runs on port 9090 and converts temperatures between Fahrenheit (°F) and Celsius (°C).
**Feature**
- Convert Fahrenheit to Celsius and vice versa.
- Handle invalid input gracefully.
- Provide clear and concise JSON responses.

#### How to run
```bash
cd unit-converter/
go mod tidy
go run .
```

#### API Endpoints

Endpoint: */converter*  
Method: *POST*  
Request Payload:
```json
{
    "value": 15,
    "from": "celsius",
    "to": "Fahrenheit"
}
```
Response:
```json
{
    "value": 59,
    "Unit": "FAHRENHEIT"
}
```

#### Example Usage
Fahrenheit to Celsius conversion using cURL
```bash

curl -X POST -H "Content-Type: application/json" -d '{"value": 15,"from": "celsius","to": "Fahrenheit"}' http://localhost:9090/converter
```
