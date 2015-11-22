package main 
import (
"fmt"
"math/rand"
"strings"
"net/http"
"github.com/julienschmidt/httprouter"
"encoding/json"
"io/ioutil"
"os"
"log"
"strconv"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"github.com/anweiss/uber-api-golang/uber"
 "sort"
 "reflect"
"bytes"
)
type (UserController struct {
	session *mgo.Session
	})
type Request struct {
	Name string `json:"name"`
	Address string `json:"address"`
	City string `json:"city" `
	State string `json:"state"`
	Zip string `json:"zip"`
}
	
type Response struct {
	Address    string `json:"address" bson:"address"`
	City       string `json:"city" bson:"city"`
	Coordinate struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
	} `json:"coordinate" bson:"coordinate"`
	ID    int   `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	State string `json:"state" bson:"state"`
	Zip   string `json:"zip" bson:"zip"`
}
//JSON struct from GooGle Maps Api
type GoogleMaps struct {
	Results []struct {
		AddressComponents []struct {
			LongName string `json:"long_name"`
			ShortName string `json:"short_name"`
			Types []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string `json:"place_id"`
		Types []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

//Function for HTTP POST
func  PostRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder:= json.NewDecoder(r.Body)

	var u Request
	err:= decoder.Decode(&u)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var c Response
    c.Name= u.Name
    c.Address= u.Address
    c.City= u.City
    c.State= u.State
    c.Zip= u.Zip 
    c.ID=rand.Intn(10000000)
    var fulladdress string 
    fulladdress= c.Address+" "+c.City
    latresponse := GetLatitude(fulladdress)
    longresponse := GetLongitude(fulladdress)
    c.Coordinate.Lat= latresponse
    c.Coordinate.Lng= longresponse
    sess:=getSession();
    collection:= sess.DB("trip-planner").C("locations")
    e:= collection.Insert(c)
    if e!=nil {
    	panic(e)
    }

	uj,_ := json.Marshal(c)
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}
//Function to obtain Latitude coordinates for a user's location
func GetLatitude(fulladdress string) float64{
	var Ad GoogleMaps
	var lat64 float64
	Baseur:= "http://maps.google.com/maps/api/geocode/json?address="
	Addur:= fulladdress
	Urlf:= Baseur + Addur
	Urlf = strings.Replace(Urlf," ","%20",-1)
    	apiRes, err:= http.Get(Urlf)
	if err!=nil {
		fmt.Printf("error occurred")
		fmt.Printf("%s", err)
		os.Exit(1)
	}else{
		defer apiRes.Body.Close()
		contents, err:= ioutil.ReadAll(apiRes.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
	   err= json.Unmarshal(contents, &Ad)
	   if err!=nil {
	   	fmt.Println("her is the error")
	   	fmt.Printf("%s", err)
	   	os.Exit(1)
	   }
	 lat64 = Ad.Results[0].Geometry.Location.Lat
	}
	 return lat64
}

//Function to access Longitude coordinates for a user's location
func GetLongitude(fulladdress string) float64{
	var s GoogleMaps
	var long64 float64
	Baseurl:= "http://maps.google.com/maps/api/geocode/json?address="
	Addurl:= fulladdress
	Url:= Baseurl + Addurl
	Url = strings.Replace(Url," ","%20",-1)
	apiResponse, err:= http.Get(Url)
	if err!=nil {
		fmt.Printf("error occurred")
		fmt.Printf("%s", err)
		os.Exit(1)
	}else{
		defer apiResponse.Body.Close()
		contents, err:= ioutil.ReadAll(apiResponse.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		err= json.Unmarshal(contents, &s)
		if err!=nil {
			fmt.Println("Here is the error from longitude")
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		long64=s.Results[0].Geometry.Location.Lng		
}
return long64
}

//Function to obtain HTTP GET Request
func  GetRequest(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	 var updatedmsg Response
	 a:= p.ByName("id")
	 ac,_ := strconv.Atoi(a)
	   sess:=getSession();
  er := sess.DB("trip-planner").C("locations").Find(bson.M{"id": ac}).One(&updatedmsg)
  if er!=nil {
  	panic(er)
  }
	uj,_ := json.Marshal(updatedmsg)
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

//Function to obtain HTTP DELETE Request
func  DeleteRequest(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	 a:= p.ByName("id")
	 ac,_ := strconv.Atoi(a)
	   sess:=getSession();
  er := sess.DB("trip-planner").C("locations").Remove(bson.M{"id": ac})
  if er!=nil {
  	panic(er)
  }
	w.WriteHeader(200)
}

//Function to access HTTP PUT Request
func  PutRequest(w http.ResponseWriter, res *http.Request, p httprouter.Params){
	decoder:= json.NewDecoder(res.Body)

	var r Request
	err:= decoder.Decode(&r)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var oldupdatingmsg Response
	var updatingmsg Response
	init:= p.ByName("id")
	abs, _:= strconv.Atoi(init)
	newsession:= getSession();
	errors:= newsession.DB("trip-planner").C("locations").Find(bson.M{"id": abs}).One(&oldupdatingmsg)
	if errors!=nil {
		panic(errors)
	}
		updatingmsg.Name= r.Name
	    updatingmsg.Address= r.Address
	    updatingmsg.City= r.City
	    updatingmsg.Zip= r.Zip
	    updatingmsg.State= r.State
        updatingmsg.ID= abs
	var updateaddress string
	updateaddress= updatingmsg.Address+updatingmsg.City
	updatelatresp:= GetLatitude(updateaddress)
	updatelongresp:= GetLongitude(updateaddress)
	updatingmsg.Coordinate.Lat= updatelatresp
	updatingmsg.Coordinate.Lng= updatelongresp
	collec:= newsession.DB("trip-planner").C("locations")
    ef:= collec.Update(oldupdatingmsg,updatingmsg)
    if ef!=nil {
    	panic(ef)
    }
    updatejson,_ := json.Marshal(updatingmsg)
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
	fmt.Fprintf(w, "%s", updatejson)
}


type UserPostRequest struct {
	StartingPositionid string`json:"starting_from_location_id"`
	OtherPositionids []string `json:"location_ids"`
}
type UserPostResponse struct {
	Id string `json:"id"`
	Status string `json:"status"`
	StartingPositionid string `json:"starting_from_location_id"`
	BestrouteIds []string `json:"best_route_location_ids"`
	Totalcost int `json:"total_uber_cost"`
	Totalduration int `json:"total_uber_duration"`
	Distance float64 `json:"total_distance"`
}
type RideReq struct {
    EndLatitude    string `json:"end_latitude"`
    EndLongitude   string `json:"end_longitude"`
    ProductID      string `json:"product_id"`
    StartLatitude  string `json:"start_latitude"`
    StartLongitude string `json:"start_longitude"`
}

type UserPutResponse struct {
  ID string `json:"id"`
  Status string `json:"status"`
  StartingFromLocationID string `json:"starting_from_location_id"`
  NextDestinationLocationID string `json:"next_destination_location_id"`
  BestRouteLocationIds []string `json:"best_route_location_ids"`
  TotalUberCost int `json:"total_uber_cost"`
  TotalUberDuration int `json:"total_uber_duration"`
  TotalDistance float64 `json:"total_distance"`
  UberWaitTimeEta int `json:"uber_wait_time_eta"`
}

type Responsefromoauth struct {
    Driver          interface{} `json:"driver"`
    Eta             int         `json:"eta"`
    Location        interface{} `json:"location"`
    RequestID       string      `json:"request_id"`
    Status          string      `json:"status"`
    SurgeMultiplier int         `json:"surge_multiplier"`
    Vehicle         interface{} `json:"vehicle"`
}

type ResponsefromDB struct {
	Address    string `json:"address" bson:"address"`
	City       string `json:"city" bson:"city"`
	Coordinate struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
	} `json:"coordinate" bson:"coordinate"`
	ID    int   `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	State string `json:"state" bson:"state"`
	Zip   string `json:"zip" bson:"zip"`
}
type Products struct {
  Latitude  float64
  Longitude float64
  Products  []Product `json:"products"`
}

// Uber product
type Product struct {
  ProductId   string `json:"product_id"`
  Description string `json:"description"`
  DisplayName string `json:"display_name"`
  Capacity    int    `json:"capacity"`
  Image       string `json:"image"`
}





var id int
var tripID int


func PlanTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var u UserPostRequest
	err:= decoder.Decode(&u)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(u.StartingPositionid)

	var options uber.RequestOptions
	options.ServerToken= "uiiQ8zd9D4GHSDdsmT_mzIw7DXS67enDl5tXuc-p"
    options.ClientId= "c6jVXK_x-UXgKRsNUdyDb2omyFjVyKPy"
    options.ClientSecret= "vVksauQ2gWAVFxlYDi4GDi65J_f-zVxmVFqibMyA"
    options.AppName= "Golang Application"
    options.BaseUrl= "https://sandbox-api.uber.com/v1/"

    client := uber.Create(&options)
    startingidstring, error := strconv.Atoi(u.StartingPositionid)
    if error != nil {
         panic(error)
     }
     var updatedmsg Response
     sess:= getSession()
	er := sess.DB("trip-planner").C("locations").Find(bson.M{"id": startingidstring}).One(&updatedmsg)
    if er!=nil {
  	panic(er)
  }

  index:= 0
  totalprice := 0
  totaldistance := 0.0
  totalduration := 0
  bestroute := make([]int, len(u.OtherPositionids))
  m := make(map[int]string)
  for _, ids := range u.OtherPositionids{
  	otherids, err1 := strconv.Atoi(ids)
  	 if err1 != nil {
            panic(err1)
        }
        var otherlocid ResponsefromDB
        session:= getSession()
	er = session.DB("trip-planner").C("locations").Find(bson.M{"id": otherids}).One(&otherlocid)
    if er!=nil {
  	panic(er)
  }  
       x:= &uber.PriceEstimates{}
        x.StartLatitude = updatedmsg.Coordinate.Lat;
        x.StartLongitude = updatedmsg.Coordinate.Lng;
        x.EndLatitude = otherlocid.Coordinate.Lat;
        x.EndLongitude = otherlocid.Coordinate.Lng;
        if e := client.Get(x); e != nil {
            fmt.Println(e);
        }
        totaldistance=totaldistance+x.Prices[0].Distance;
        totalduration=totalduration+x.Prices[0].Duration;
        totalprice=totalprice+x.Prices[0].LowEstimate;
        bestroute[index]=x.Prices[0].LowEstimate;
        m[x.Prices[0].LowEstimate]=ids;
        index=index+1;
  }
  sort.Ints(bestroute)
  var postuber UserPostResponse
  tripID= tripID+1
  postuber.Id=strconv.Itoa(tripID)
     postuber.Distance=totaldistance
     postuber.Totalcost=totalprice
     postuber.Totalduration=totalduration
     postuber.Status="Planning"
     postuber.StartingPositionid=strconv.Itoa(startingidstring)
     postuber.BestrouteIds=make([]string,len(u.OtherPositionids))
     index=0;
     for _, ind := range bestroute{
        postuber.BestrouteIds[index]=m[ind];
        index=index+1;
     }
     sessions:=getSession();
    collection:= sessions.DB("trip-planner").C("trip")
    e:= collection.Insert(postuber)
    if e!=nil {
    	panic(e)
    }
     js,err := json.Marshal(postuber)
    if err != nil{
       fmt.Println("Error")
       return
    }
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(201)
    fmt.Fprintf(w, "%s", js)
}

func GetTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var getquery UserPostResponse
	ac:= p.ByName("id")
	
	sessionss:=getSession();
  er := sessionss.DB("trip-planner").C("trip").Find(bson.M{"id": ac}).One(&getquery)
  if er!=nil {
  	panic(er)
  }
	uj,_ := json.Marshal(getquery)
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

var currentPos int
var otherint int 


func PutTrip(w http.ResponseWriter, r*http.Request, p httprouter.Params) {
  //ac:= p.ByName("id")
  ac:= p[0].Value
  var dbresult UserPostResponse
  var finalresponse UserPutResponse
  var r1  ResponsefromDB
  var r2  ResponsefromDB
  sessionss:=getSession();
  er := sessionss.DB("trip-planner").C("trip").Find(bson.M{"id": ac}).One(&dbresult)
  if er!=nil {
    panic(er)
  }else{
  
  var ids int 
  if currentPos ==0 {
    ids, err := strconv.Atoi(dbresult.StartingPositionid)
    if err != nil {
        // handle error
            fmt.Println(err)
        }
        fmt.Println(ids)
        fmt.Println(reflect.TypeOf(ids))
        
   }else {
       ids, err := strconv.Atoi(dbresult.BestrouteIds[currentPos-1])
   
        if err != nil {
        // handle error
            fmt.Println(err)
        }
        fmt.Println(ids)
      
   }

  err := sessionss.DB("trip-planner").C("locations").Find(bson.M{"id": ids}).One(&r1)
  if err!=nil {
    fmt.Println("Here is the error2")
    log.Println(err.Error())
    panic(err)
  }
  
  ids, err = strconv.Atoi(dbresult.BestrouteIds[currentPos])
    if err != nil {
        // handle error
        fmt.Println(err)
    }
     newsesss:=getSession();
  err = newsesss.DB("trip-planner").C("locations").Find(bson.M{"id": ids}).One(&r2)
  if err!=nil {
    panic(err)
  }
}
finalresponse.ID = strconv.Itoa(otherint)
finalresponse.BestRouteLocationIds= dbresult.BestrouteIds
finalresponse.StartingFromLocationID = dbresult.StartingPositionid
finalresponse.TotalUberCost= dbresult.Totalcost
finalresponse.TotalUberDuration= dbresult.Totalduration

if finalresponse.Status != "completed" {
  if finalresponse.Status== "planning" && len(finalresponse.NextDestinationLocationID)==0 {
    finalresponse.Status="requesting"
    finalresponse.NextDestinationLocationID= finalresponse.BestRouteLocationIds[0]
  }else if finalresponse.StartingFromLocationID== finalresponse.NextDestinationLocationID{
    finalresponse.Status= "completed"
  }else{
    finalresponse.Status="requesting"
    finalresponse.NextDestinationLocationID= dbresult.BestrouteIds[currentPos]
  }
}

  var options uber.RequestOptions
  options.ServerToken= "uiiQ8zd9D4GHSDdsmT_mzIw7DXS67enDl5tXuc-p"
    options.ClientId= "c6jVXK_x-UXgKRsNUdyDb2omyFjVyKPy"
    options.ClientSecret= "vVksauQ2gWAVFxlYDi4GDi65J_f-zVxmVFqibMyA"
    options.AppName= "Golang Application"
    options.BaseUrl= "https://sandbox-api.uber.com/v1/"
     client := uber.Create(&options)

     pl := &uber.Products{}
     pl.Latitude=r1.Coordinate.Lat
     pl.Longitude=r1.Coordinate.Lng
     if e := client.Get(pl); e != nil {
         fmt.Println(e)
    }
    var productid string;
    i:=0
    for _, product := range pl.Products {
         if(i == 0){
             productid = product.ProductId
        }
    }

    var ride RideReq
    ride.StartLatitude=strconv.FormatFloat(r1.Coordinate.Lat, 'f', 6, 64)
    ride.StartLongitude=strconv.FormatFloat(r1.Coordinate.Lng, 'f', 6, 64)
    ride.EndLatitude=strconv.FormatFloat(r2.Coordinate.Lat, 'f', 6, 64)
    ride.EndLongitude=strconv.FormatFloat(r2.Coordinate.Lng, 'f', 6, 64)
    ride.ProductID=productid
    ridejson, _ := json.Marshal(ride)
    url:= "https://sandbox-api.uber.com/v1/requests"
    requ, err:= http.NewRequest("POST", url, bytes.NewBuffer(ridejson))
    if err!=nil {
      panic(err)
    }
    requ.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3QiLCJyZXF1ZXN0X3JlY2VpcHQiLCJoaXN0b3J5X2xpdGUiXSwic3ViIjoiZDdjY2JjYmMtZmU0MS00NTYwLThiNTQtN2QxY2E0Y2MzOGZkIiwiaXNzIjoidWJlci11czEiLCJqdGkiOiIwYjk3MThlMi0yOGViLTQ4YjctODc4MS0xM2ZiYWI2NDU5ZGUiLCJleHAiOjE0NTA3NjUwNDIsImlhdCI6MTQ0ODE3MzA0MiwidWFjdCI6ImdjcHJ1SHBlcmsxcm9ObWNJcUc5MDh5OTloYW1aZiIsIm5iZiI6MTQ0ODE3Mjk1MiwiYXVkIjoiYzZqVlhLX3gtVVhnS1JzTlVkeURiMm9teUZqVnlLUHkifQ.iu0wdQ-xEhu3uxNkFFz1N2vEVA1f-bil7_UzSMdkP8mBh-Ao4r3lGRrOCGiYg9fZuHjxNI8_ijpa8V-iJtD3lfkLZ1y-9uNSNLUcbJBGYQyNo3ucJ3BqIXVx6XIEw7AYP72kz9n5kjIEszYNG25yIF4-I2cYtgWxZrt_nv1_dCkdMVtRQ7WrAZFOkjv45tPH1yLYyuH8x56kZ5wEvrHkZKaLgzsgI9-UOBpy9dhaYuHDVqqbltz94N8wdw_RMGaUc4xKAsX2nSAMrnh-QGHl3L9c0ZzjAtgGPZ1RgweDdfrSc4vPUbI72HQLarwuVyhIcJ5LE7ArT0mHM_ZtS7ZDrA")
    requ.Header.Set("Content-Type", "application/json")
    outdata, err:= ioutil.ReadAll(requ.Body)
    var abc Responsefromoauth
    err= json.Unmarshal(outdata, &abc)
    finalresponse.UberWaitTimeEta= abc.Eta
    updatejson,err := json.Marshal(finalresponse)
    if err != nil{
       fmt.Println("Error")
       return
    }
     otherint=otherint+1;
    currentPos=currentPos+1;
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(201)
    fmt.Fprintf(w, "%s", updatejson)

}




func main() {
	uberrouter:= httprouter.New()
	uberrouter.GET("/:id", GetRequest)
	uberrouter.POST("/", PostRequest)
	uberrouter.DELETE("/:id", DeleteRequest)
	uberrouter.PUT("/:id", PutRequest)

	uberrouter.POST("/trip", PlanTrip) //For First Requirement
	uberrouter.GET("/:id/trip", GetTrip)//For Second Requirement
	uberrouter.PUT("/:id/request", PutTrip)//for Third Requirement
	log.Fatal(http.ListenAndServe(":8011", uberrouter))
}

//Function to connect to MongoDB
func getSession() *mgo.Session {
	connect, err:= mgo.Dial("mongodb://suyash:123@ds031531.mongolab.com:31531/trip-planner")
	if err!= nil {
			panic(err)
		}	
		return connect
}