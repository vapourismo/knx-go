package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
)

type Range struct {
	Low  float32 `json:"Low"`
	High float32 `json:"High"`
}

type DPTDescription struct {
	Id       string            `json:"ID"`
	Name     string            `json:"Name"`
	Type     string            `json:"Type"`
	Range    Range             `json:"Range"`
	Unit     string            `json:"Unit"`
	Format   string            `json:"Format"`
	ValueMap map[string]string `json:"ValueMap"`
}

func replace(input string, from string, to string) string {
	return strings.Replace(input, from, to, -1)
}

func isset(f interface{}) bool {
	if !reflect.ValueOf(f).IsValid() {
		return false
	}

	if reflect.ValueOf(f).Kind() == reflect.Struct {
		if reflect.DeepEqual(f, reflect.Zero(reflect.TypeOf(f)).Interface()) {
			return false
		}
	}

	return true
}

func iseqf(a float32, b int) bool {
	return a == float32(b)
}

func main() {
	tmpl_path := flag.String("template", "", "path to the template file")
	flag.Parse()

	function_map := template.FuncMap{
		"replace": replace,
		"isset":   isset,
		"iseqf":   iseqf,
		"title":   strings.Title,
	}

	// Parse the template file
	t := template.Must(template.New(*tmpl_path).Funcs(function_map).ParseFiles(*tmpl_path))

	// Read and parse the json file
	// m := map[string]interface{}{}
	m := make([]DPTDescription, 0)

	json_file, err := os.Open("DataPointTypes.json")
	if err != nil {
		fmt.Println(err)
	}
	json_data, _ := ioutil.ReadAll(json_file)
	json_file.Close()

	if err := json.Unmarshal([]byte(json_data), &m); err != nil {
		panic(err)
	}

	if err := t.Execute(os.Stdout, m); err != nil {
		panic(err)
	}
}
