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
	column := col{}
	for _, DayNumber := range keys{
		week := DayNumber/7
		weekday := DayNumber%7

		if weekday == 0{
			column = col{}
		}
		column = append(column, commits[DayNumber])

		if weekday == 6{
			cols[week] = column
		}
	}
	if len(column) > 0 {
		cols[keys[len(keys)-1]/7] = column
	}
	return cols
}

func printCells(cols map[int]col){
	printMonths()
	currOffset := CalcOffset()
	for i := 6; i >= 0; i--{
		for j := 27; j >= 0; j--{

			if j == 27{
				printCols(i)
			}
			if col,ok := cols[j]; ok{
				if j == 0 && i == currOffset {//special cell for today
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