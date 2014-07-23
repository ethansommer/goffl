package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var ( 
		qbs []byte
	)
	qbs=getURL("http://www.fantasypros.com/nfl/projections/qb.php?export=xls")
	fmt.Printf("%s",qbs)
	
}


func getURL(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	retbs, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return retbs
}