# Luno Market Data Application


This is a CLI style application that uses golang to connect to a webservice to extract crypto market data 

This application connects to the Luno Website (https://www.luno.com/) and it can retrieve Market Data from their website, via WebSocket interface.


# Pre-requisites
1. Golang version 1.24

# To install
1. Create a folder, for the source code
2. Clone the repo
3. Run go install

# Configuration
Create a director, .luno, your user director, and add a file keys.json, with the following content:

```json
{
"key": "",
"secret": ""
}
```


Run application
