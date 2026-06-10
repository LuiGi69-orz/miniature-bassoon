package main

import (
	"bufio"
	"fmt"
	"io" 
	"log"
	"os"
	"os/user"
	"strings"
)

func scan(path string){
	fmt.Printf("Found Folders :\n\n")
	NewRepo := GenRepoPath(make([]string,0),path)
	DotFilePath := GetDotFile()
	Merge(NewRepo,DotFilePath)
	fmt.Printf("Scanned all Folders successfully !\n\n")
}

func GenRepoPath(folders []string,path string) []string{
	// path = strings.TrimSuffix(path,"/")

	fptr,err := os.Open(path)
	if err != nil{
		log.Fatal(err)
	}
	defer fptr.Close()

	entry,err := fptr.ReadDir(-1)
	if err != nil{
		log.Fatal(err)
	}

	for _ , entries := range entry{
		if(entries.IsDir()){
			var file string 
			file = path + "/" + entries.Name()
			if entries.Name() == ".git"{
				file = strings.TrimSuffix(file,"/.git")
				folders = append(folders,file)
				continue
			}
			if entries.Name() == "vendor" || entries.Name() == "node_modules"{
				continue
			}
			folders = GenRepoPath(folders,file)
		}
	}
	return folders
}

func GetDotFile() string{
	usr,err := user.Current()
	if err != nil{
		log.Fatal(err)
	}
	var path = usr.HomeDir + "/.gogitloccalstats"
	return path
}

func Merge(NewRepo []string,DotFilePath string) {
	ExistingRepo := Read(DotFilePath)
	CombinedRepo := joinSlices(NewRepo,ExistingRepo)
	Content := strings.Join(CombinedRepo, "\n")
	os.WriteFile(DotFilePath , []byte(Content), 0755)
}

func Read(path string) []string{
	fptr := Open(path)
	defer fptr.Close()

	var lines []string
	scan := bufio.NewScanner(fptr);
	for scan.Scan(){
		lines = append(lines,scan.Text())
	}
	if err := scan.Err(); err != nil{
		if err != io.EOF{
			panic(err)
		}
	}
	return lines
}

func Open(path string) *os.File{
	fptr,err := os.OpenFile(path,os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err){
			var err1 error
			fptr ,err1 = os.Create(path)
			if err1 != nil{
				panic(err1)
			}
		} else {
			panic(err)
		}
	}
	return fptr
}

func joinSlices(New,Old []string) []string{
	for _,path := range New{
		if SliceContains(Old,path){
			continue
		}
		Old = append(Old, path)
	}
	return Old
}

func SliceContains(slice []string,path string) bool{
	for _,it := range slice{
		if it == path{
			return true
		}
	}
	return false
}