package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

// ReqBody is type for input request
type ReqBody struct {
	Input string `json:"input"`
}

// Resp is type to provide the req resp
type Resp struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type WordstringType []Resp

func (w WordstringType) Len() int           { return len(w) }
func (w WordstringType) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w WordstringType) Less(i, j int) bool { return w[i].Count > w[j].Count }

// RegisterRoute register route and created server
func RegisterRoute() {
	r := mux.NewRouter()
	r.HandleFunc("/count", handlerCount).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

// handlerCount is handler to handle the request
func handlerCount(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var rBody ReqBody
	err = json.Unmarshal(data, &rBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if rBody.Input != "" {
		output := countWordsAndOccurnace(rBody.Input)
		_ = json.NewEncoder(w).Encode(output)
	}
}

// countWordsAndOccurnace helper function to calculate the count and its occurances
func countWordsAndOccurnace(input string) []Resp {
	splitInputArray := strings.Split(input, " ")
	mapForKeyAndConut := make(map[string]int)
	if len(splitInputArray) > 0 {
		for _, item := range splitInputArray {
			mapForKeyAndConut[item] += 1
		}
	}
	sortArray := make(WordstringType, len(mapForKeyAndConut))
	i := 0
	for key, value := range mapForKeyAndConut {
		sortArray[i] = Resp{key, value}
		i++
	}
	sort.Sort(sortArray)
	arrayLen := len(sortArray)
	if arrayLen > 10 {
		return sortArray[0:10]
	}
	return sortArray[0:arrayLen]
}
