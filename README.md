# location-server

A general "real-time" location server based on GPS position. The principle is to have a server that is able for a certain time interval handle GPS position based requests. The entries are disposed of after time T and clients are identified using universal identifiers. The server is currently configured to handle two types of objects identified as vehicle and bike but the server itself does not handle any logic apart from GPS position and time. For each entry,  clients can/should add a payload which is a string of any content - in our examples we have been using JSON.

The idea behind using universal identifiers is that the clients handle and decide how they are identified and for how long. It also anonymizes and removes the need to use vehicle identification numbers, but not only that it also gives other types of "objects" means of identification. There is no persistence once time T has passed the data cannot be retrieved again.

There are currently 4 apis:
```
  /addposnoret
  
	/addposition
  /retrieve
  /version
  

