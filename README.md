# Uber-Trip-Planning-
Trip planning api that tells its users about the best possible route available(based on cost as per uber api calculations) among the places the user is interested in.
The api has three services:-
Plan a Trip - Plan a Trip based on the interested location IDs
Request- POST        /trips 
{
    "starting_from_location_id: "999999",
    "location_ids" : [ "10000", "10001", "20004", "30003" ] 
}
Response- 
{
     "id" : "1122",
     “status” : “planning”,
     "starting_from_location_id: "999999",
     "best_route_location_ids" : [ "30003", "10001", "10000", "20004" ],
  "total_uber_costs" : 125,
  "total_uber_duration" : 640,
  "total_distance" : 25.05 
}

Get A Trip Information - Getting trip details based on ID
Request:  GET    /trips/1122
Response- {
     "id" : "1122",
     "status" : "planning",
     "starting_from_location_id: "999999",
     "best_route_location_ids" : [ "30003", "10001", "10000", "20004" ],
  "total_uber_costs" : 125,
  "total_uber_duration" : 640,
  "total_distance" : 25.05 
}

Start the trip by requesting UBER for the first destination. 
You will call UBER request API to request a car from starting point to the next destination.
Request:  PUT             /trips/1122/request
Response- 
{
     "id" : "1122",
     "status" : "requesting",
     "starting_from_location_id”: "999999",
     "next_destination_location_id”: "30003",
     "best_route_location_ids" : [ "30003", "10001", "10000", "20004" ],
  "total_uber_costs" : 125,
  "total_uber_duration" : 640,
  "total_distance" : 25.05,
  "uber_wait_time_eta" : 5 
}



