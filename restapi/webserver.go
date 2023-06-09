package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type student struct {
	ID    string
	Name  string
	Grade int
}

var data = []student{
	{
		ID:    "B001",
		Name:  "Testie",
		Grade: 1,
	},
	{
		ID:    "A002",
		Name:  "Kurisu",
		Grade: 3,
	},
	{
		ID:    "A003",
		Name:  "Kyoma",
		Grade: 2,
	},
}

const BASEURL = "http://localhost:8080"

func fetchUsers(client http.Client) ([]student, error) {
	var data []student
	var err error

	request, err := http.NewRequest("GET", BASEURL+"/users", nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func fetchUser(ID string, client http.Client) (student, error) {
	var data student
	var err error

	// for POST method
	// var param = url.Values{}
	// param.Set("id", ID)
	// var payload = bytes.NewBufferString(param.Encode())

	// request, err := http.NewRequest("POST", BASEURL+"/user", payload)
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	request, err := http.NewRequest("GET", BASEURL+"/user", nil)
	q := request.URL.Query()
	q.Add("id", ID)
	request.URL.RawQuery = q.Encode()

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

func main() {
	var err error
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}

	go RestFulAPIServer(ln)

	client := &http.Client{}

	users, err := fetchUsers(*client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, each := range users {
		fmt.Printf("ID: %s\t Name: %s\t Grade: %d\n", each.ID, each.Name, each.Grade)
	}

	user, err := fetchUser("B001", *client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("fetch spesific user")
	fmt.Printf("ID: %s\t Name: %s\t Grade: %d\n", user.ID, user.Name, user.Grade)

}

func RestFulAPIServer(ln net.Listener) {
	http.HandleFunc("/user", user)
	http.HandleFunc("/users", users)

	fmt.Println("Server running on http://localhost:8080")
	http.Serve(ln, nil)
}

func user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		id := r.FormValue("id")
		for _, user := range data {
			if user.ID == id {
				var encoded, err = json.Marshal(user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Write(encoded)
				return
			}
		}

		http.Error(w, "User Not Found", http.StatusBadRequest)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}
func users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		encoded, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(encoded)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
