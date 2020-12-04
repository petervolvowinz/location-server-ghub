/*
 * Copyright Volvo Cars (c) 2020. Author: Peter Winzell, peter.winzell@volvocars.com, Sunnyvale
 */


// REST APIS
var RETRIEVE_API = 'retrieve?search=';

function FetchURL(api){
    hostname = location.hostname
    if (hostname === "localhost"){
        return "http://"+ location.host + "/" + api
    }else {
        return "https://" + hostname +"/" + api
    }
}