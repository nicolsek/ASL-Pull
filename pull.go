/* The goal of this is going to be pulling a simple webpage off of a website that handles ASL definitions, possible at sompoint implementing multiple manners of managing that. */

package main

import (
	"errors"
	"io"
	"net/http"
	"os"
)

// main ... The main function?
func main() {
	//Setting up the default layout.
	address := "http://www.lifeprint.com/asl101/pages-signs/" //Example URL for apple --> http://www.lifeprint.com/asl101/pages-signs/a/apple.htm
	//Getting the url from the getUrl function.
	url, word, err := getURL()
	//Checking the error if exists panic.
	check(err)
	//Parse the url.
	resp, err := http.Get(address + url)
	//Checks the error.
	check(err)
	//Creating the file and writing the data from the HTTP request to it.
	writeFile(resp, word)

}

// writeFile ... Creates the file with the name of word and copies the body of the HTTP request to it.
func writeFile(resp *http.Response, word string) {
	//Creates the file with the name word.
	file, err := os.Create(word)
	//Checks the error.
	check(err)
	//Writes to the file by copying the data.
	io.Copy(file, resp.Body)
	//Closes the file when done.
	defer file.Close()
}

// getURL ... Gets the url by parsing the terminal input using Os.Args[], returns the url.
func getURL() (string, string, error) {
	//Creating the error to be returned.
	err := errors.New("Unable to retrieve URL")
	//Getting the URL by accessing an element of an array from os.Args[] @ position 1.
	word := os.Args[1]

	//Checking if os.Args[1] exists.
	if os.Args[1] != "" {
		err = nil
	}

	//Format the url --> example URL for apple --> http://www.lifeprint.com/asl101/pages-signs/a/apple.htm
	//What we want /a/apple.htm

	/*

	   1. Get the first letter --> Using string slicing.
	   2. Take the entire word --> Already gotten from array slicing the Os.Args[]
	   3. Add .htm

	*/

	//Creating an array to hold all the values which at the end we will pull together to create the final string.
	urlVars := make([]string, 3)

	//Getting the first letter
	urlVars[0] = word[0:1] + "/" //Accessing the first element of the string.
	//We already know the word so just set that directly.
	urlVars[1] = word
	//Just set [2] to .htm
	urlVars[2] = ".htm"

	//Creating the final string for the URL
	url := ""

	//Set the original URL value to what we want.
	//Creats a for each loop and itterates through the entire thing.
	for _, s := range urlVars {
		url += s
	}

	//Returns the values.
	return url, word, err
}

// check ... Checks if there is an error and if there is panic the error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
