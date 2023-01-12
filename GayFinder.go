package main

import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
    "net/http"
    "encoding/json"
    "github.com/googlemaps/google-maps-services-go/maps"
)

type Person struct {
    ID int
    Name string
    Profile string
    Gender string
    SexualOrientation string
    Services map[string]interface{}
    Position maps.LatLng
}

func main() {
    db, err := sql.Open("postgres", "user=username password=password dbname=dbname sslmode=disable")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    people, err := getPeopleWithConsent(db)
    if err != nil {
        panic(err.Error())
    }

    for _, person := range people {
        fmt.Println(person.Name, person.Profile, person.Gender, person.SexualOrientation)
        for service, data := range person.Services {
            fmt.Printf("%s: %v\n", service, data)
        }
        fmt.Println("Position: ", person.Position)
    }
    friendsMap := findConnections(people)
    for person, friends := range friendsMap {
        fmt.Printf("%s is friends with %v\n", person.Name, friends)
    }

    func getPeopleWithConsent(db *sql.DB) ([]Person, error) {
    rows, err := db.Query("SELECT id, name, profile, gender, sexual_orientation, services, address FROM people WHERE has_consent = 1")
    if err != nil {
    return nil, err
    }
    defer rows.Close()
    people := []Person{}
for rows.Next() {
    var person Person
    var servicesJSON string
    var address string
    err := rows.Scan(&person.ID, &person.Name, &person.Profile, &person.Gender, &person.SexualOrientation, &servicesJSON, &address)
    if err != nil {
        return nil, err
    }
    json.Unmarshal([]byte(servicesJSON), &person.Services)
    person.Position, err = getLocationFromAddress(address)
    if err != nil {
        return nil, err
    }
    people = append(people, person)
}
return people, nil
}

func fetchDataFromService(serviceName string, username string) (map[string]interface{}, error) {
resp, err := http.Get("https://vk.com/" + serviceName + "/" + username)
if err != nil {
return nil, err
}
defer resp.Body.Close()
var data map[string]interface{}
json.NewDecoder(resp.Body).Decode(&data)

return data, nil
}

func findConnections(people []Person) map[Person][]Person {
friendsMap := make(map[Person][]Person)
for i, person := range people {
friendsList := []Person{}
for j, friend := range people {
if i == j {
continue
}
if haveCommonInterest(person, friend) {
friendsList = append(friendsList, friend)
}
}
friendsMap[person] = friendsList
}
return friendsMap
}

func haveCommonInterest(p1, p2 Person) bool {
// Write your own code here to check if two people have common interests
// based on the information in their profiles
// Example:
// return p1.Profile == p2.Profile
//
// Return true if they have common interests, otherwise return false.
return true
}

func getLocationFromAddress(address string) (maps.LatLng, error) {
c, err := maps.NewClient(maps.WithAPIKey("YOUR_API_KEY"))
if err != nil {
return maps.LatLng{}, err
}
r := &maps.GeocodingRequest{
    Address: address,
}

resp, err := c.Geocode(context.Background(), r)
if err != nil {
    return maps.LatLng{}, err
}

return resp.Results[0].Geometry.Location, nil
}


