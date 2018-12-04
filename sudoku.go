package main

import (
	"fmt"
	"github.com/jinzhu/copier"
	"strings"

	//"log"
	"math/rand"
	//"os"
	//"runtime"
	//"runtime/pprof"
	"time"
)

// structure for a possibility
type Possibility struct{
	Numbers []int
	X int
	Y int
}

//structure for the grid
type Grid struct{
	Values [9][9]int
	Possibilities []Possibility
}

//some methods to interact with our grid

//a getter to get the row, the column and the square of two given coords. Couldn't find a really fitting name...
func (this *Grid) getIncludingSets(x int, y int) (rowRet [9]int, colRet [9]int, squareRet [9]int){
	rowRet = this.Values[y]

	//let's get the whole column as a slice
	for index, row := range this.Values{
		colRet[index] = row[x]
	}

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

//a getter of all the legals numbers at a given position
func (this *Grid) getLegalNumbersAtPos(x int, y int)(legalValues []int){
	//because no value can be legal if the box is already filled
	if this.Values[y][x] != 0{
		return legalValues
	}

	arrayIter := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rowForbidden, colForbidden, squareForbidden := this.getIncludingSets(x, y)

	for _, solvingValue := range arrayIter{
		if !intInSlice(solvingValue, rowForbidden) && !intInSlice(solvingValue, colForbidden) && !intInSlice(solvingValue, squareForbidden){
			legalValues = append(legalValues, solvingValue)
		}
	}

	return legalValues
}

//an adder of number which checks if the number can be added
func (this *Grid) addSolvingNumber(solving int, x int, y int, modifyOriginal bool) (isValid bool, modifiedGrid Grid){
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

	if modifyOriginal{
		this.Values[y][x] = solving
		return true, *this
	}else{
		copiedGrid.Possibilities = copiedGrid.Possibilities[:len(copiedGrid.Possibilities)-1]
		copiedGrid.Values[y][x] = solving
		return true, copiedGrid
	}
}

//a pretty print
func (this *Grid) prettyPrint(){
	fmt.Println(strings.Repeat("_ ", 13))
	for indexY, row := range this.Values {
		fmt.Print("|")
		for indexX, number := range row{
			fmt.Print(" ", number)
			if (indexX + 1) % 3 == 0{
				fmt.Print(" |")
			}
		}
		if (indexY + 1) % 3 == 0{
			fmt.Print("\n")
			fmt.Println(strings.Repeat("_ ", 13))
		}else{
			fmt.Print("\n")
		}
	}

	fmt.Println("*************************************************************************************************")
}

//a filler for the grid
func (this *Grid) fillGrid(howManyValues int){
	this.emptyGrid()

	_, solvedGrid := recursivelySolveGrid(*this, true, true)
	//make sure random changes
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	nbRepetition := 0
	for nbRepetition != howManyValues{
		x := r1.Intn(9)
		y := r1.Intn(9)

		addSuccessfull, _ := this.addSolvingNumber(solvedGrid.Values[y][x], x, y, true)
		if addSuccessfull{
			nbRepetition++
		}
	}
}

//a helper to empty the grid
func (this *Grid) emptyGrid(){
	for y := range make([]int, 9){
		for x := range make([]int, 9){
			this.Values[y][x] = 0
		}
	}
}

//a helper to fill the possibles list for our grid
func (this *Grid) preparePossibles(){
	for y := range make([]int, 9){
		for x := range make([]int, 9){
			if this.Values[y][x] == 0{//because otherwhise there wouldn't be any possibility.
				this.Possibilities = append(this.Possibilities, Possibility{X: x, Y: y, Numbers: this.getLegalNumbersAtPos(x, y)})
			}
		}
	}
}

//a helper to get the coords AND the list of elems of the best cell
func (this *Grid) getNextPossibility() (x int, y int, legalValues []int, isSolved bool){
	possibilityNumberLeft := len(this.Possibilities)

	if possibilityNumberLeft == 0{//there isn't any possibility left
		return -1, -1, []int{}, true
	}

	nextPossibility := this.Possibilities[possibilityNumberLeft - 1]

	return nextPossibility.X, nextPossibility.Y, nextPossibility.Numbers, false
}

//solver for the grid
func recursivelySolveGrid(grid Grid, randomly bool, firstTime bool) (isSolved bool, solvedGrid Grid){
	if firstTime{
		grid.preparePossibles()
	}
	x, y, arrayIter, isSolved := grid.getNextPossibility()

	//if the sudoku is already filled, nothing to do
	if isSolved{
		return true, grid
	}

	//shuffle the array of numbers to test in case the user wants random
	if randomly{
		shuffle(arrayIter)
	}

	//test each value in the array of valid numbers.
	for _, solvingValue := range arrayIter{
		addSuccessfull, newGrid := grid.addSolvingNumber(solvingValue, x, y, false)
		if addSuccessfull{
			isSolved, solvedGrid := recursivelySolveGrid(newGrid, randomly, false)
			if isSolved{
				return true, solvedGrid
			}
		}
	}

	//fmt.Println("sorry...")
	//if no value was sent, then the grid can't be solved.
	return false, grid


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

//helper function to shuffle an array
func shuffle(vals []int){
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}


func programMain()  {
	//filling the grid with numbers
	var currentGrid = Grid{Values:[9][9]int{
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 3, 0, 8, 5},
		//{0, 0, 1, 0, 2, 0, 0, 0, 0},
		//{0, 0, 0, 5, 0, 7, 0, 0, 0},
		//{0, 0, 4, 0, 0, 0, 1, 0, 0},
		//{0, 9, 0, 0, 0, 0, 0, 0, 0},
		//{5, 0, 0, 0, 0, 0, 0, 7, 3},
		//{0, 0, 2, 0, 1, 0, 0, 0, 0},
		//{0, 0, 0, 0, 4, 0, 0, 0, 9},

		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 6, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}}

	currentGrid.prettyPrint()

	locBegin := time.Now()
	_, solvedGrid := recursivelySolveGrid(currentGrid, false, true)
	fmt.Println(time.Since(locBegin))


	solvedGrid.prettyPrint()
	currentGrid.fillGrid(5)
	currentGrid.prettyPrint()
	_, solvedGrid = recursivelySolveGrid(currentGrid, false, true)
	solvedGrid.prettyPrint()
}

func main(){
	//f, err := os.Create("perf_cpu.perf")
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}
	//defer pprof.StopCPUProfile()
	//
	begin := time.Now()
	programMain()
	fmt.Println(time.Since(begin))
	//
	//
	//f, err = os.Create("mem_profile.perf")
	//if err != nil {
	//	log.Fatal("could not create memory profile: ", err)
	//}
	//runtime.GC() // get up-to-date statistics
	//if err := pprof.WriteHeapProfile(f); err != nil {
	//	log.Fatal("could not write memory profile: ", err)
	//}
	//f.Close()
}
