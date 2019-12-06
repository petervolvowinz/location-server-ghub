# location-server

A general "real-time" location server based on GPS position. The principle is to have a server that is able for a certain time interval handle GPS position based requests. The entries are disposed of after time T and clients are identified using universal identifiers. The server is currently configured to handle two types of objects identified as vehicle and bike but the server itself does not handle any logic apart from GPS position and time. For each entry,  clients can/should add a payload which is a string of any content - in our examples we have been using JSON.

The idea behind using universal identifiers is that the clients handle and decide how they are identified and for how long. It also anonymizes and removes the need to use vehicle identification numbers, but not only that it also gives other types of "objects" means of identification. There is no persistence once time T has passed the data cannot be retrieved again.

There are currently 4 apis:
```
  /addposnoret
  Adds the new position , does not return a response.
  parameters:
  gps ->
  Ex:
  gps={"Location":
   {"lat":37.384997,
    "lng":-122.03452399999999,
    "accuracy":1,
    "payload":{\"ambientemp\":7.4,\"cabintemp\":23.3,\"drivertemp\":25.7,\"parkingspots\":36\"}"},
     "Gpsobject":1,
     "UUID":"c592fe16-5b4d-4fd7-ab58-7441b20299cd",
     "timestamp":1
     }
   } 
  
  /addposition
  Adds the new position, returns gps positions for objects within d meters during a t timespan.
  parameters:
  gps ->
  
  /retrieve
  
  /version
  returns the current deployed server version as a string
  ```
  
  The json that client needs to parse in Go:
  ```
  type Warninglst struct{
	  Warnings []GPSLocation	`json:"warnings"`
  }

  type GPSLocation struct{
	  Location Locationdata `json:"Location"`
	  Gpsobject int	  `json:"Gpsobject"`
	  UUID uuid.UUID       `json:"UUID"`
	  Timestamp int64    `json:"timestamp"`
  }

  type Locationdata struct {
	  // Payload   string  `json:"payload"`
	  Latitude  float64 `json:"lat"`
	  Longitude float64 `json:"lng"`
	  Accuracy  float64 `json:"accuracy"`
	  Payload   string  `json:"payload"`
  }
  
  ```
