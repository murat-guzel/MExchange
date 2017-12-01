package main

import(
 
	"time"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

)
var Msession *mgo.Session
 
func main() {

	Msession = BuildMongo() 
	HandleRequests()
	
	
}

type Task struct {

	Id bson.ObjectId `json:"_id" bson:"_id"`
	Name string `json:"Name"`
	Status string `json:"Status"`
	Date time.Time `json:"Date"`

}

func CreateTask(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Println(Msession)
	c := Msession.DB("Exchange").C("task") 
	// Index
	index := mgo.Index{
		Key:        []string{"id", "name","status","date"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	err = c.Insert(&Task{Id :bson.NewObjectId() , Name: name,Status : "InComplete", Date : time.Now()})
	
		if err != nil {
			panic(err)
		}else{
	
			fmt.Println("Insert Process Done")
		}

		json.NewEncoder(w).Encode("Task is created")

}

func UpdateTask(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	id := vars["id"]
	name := vars["name"]
	c := Msession.DB("Exchange").C("task") 
	 
	if id == ""{
		fmt.Println("Heeey id is null baby")
	}else{

		_id := bson.ObjectIdHex(id)
		fmt.Println(_id)
	   
		colQuerier := bson.M{"_id": _id}
		change := bson.M{"$set": bson.M{"name": name}}
		err := c.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}


	}

	fmt.Println("Task is changed")
	json.NewEncoder(w).Encode("Task is changed")
}

func DeleteTask(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	id := vars["id"]


	c := Msession.DB("Exchange").C("task") 
	var tasks []Task
	 
	if id != ""{
		
		_id := bson.ObjectIdHex(id)
		fmt.Println(_id)
		err := c.FindId(_id).All(&tasks)
		if err != nil {
		 
			panic(err)
		 
		}
		err = c.Remove(tasks[0])
		if err != nil {
			
			   panic(err)
			
		   }
		json.NewEncoder(w).Encode("Removed")
	}else{

		fmt.Println("Heeey id is null baby")
		json.NewEncoder(w).Encode("Heeey id is null baby")
	}

}

func ChangeStatusTask(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	id := vars["id"]
	c := Msession.DB("Exchange").C("task") 
	var tasks []Task
	if id == ""{
		fmt.Println("Heeey id is null baby")
	}else{

		_id := bson.ObjectIdHex(id)
		fmt.Println(_id)
	    c.FindId(_id).All(&tasks)
		statusVal := tasks[0].Status
		if	statusVal=="Done"{
			statusVal = "InComplete"
		}else{

			statusVal = "Done"
		}
		colQuerier := bson.M{"_id": _id}
		change := bson.M{"$set": bson.M{"status": statusVal}}
		err := c.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}


	}

	fmt.Println("Status is changed")
	json.NewEncoder(w).Encode("Status is changed")

	

}

func GetAllTasks(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	c := Msession.DB("Exchange").C("task") 
	var tasks []Task
	var err error
	if id != ""{
		
		_id := bson.ObjectIdHex(id)
		fmt.Println(_id)
		err = c.FindId(_id).All(&tasks)

	}else{

		err = c.Find(nil).All(&tasks)
	}
	

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", tasks)
	json.NewEncoder(w).Encode(tasks)
	
}

func BuildMongo() *mgo.Session{

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
 
	// Collection People
	c := session.DB("Exchange").C("task")

	// Index
	index := mgo.Index{
		Key:        []string{"id", "name","status","date"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	fmt.Println("Build Mongo Db Done!")
	
	return session.Clone()

}

func Test(w http.ResponseWriter, r *http.Request){

	fmt.Fprintf(w, "Success")
}


func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/test", Test)
	myRouter.HandleFunc("/ct/{name}", CreateTask)
	myRouter.HandleFunc("/dt/{id}", DeleteTask)
	myRouter.HandleFunc("/ut/{id}/{name}", UpdateTask)
	myRouter.HandleFunc("/cst/{id}", ChangeStatusTask) 
	myRouter.HandleFunc("/g/{id}", GetAllTasks) 
	myRouter.HandleFunc("/g", GetAllTasks) 

	fmt.Println("Build HandleRequests Done!")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}


