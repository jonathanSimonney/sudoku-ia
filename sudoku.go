package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

//structure for the grid
type Grid struct{
	Values [9][9]int
}

//some methods to interact with our grid
func (this *Grid) addSolvingNumber(solving int, x int, y int) (bool, Grid){
	//first we'll copy the current grid to get a new one
	copiedGrid := Grid{}

	copier.Copy(&copiedGrid, &this)

	//check the box is empty
	if this.Values[y][x]  != 0{
		return false, copiedGrid
	}

	//check if the number is in the row
	if intInSlice(solving, this.Values[y]){
		return false, copiedGrid
	}

	//todo check also for square and column

	copiedGrid.Values[y][x] = solving
	return true, copiedGrid
}

func (this *Grid) prettyPrint(){
	for _, row := range this.Values {
		fmt.Println(row)
	}

	fmt.Println("*************************************************************************************************")
}

//helper function to check if int is in list
func intInSlice(a int, list [9]int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main()  {
	fmt.Println("hey")

	//filling the grid with numbers
	var currentGrid = Grid{Values:[9][9]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}}

	_, nextGrid := currentGrid.addSolvingNumber(6, 2, 0)

	currentGrid.prettyPrint()
	nextGrid.prettyPrint()
}
