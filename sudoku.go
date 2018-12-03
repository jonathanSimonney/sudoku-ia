package main

import (
	"fmt"
	"github.com/jinzhu/copier"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
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
func (this *Grid) addSolvingNumber(solving int, x int, y int, modifyOriginal bool) (bool, Grid){
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
		copiedGrid.Values[y][x] = solving
		return true, copiedGrid
	}
}

//a pretty print
func (this *Grid) prettyPrint(){
	for _, row := range this.Values {
		fmt.Println(row)
	}

	fmt.Println("*************************************************************************************************")
}

//a filler for the grid
func (this *Grid) fillGrid(howManyValues int){
	this.emptyGrid()

	//make sure random changes
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	nbRepetition := 0
	for nbRepetition != howManyValues{
		currentValue := r1.Intn(9) + 1
		x := r1.Intn(9)
		y := r1.Intn(9)

		addSuccessfull, _ := this.addSolvingNumber(currentValue, x, y, true)
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

//solver for the grid
func recursivelySolveGrid(grid Grid) (bool, Grid){
	for y := range make([]int, 9){
		for x := range make([]int, 9){
			if grid.Values[y][x] == 0{
				//fmt.Println("trying to fill ", x, y)
				//var successVar []int
				for solvingValueMinus := range make([]int, 9){
					solvingValue := solvingValueMinus + 1 //because it is currently between 0 and 8, not 1 and 9
					addSuccessfull, newGrid := grid.addSolvingNumber(solvingValue, x, y, false)
					if addSuccessfull{
						isSolved, solvedGrid := recursivelySolveGrid(newGrid)
						if isSolved{
							return true, solvedGrid
						}
					}
				}
				return false, grid
			}
		}
	}

	return true, grid
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

func programMain()  {
	//filling the grid with numbers
	var currentGrid = Grid{Values:[9][9]int{
		{4, 0, 2, 0, 1, 7, 0, 3, 0},
		{3, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 9, 0, 6, 0, 0},
		{0, 2, 0, 0, 4, 0, 7, 0, 8},
		{0, 0, 4, 0, 0, 0, 5, 0, 0},
		{1, 0, 7, 0, 6, 0, 0, 9, 0},
		{0, 0, 6, 0, 3, 0, 0, 0, 0},
		{9, 0, 0, 0, 0, 0, 0, 0, 3},
		{0, 4, 0, 1, 8, 0, 9, 0, 6},
	}}

	_, nextGrid := currentGrid.addSolvingNumber(9, 2, 6, false)

	currentGrid.prettyPrint()
	nextGrid.prettyPrint()

	//currentGrid.fillGrid(30)
	//currentGrid.prettyPrint()

	_, solvedGrid := recursivelySolveGrid(currentGrid)
	solvedGrid.prettyPrint()

	//bestGrid.prettyPrint()
}

func main(){
	f, err := os.Create("perf_cpu.perf")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	begin := time.Now()
	programMain()
	fmt.Println(time.Since(begin))


	f, err = os.Create("mem_profile.perf")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()
}
