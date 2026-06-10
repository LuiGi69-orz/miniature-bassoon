package main

import (
	"fmt"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func stats(email string) {
	fmt.Println("Given mail -> ", email)
	CommitsHistory := ProcessRepo(email)
	PrintHistory(CommitsHistory)
}

const TotDays = 183

func ProcessRepo(email string) map[int]int {
	DotFilePath := GetDotFile()
	AllRepo := Read(DotFilePath)

	Commits := make(map[int]int, TotDays)

	//When printstats.go loops through the sorted keys to build the columns, 
	// any missing days completely throw off the calendar grid calculation. 
	// If you go three days without committing, the grid slips out of alignment 
	// 			because : 
	//			it assumes the data points are perfectly consecutive.
	
	for i := 0; i < TotDays; i++ {
		Commits[i] = 0
	}

	for _, Repo := range AllRepo {
		Commits = FillIt(Commits, email, Repo)  //no refrencing only passing by value 
	}
	return Commits
}

func FillIt(Commits map[int]int, email, Repo string) map[int]int {
	repo, err := git.PlainOpen(Repo)
	if err != nil {
		return Commits  // If a repo was moved or deleted, skip it instead of crashing
	}
	ref, err := repo.Head()
	if err != nil {
		//FIX: If the repo is empty (no commits yet), HEAD doesn't exist.
		// We catch the "reference not found" error here and skip the repo safely
		return Commits
	}
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return Commits
	}

	OffSet := CalcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := CountDist(c.Author.When) + OffSet 
		if c.Author.Email != email {
			return nil
		}
		if daysAgo < TotDays {
			Commits[daysAgo]++
		}
		return nil
	})

	return Commits
}
//Your grid has fixed rows for the days of the week:
// Row 0 = Sunday
// Row 1 = Monday
// Row 2 = Tuesday
// Row 3 = Wednesday... and so on.
// If you don't use an offset, 
// the program uses the raw "days ago" number to determine the row using DayNumber % 7. 
// Look at what would happen to a commit you make : 
//		Today(which is a Wednesday):
//			Today's commit = 0 days ago => 0 % 7 = 0 -> The program prints today's commit on the Sunday row.
// 		Tomorrow (Thursday), that same commit becomes 1 day ago:1 % 7 = 1 ->
// 			 The program now prints it on the Monday row.
// 
// Every time you wake up, your entire history would slide down by one row!

func CalcOffset() int {
	var OffSet int
	switch time.Now().Weekday() {
	case time.Sunday:
		OffSet = 0
	case time.Monday:
		OffSet = 1
	case time.Tuesday:
		OffSet = 2
	case time.Wednesday:
		OffSet = 3
	case time.Thursday:
		OffSet = 4
	case time.Friday:
		OffSet = 5
	case time.Saturday:
		OffSet = 6
	}
	return OffSet
}

func CountDist(date time.Time) int {
	days := 0
	Present := CalcPresentDate(time.Now())

	for date.Before(Present) {
		date = date.Add(time.Hour * 24)
		days++

		if days > TotDays {
			return 99999 //to indicate out of range
		}
	}
	return days
}

func CalcPresentDate(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location()) 
	return startOfDay
}
