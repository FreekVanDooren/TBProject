# Usage
## Prerequisites
Go version 1.15 
## Test
```
go test ./...
```
#### Artifacts
Running tests will create `test_data`-folders
## Build
```
go build
```
## Run
After build step
```
./hello
```
Service is listening on port 8080
#### Artifacts
Running the service will create folders `data` and `logs`
# Endpoints
| Method | Path |
| ------ | ---- |
| GET | '/' |
| GET | '/history' |
| GET | '/primes/{number:[0-9]+}' |
| GET | '/messages' |
| POST | '/messages' |

## Example requests
```
curl localhost:8080/
# Home!% 
curl localhost:8080/history
# {"requests":[{"number":9002,"count":12}]}
curl localhost:8080/primes/9002
# {"isPrime":false,"message":"No, and we already told you so!"}
curl localhost:8080/messages
# {"messages":[{"lowerLimit":3,"message":"No, and we already told you so!"},{"lowerLimit":0,"message":"No"}]}
curl -H "application/json" -X POST localhost:8080/messages -d "{\"messages\":[{\"lowerLimit\":0,\"message\":\"No no no no no...\"}]}"
#
```

# TODO
* Create Frontend
* Make positive feedback messages adaptable
* Validate POST to /messages to contain only benign data  
* Figure out if there are any memory leaks
* Rename to something else than "hello"
* Figure out how to write cleaner (test) code
* Stop tests from creating folders or clean them up afterwards.

# Assignment
Live coding
We have had an increasing demand of knowing if a number X is a prime number or not! It has exhausted our recourses on tables of primes in ancient lore. We also get annoying request on the same number, over and over, mostly when it clearly is not a prime. This situation needs to be handled and handled quickly!

Therefore, we want you to create a service, preferable in go and most possible implemented in microservice style that can decide if a number sent to it is a prime or not! 
We also want you to keep track of the requests coming in, since we really hate those annoying repeating questions and want to ensure that after for example a third call on the same number, we want the service to reply with a more explicit answer then “no” more in the terms like “No, and we already told you so!”. And we want the service to be as dependable as the books and charts we currently use!


It would be great to call your service with curl, but it is up to you how you handle request.
If you create a possibility to change the “No, and we already told you so!” with a call to your code, that is indeed nice! Maybe even change the frequency it will be displayed?
It would also be really cool to call your service and get all prime the service knows of in a specific number range!
For statistical purpose is would also be interesting to know all numbers that have been asked for in a range, maybe how many times?
If you opt to create a UI, we would be happy!

Constrains:
More or less none, use anything you need, to get the job done!
For initial iteration your service would only need to work on numbers up to 100.
Can you have it up and running in less than 40 minutes?
