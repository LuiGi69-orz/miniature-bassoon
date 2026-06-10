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
	for i := 0; i < TotDays; i++ {
		Commits[i] = 0
	}
	for _, Repo := range AllRepo {
		Commits = FillIt(Commits, email, Repo)
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
