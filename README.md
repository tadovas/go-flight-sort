## Flight path calculation microservice
Given list of user direct flights as input, service calculates initial and final destinations if possible.
### To run service
`go run cmd/main.go`

It will start serving http requests at 8080 port

### API documentation
#### POST /calculate
##### Request
Expected Content-type: application/json

Expects json body with single field `flights` which contains list of users flights.
Example:
```json
{
    "flights": [
        {
            "source": "ATL",
            "dest": "EWR"
        }, { 
            "source": "SFO", 
            "dest": "ATL"
        }
    ]
}
```
Each flight object consists of `source` and `dest` airports.
##### Response
Content-type: application/json

A json object containing single `flight` field with initial `source` and final `destination` airports for given user input.
Example:
```json
{
    "flight": {
        "source": "SFO",
        "dest": "EWR"
    }
}
```
#### Service errors
In case of http code other than 200 OK, Service returns generic error json object with single field `error` of string type with error description and appropriate http error code.
Example:
```json
{
    "error": "input sanitization: flights[0] source is empty"
}
```
##### Http error codes
| Http code  | Description                                                         |
| ---------- |---------------------------------------------------------------------|
| 200        | Everything went well. Normal response with flight should be present |
| 405 | Only HTTP POST method is allowed                                    |
 | 415 | Only application/json body is allowed |
 | 400 | Cannot parse JSON body |
 | 422 | Either input validation failed or flight sorting function returned error |

