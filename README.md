# location-server

A "real-time" location server based on GPS position. The principle is to implement a server solution that is able for a certain time interval handle GPS position based requests. The client request is to identify objects such as vehicles and vrus within a certain distance d and time depth t from when the request was recieved. The entries are disposed of after time T. Clients are identified using universal identifiers. The server is currently configured to handle two types of objects identified as vehicle and bike. The server itself does not handle any logic apart from GPS position and time. For each entry,  clients can/should add a payload which is a string of any content - in our examples we have been using a JSON payload.

The idea behind using universal identifiers is that the clients handle and decide how they are identified and for how long. It also anonymizes and removes the need to use vehicle identification numbers, but not only that it also gives other types of "objects" means of identification. There is no persistence, and once time T has passed the data cannot be retrieved again.

## Client apis
There are currently 4 apis:
```
  /addposnoret
  Adds the new position , does not return a response.
  parameters:
  gps ->
  Ex:
  {"Location":
   {"lat":37.384997,
    "lng":-122.03452399999999,
    "accuracy":1,
    "payload":{\"ambientemp\":7.4,\"cabintemp\":23.3,\"drivertemp\":25.7,\"parkingspots\":36\"}"},
     "Gpsobject":1,
     "UUID":"c592fe16-5b4d-4fd7-ab58-7441b20299cd",
     "timestamp":1
     }
   } 
  
  => https:/serverurl/addposnoret?gps=<json>
  
  /addposition
  Adds the new position, returns gps positions for objects within d meters during a t timespan.
  parameters:
  gps ->
  type:string , json format.
  
  Ex:
  {"Location":
   {"lat":37.384997,
    "lng":-122.03452399999999,
    "accuracy":1,
    "payload":{\"ambientemp\":7.4,\"cabintemp\":23.3,\"drivertemp\":25.7,\"parkingspots\":36\"}"},
     "Gpsobject":1,
     "UUID":"c592fe16-5b4d-4fd7-ab58-7441b20299cd",
     "timestamp":1
     }
   }
   timespan ->
   type:int, seconds
   Ex:
   5
   distance ->
   type:int, meters
   200
   
   => https:/serverurl/addposnoret?gps=<json>&timespan=5,distance=200
   <=
    {"warnings":[{"Location":{"lat":37.385997,"lng":-122.03923636959762,"accuracy":1,
    	"payload":"{\"ambientemp\":5.9,\"cabintemp\":18.5,\"drivertemp\":21.8,\"parkingspots\":94,\"vehicleid\":\"1\"}"},
    	"Gpsobject":1,
    	"UUID":"310dcda3-86c7-447b-a104-4061924dbc94",
    	"timestamp":1575672675636676000
    	}]
    } 

  
  /retrieve
  Retrieves all registered gps objects within distance d from position p during timespan t.
  parameters: <TO BE ADDED>
  
  /version
  returns the current deployed server version as a string
  ```
  
  Go structs as an example. The structs are JSON encoded and then passed to the client in the http reponse. Corresponding data structures is needed on the client side to encode the reponse:
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
	  string  `json:"payload"` // Payload   
	  Latitude  float64 `json:"lat"`
	  Longitude float64 `json:"lng"`
	  Accuracy  float64 `json:"accuracy"`
	  Payload   string  `json:"payload"`
  }
  
  ```
## Deployment

The server is built using a Dockerfile which is subsequently used for deployment:

```
FROM golang:alpine AS build-env

ADD  ./main.go /go/src
ADD ./vendor/queue /go/src/queue

RUN apk add --no-cache git
RUN go get -d -v github.com/google/uuid
RUN go get -d -v github.com/sirupsen/logrus
RUN go get -d -v github.com/emirpasic/gods/lists/doublylinkedlist


RUN cd /go/src && CGO_ENABLED=0 go build -o locationserver

FROM alpine
ADD static /app/static
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/locationserver /app/

RUN addgroup --gid 1000 go && adduser -D -G go -u 100 go
RUN chown go ./locationserver
USER go

EXPOSE 8081
ENTRYPOINT  ./locationserver
```
Furthermore, the github account is currently tied to quay.io which builds the Docker image and then spinnaker and kubernetes handles the deployment to the cloud. 
Tagging scheme, which is used as the image build trigger.
```
git tag release-x.y.z-demo
```
