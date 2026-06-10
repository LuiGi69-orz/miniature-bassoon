package main

import (
	"fmt"
	// "math"
	"sort"
	"time"
)

func PrintHistory(CommitsHistory map[int]int) {
	keys := Sort(CommitsHistory)
    cols := BuildCols(keys, CommitsHistory)
    printCells(cols)
}

func Sort(commits map[int]int) []int{
	var res []int
	for key := range commits{
		res = append(res,key)
	}
	sort.Ints(res)
	return res
}

type col []int
func BuildCols(keys []int,commits map[int]int) map[int]col{
	cols := make(map[int]col)
	today := int(time.Now().Weekday())
	for _, daysAgo := range keys{
		weekday := (today - (daysAgo % 7) + 7) % 7 //weekday := int(date.Weekday())
        week := (daysAgo + today) / 7

		if _, ok := cols[week]; !ok {
            cols[week] = make(col, 7)
        }
        cols[week][weekday] = commits[daysAgo]
	}
	return cols
}

func printCells(cols map[int]col){
	printMonths()
	today := int(time.Now().Weekday())
	for i := 0; i <= 6; i++{
		for j := 27; j >= 0; j--{

			if j == 27{
				printCols(i)
			}
			if col,ok := cols[j]; ok{
				if j == 0 && i == today {//special cell for today
					PrintCell(col[i],true)
					continue
				}else {
					if len(col) > i{ //if we access col[j] being 0 or not and j>len(col) then panic occurs
						PrintCell(col[i],false)
						continue
					}
				}
			}
			PrintCell(0,false)//better to keep this way for days with no commits
		}
		fmt.Printf("\n")
	}
}
func printMonths(){
	week := CalcPresentDate(time.Now()).Add(-(TotDays * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")

	for{
		if week.After(time.Now()) {
			break
		}

		if week.Month() != month {
			fmt.Printf("%s",week.Month().String()[:3])
			month = week.Month()
		}else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
	}
	fmt.Printf("\n")
}

func printCols(day int){
	switch day{
		case 0:
			fmt.Printf("Sun")
		case 1 :
			fmt.Printf("Mon")
		case 2:
			fmt.Printf("Tue")
		case 3:
			fmt.Printf("Wed")
		case 4 :
			fmt.Printf("Thu")
		case 5:
			fmt.Printf("Fri")
		case 6:
			fmt.Printf("Sat")
	}
}

func PrintCell(day int , today bool){ 
    escape := "\033[0;37;30m"
    switch {
		case day > 0 && day < 5:
			escape = "\033[1;30;47m"
		case day >= 5 && day < 10:
			escape = "\033[1;30;43m"
		case day >= 10:
			escape = "\033[1;30;42m"
    }
    if today {
        escape = "\033[1;37;45m"
    }
    if day == 0 {
        fmt.Printf("%s  - %s", escape, "\033[0m")
        return
    }
    str := "  %d "
    switch {
		case day >= 10:
			str = " %d "
		case day >= 100:
			str = "%d "
    }
    fmt.Printf(escape + str + "\033[0m", day) 
}


// The simpler approach is conceptually cleaner:
// Don't store commits by "days ago". Store them by actual date.

// Instead Commits[daysAgo]++

// which forces  to constantly convert:
// daysAgo -> weekday
// daysAgo -> week column
// daysAgo -> calendar date
// and that's where all the offset headaches come from.

// Store by date instead

// Change the map to:
// map[time.Time]int

// where the key is the start of the day:
// day := CalcPresentDate(c.Author.When)
// Commits[day]++


// In ProcessRepo
// func ProcessRepo(email string) map[time.Time]int {
//     commits := make(map[time.Time]int)
//     ...
// }


// In FillIt
// Instead of:
// daysAgo := CountDist(c.Author.When)
// if daysAgo < TotDays {
//     Commits[daysAgo]++
// }

// do:
// day := CalcPresentDate(c.Author.When)
// cutoff := CalcPresentDate(time.Now()).AddDate(0, 0, -TotDays)
// if !day.Before(cutoff) {
//     Commits[day]++
// }

// Building columns becomes trivial
// For each date:
// weekday := int(day.Weekday())
// No offset.
// No modulo tricks.
// No dependence on today's weekday.

// The weekday is literally stored in the date.

// Computing the week column
// Let
// start := CalcPresentDate(time.Now()).AddDate(0, 0, -TotDays)

// Then:
// days := int(day.Sub(start).Hours() / 24)
// week := days / 7


// Example

// Suppose today is Thursday June 11.
// A commit from Wednesday June 10:

// day = 2026-06-10
// weekday := int(day.Weekday())

// gives:
// 3 (Wednesday)

// forever.
// Tomorrow, next week, next month:
// day.Weekday()
// still returns Wednesday.
// Nothing shifts.