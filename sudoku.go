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

//a getter to get the row, the column and the square of two given coords. Couldn't find a really fitting name...
func (this *Grid) getIncludingSets(x int, y int) ([9]int, [9]int, [9]int){
	rowRet := this.Values[y]

	//let's get the whole column as a slice
	var colRet [9]int

	for index, row := range this.Values{
		colRet[index] = row[x]
	}

	var squareRet [9]int

	squareIndex := 0

	//first, know which square is targetted
	xSquareStart := (x / 3) * 3
	ySquareStart := (y / 3) * 3

	//then get all the numbers in this square
	for xIncrement, _ := range make([]int, 3){
		for yIncrement, _ := range make([]int, 3){
			squareRet[squareIndex] = this.Values[ySquareStart + yIncrement][xSquareStart + xIncrement]
			squareIndex++
		}
	}

	return rowRet, colRet, squareRet
}

//an adder of number which checks if the number can be added
func (this *Grid) addSolvingNumber(solving int, x int, y int) (bool, Grid){
	//first we'll copy the current grid to get a new one
	copiedGrid := Grid{}

	copier.Copy(&copiedGrid, &this)

	//check the box is empty
	if copiedGrid.Values[y][x]  != 0{
		return false, copiedGrid
	}

	rowForbidden, colForbidden, squareForbidden := copiedGrid.getIncludingSets(x, y)

	if intInSlice(solving, rowForbidden) || intInSlice(solving, colForbidden) || intInSlice(solving, squareForbidden){
		return false, copiedGrid
	}

	copiedGrid.Values[y][x] = solving
	return true, copiedGrid
}

//a pretty print
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

	_, nextGrid := currentGrid.addSolvingNumber(9, 2, 6)

	currentGrid.prettyPrint()
	nextGrid.prettyPrint()
}
