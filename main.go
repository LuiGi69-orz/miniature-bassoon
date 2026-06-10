package main

import (
	"flag" 
	// "fmt"
)

//1) scan

//2) stats

func main(){
	var folder , email string
	
	flag.StringVar(&folder,"add","","add folders to scan")
	flag.StringVar(&email,"email","your@email.com","email to scan")
	flag.Parse()
	
	if folder!=""{
		scan(folder)
	}
	stats(email)
	// fmt.Println(flag.Args())
}
