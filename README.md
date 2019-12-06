# location-server

A general "real-time" location server based on gps position. The principle is to have server that are able for a certain time interval handle gps position based requests. The entries are disposed after time T and clients are identified using universal idenftifiers. The server currently is configured to handle two types of objects identfied as vehicle and bike but the server itself does not handle any logic apart from gps positition and time. For each entry we clients can/should add a payload which is basically string of any content - in our examples we have been using json.

The idea behind using universal identifiers is that the clients handle and decide how they are identified, or recognizable if you will, from other objects. It also anonomises 
