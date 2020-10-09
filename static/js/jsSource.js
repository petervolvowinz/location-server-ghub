var carCount = 0;
var bicycleCount = 0;
var wholeShareLocations = "";
var position = [];
var carCount = 0;
var bicycleCount = 0;
var alreadyExsist = 0;
var ambientTempInfo = [];
var cabinTempInfo = [];
var driverTempInfo = [];
var UUIDorigin = [];
var filterType = ['match', ['get', 'Icontype'], ["car","bicycle"], true, false];
var layerTypeObject = document.getElementById('layerTypeObject');
var distanceItem = document.getElementById('distance');

var lat = 0
var lon = 0

// Utility function checks if we have a uuid in sharedLocations array
function uuidMemberTest(sharedLocations,uuid){
  let found = false
  for (let el of sharedLocations){
    if (el.properties['UUID'] === uuid){
      found = true
      break
    }
  }
  return found
}

//  Increase counter of icontype
function countCarsAndBikes(icontype) {
  let IcontypeObject = icontype
  if (IcontypeObject === "car") {
    realCarCount++;
  } else if (IcontypeObject === "bicycle") {
    realBicycleCount++;
  }
}


mapboxgl.accessToken = 'pk.eyJ1IjoibWFyaWEtaTI4MyIsImEiOiJjamVqMDVxdm4zYzl5MnBsbnZhZjV1MDkyIn0.bl4Y0AewatInkbnEfTs6Pg'


var metersPerPixel = function(latitude, zoomLevel) {
  var earthCircumference = 40075017;
  var latitudeRadians = latitude * (Math.PI/180);
  return earthCircumference * Math.cos(latitudeRadians) / Math.pow(2, zoomLevel + 8);
};

var pixelValue = function(latitude, meters, zoomLevel) {
  return meters / metersPerPixel(latitude, map.getZoom);
};

const metersToPixels = function (meters,latitude){
  return meters * metersPerPixel(latitude,map.getZoom())
}
//we need it to calculate the size of the grey range
const metersToPixelsAtMaxZoom = (meters, latitude) =>
  meters / 0.075 / Math.cos(latitude * Math.PI / 180)

var map = new mapboxgl.Map({
  container: 'map',
  style: 'mapbox://styles/mapbox/dark-v9',//'mapbox://styles/mapbox/satellite-streets-v10',
  center: [-121.076168,37.541611], 
  zoom: 5.6,
  maxZoom: 20
});
  
// console.log("here", map.getCenter());
distanceItem.addEventListener('change', function(e) {
  distance = document.getElementById('distance').value;
  pixels = metersToPixels(distance,lat)
  var radiusCircle = {
    stops: [
              [0, 0],
              [20, metersToPixelsAtMaxZoom(distance, lat)]
            ],
            base: 2
  }

  map.setPaintProperty("circleRange", 'circle-radius', radiusCircle);
});


var geolocate = new mapboxgl.GeolocateControl(({
  positionOptions: {
      enableHighAccuracy: true
  },
  trackUserLocation: false
}));

map.addControl(geolocate);


map.on('load', function() {
  document.getElementById('layerTypeObject').value = "noneObject";

  //add fake data
  map.addLayer({
    id: 'fakeCoord',
    type: 'circle',
    source: {
      type: 'geojson',
      data: fakeCoord
    },
    paint: {
      'circle-color': [       
        'match',
          ['get', 'Icontype'],
          'car',  [ 'case',
            ["<=", ['number', ['get', 'cabintemp']], 50], '#F87431',
            ["<=", ['number', ['get', 'cabintemp']], 70], '#E55451',
            ["<=", ['number', ['get', 'cabintemp']], 80], '#F778A1',
            '#F87431'
          ],
          /* other */ '#1589FF' //bicycle
      ],
      "circle-radius": 3
    }
  });

  var dbclick = false

  map.on('dblclick', function(e) {
    console.log('A dblclick event has occurred at ' + e.lngLat);
    lat = e.lngLat.lat
    lon = e.lngLat.lng
    dbclick = true
  });

  map.on('sourcedata', function(e) {
    if (e.isSourceLoaded && dbclick) {
      renderDistance()
      dbclick = false
    }
  });

  //Count car and bike on fake gps data
  fakeCoord.features.forEach(function(feature,rowIndex) {
    var IcontypeObject = feature.properties['Icontype'];     

    if(IcontypeObject == "car"){
      carCount++;
    }
    else if(IcontypeObject == "bicycle")
      bicycleCount++;
  });
  document.getElementById("totCars").innerHTML = "  " + carCount;
  document.getElementById("totBicycle").innerHTML = "  " + bicycleCount;

  // geolocate.trigger();
  geolocate.on('geolocate', function(e) {

  });
  
  //Add change event on selection of car or bicycle in dropdown menu
  layerTypeObject.addEventListener('change', function(e) {
    if(document.getElementById('layerTypeObject').value == "carObject")
      filterType = ['==', ['get', 'Icontype'], "car"];
    else if(document.getElementById('layerTypeObject').value == "bicycleObject")
      filterType = ['==', ['get', 'Icontype'], "bicycle"];
    else
      filterType = ['match', ['get', 'Icontype'], ["car","bicycle"], true, false];
    map.setFilter('fakeCoord', ['all', filterType]);
    if(map.getLayer('shareLocationsDot')){
      map.setFilter('shareLocationsDot', ['all', filterType]);
    }

  });


  var dataLayer;

  //Trigger the API every 10 seconds
  setInterval(function(){ 
    if(map.getLayer('shareLocationsDot')){
      $.ajaxSetup({
        async: false
      });
      //$.getJSON('http://localhost:8081/retrieve?search={"lat":'+ lat +',"lng":' + lon + ',"timespan":10,"distance":'+ distance +'}', function(data) {
      $.getJSON('https://locationserver.uswest2.development.volvo.care/retrieve?search={"lat":'+ lat+',"lng":' + lon + ',"timespan":10,"distance":'+ distance +'}', function(data) {
        shareLocations = data;
        console.log(data);

        //Remove UUID duplicates 
        // UUIDorigin = [];
        realCarCount = 0;
        realBicycleCount = 0;
        var newShareLocations = [];
        if(shareLocations == null){
          console.log("shareLocations data is null");
          map.setLayoutProperty('shareLocationsDot', 'visibility', 'none'); 
        }
        else{
          map.setLayoutProperty('shareLocationsDot', 'visibility', 'visible');

          shareLocations.forEach(function(element,rowIndex) {
            if (!uuidMemberTest(newShareLocations,element.properties['UUID'])) {
              (function(dataelement){
                newShareLocations.push(dataelement)
                countCarsAndBikes(element.properties['Icontype'])
              })(shareLocations[rowIndex]);
            }

          });



          console.log("New shareLocations array after removing the double uuid:");
          console.log(newShareLocations);
          dataLayer={
                  type: "FeatureCollection",
                  features: newShareLocations
                }
          
          //To update the map layer
          map.getSource('shareLocationsDot').setData(dataLayer);
        }

        document.getElementById("totCars").innerHTML = "  " + (carCount + realCarCount);
        document.getElementById("totBicycle").innerHTML = "  " + (bicycleCount + realBicycleCount);
      });

      $.ajaxSetup({
        async: true
      });
    }
  }, 3000);


});


// Create a popup, but don't add it to the map yet.
var popup = new mapboxgl.Popup({
    closeButton: false,
    closeOnClick: false
  });

map.on('mouseenter', "fakeCoord", function(e) {
  // Change the cursor style as a UI indicator.
  map.getCanvas().style.cursor = 'pointer';

  var coordinates = e.features[0].geometry.coordinates.slice();
  var description = "<strong>Climate: </strong>"+e.features[0].properties.cabintemp + "<p><strong>Coordinates: </strong>"+coordinates+"</p>";  

  // Ensure that if the map is zoomed out such that multiple
  // copies of the feature are visible, the popup appears
  // over the copy being pointed to.
  while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
    coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
  }

  // Populate the popup and set its coordinates
  // based on the feature found.
  popup.setLngLat(coordinates)
      .setHTML(description)
      .addTo(map);
});

map.on('mouseenter', "shareLocationsDot", function(e) {
  map.getCanvas().style.cursor = 'pointer';

  var coordinates = e.features[0].geometry.coordinates.slice();
  var description = "<strong>Climate: </strong>"+e.features[0].properties.cabintemp + "<p><strong>UUID: </strong>"+e.features[0].properties.UUID+"</p>";  

  while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
    coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
  }
  popup.setLngLat(coordinates)
      .setHTML(description)
      .addTo(map);
});

map.on('mouseleave', "fakeCoord", function() {
  map.getCanvas().style.cursor = '';
  popup.remove();
});   
map.on('mouseleave', "shareLocationsDot", function() {
  map.getCanvas().style.cursor = '';
  popup.remove();
});  

function renderDistance(){
  if (map.getSource('source_range_circle')){

    if (map.getLayer('circleRange')){
      map.removeLayer('circleRange')
    }
    map.removeSource('source_range_circle')
  }

  if (map.getSource('source_range_circle')) {
    map.getSource('source_range_circle').setData({
      "type": "geojson",
      "data": {
        "type": "FeatureCollection",
        "features": [{
          "type": "Feature",
          "geometry": {
            "type": "Point",
            "coordinates": [lon, lat],
          }
        }]
      }
    });
  } else{
      map.addSource("source_range_circle", {
      "type": "geojson",
      "data": {
       "type": "FeatureCollection",
        "features": [{
          "type": "Feature",
         "geometry": {
            "type": "Point",
            "coordinates": [lon, lat],
          }
        }]
      }
    });
  }

  circleRange = map.getLayer('circleRange')
  if (!(circleRange === undefined)) {

    var radiusCircle = {
      stops: [
        [0, 0],
        [20, metersToPixelsAtMaxZoom(distance, lat)]
      ],
      base: 2
    }
    map.setPaintProperty("circleRange", 'circle-radius', radiusCircle);

  }else {
    map.addLayer({
      "id": "circleRange",
      "type": "circle",
      "source": "source_range_circle",
      "paint": {
        "circle-radius": {
          stops: [
            [0, 0],
            [20, metersToPixelsAtMaxZoom(distance, lat)]
          ],
          base: 2
        },
        "circle-color": "white",
        "circle-opacity": 0.3
      }
    });
  }
    //END to draw the gray range circle


    if (map.getSource("shareLocationsDot")){

    }else{
      map.addLayer({
        id: 'shareLocationsDot',
        type: 'circle',
        source: {
          type: 'geojson',
          data: {
            type: "FeatureCollection",
            features: shareLocations
          }
        },
        filter: ['all', filterType],
        paint: {
          'circle-color': [
            'match',
            ['get', 'Icontype'],
            'car', ['case',
              ["<=", ['number', ['get', 'cabintemp']], 50], '#F87431',
              ["<=", ['number', ['get', 'cabintemp']], 70], '#E55451',
              ["<=", ['number', ['get', 'cabintemp']], 80], '#F778A1',
              '#F87431'
            ],
            /* other */ '#1589FF' //bicycle
          ],
          "circle-radius": 3
        }
      });
    }
}
