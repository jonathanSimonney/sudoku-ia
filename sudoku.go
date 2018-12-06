package main

import (
	"fmt"
	"sort"
	"strings"

	//"log"
	"math/rand"
	//"os"
	//"runtime"
	//"runtime/pprof"
	"time"
)

type byPossibility []Possibility

func (a byPossibility) Len() int           { return len(a) }
func (a byPossibility) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPossibility) Less(i, j int) bool { return len(a[i].Numbers) < len(a[j].Numbers) }


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

//a getter to get the valid or invalid values at a given position
func (this *Grid) getIncludingSets(x int, y int, returnValid bool) (mergedRet []int){
	if returnValid{
		mergedRet = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	}

	for i := 0; i < 9; i++{
		inLineForbidden := this.Values[y][i]
		if inLineForbidden != 0 && i != x{
			mergedRet = addOrRemoveFromUniqSlice(mergedRet, !returnValid, inLineForbidden)
		}

		inColumnForbidden := this.Values[i][x]
		if inColumnForbidden != 0 && i != y{
			mergedRet = addOrRemoveFromUniqSlice(mergedRet, !returnValid, inColumnForbidden)
		}
	}

	if (len(mergedRet) == 0 && returnValid) || (len(mergedRet) == 9 && !returnValid){
		return mergedRet
	}

	//first, know which square is targetted
	xSquareStart := (x / 3) * 3
	ySquareStart := (y / 3) * 3

	xSquareEnd := ((x / 3) * 3) + 3
	ySquareEnd := ((y / 3) * 3) + 3

	//then get all the numbers in this square
	for yIndex := ySquareStart; yIndex < ySquareEnd; yIndex++{
		for xIndex := xSquareStart; xIndex < xSquareEnd; xIndex++{
			forbidden := this.Values[yIndex][xIndex]
			if forbidden != 0 && (xIndex != x || yIndex != y){
				mergedRet = addOrRemoveFromUniqSlice(mergedRet, !returnValid, forbidden)
			}
		}
	}

	return mergedRet
}

//a getter of all the legals numbers at a given position
func (this *Grid) getLegalNumbersAtPos(x int, y int)(legalValues []int){
	//because no value can be legal if the box is already filled
	if this.Values[y][x] != 0{
		return legalValues
	}

	legalValues = this.getIncludingSets(x, y, true)

	return legalValues
}

//an adder of number which checks if the number can be added
func (this *Grid) addSolvingNumber(solving int, x int, y int, trustSolving bool) (isValid bool, modifiedGrid Grid){
	//first we'll copy the current grid to get a new one
	copiedGrid := *this

	//check the box is empty
	if this.Values[y][x]  != 0{
		return false, *this
	}

	if trustSolving{
		copiedGrid.Possibilities = this.Possibilities[1:]
		copiedGrid.Values[y][x] = solving
		return true, copiedGrid
	}else{
		invalidValues := this.getIncludingSets(x, y, false)

		if intInSlice(solving, invalidValues){
			return false, *this
		}

		this.Values[y][x] = solving
		return true, *this
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

		addSuccessfull, _ := this.addSolvingNumber(solvedGrid.Values[y][x], x, y, false)
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
func (this *Grid) prepare() (isValid bool){
	for y := range make([]int, 9){
		for x := range make([]int, 9){
			if this.Values[y][x] == 0{//because otherwhise there wouldn't be any possibility.
			    legalNumbers := this.getLegalNumbersAtPos(x, y)
			    if len(legalNumbers) == 1{
			    	this.Values[y][x] = legalNumbers[0]
				}else{
					this.Possibilities = append(this.Possibilities, Possibility{X: x, Y: y, Numbers: legalNumbers})
				}
			}else{
				forbiddenValues := this.getIncludingSets(x, y, false)
				if intInSlice(this.Values[y][x], forbiddenValues){
					return false
				}
			}
		}
	}

	sort.Sort(byPossibility(this.Possibilities))

	return true
}

//a helper to get the coords AND the list of elems of the best cell
func (this *Grid) getNextPossibility() (x int, y int, legalValues []int, isSolved bool){
	possibilityNumberLeft := len(this.Possibilities)

	if possibilityNumberLeft == 0{//there isn't any possibility left
		return -1, -1, []int{}, true
	}

	nextPossibility := this.Possibilities[0]

	return nextPossibility.X, nextPossibility.Y, nextPossibility.Numbers, false
}

//solver for the grid
func recursivelySolveGrid(grid Grid, randomly bool, firstTime bool) (isSolved bool, solvedGrid Grid){
	if firstTime{
		prepareBegin := time.Now()
		validGrid := grid.prepare()
		if !validGrid{
			return false, grid
		}
		fmt.Println("grid prepared in ")
		fmt.Println(time.Since(prepareBegin))
	}
	x, y, arrayIter, isSolved := grid.getNextPossibility()

	//it avoids to check if the number is valid in addSolvingNumber...
	if x != -1{
		arrayIter = grid.getLegalNumbersAtPos(x, y)
	}


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
		addSuccessfull, newGrid := grid.addSolvingNumber(solvingValue, x, y, true)
		if addSuccessfull{
			isSolved, solvedGrid := recursivelySolveGrid(newGrid, randomly, false)
			if isSolved{
				return true, solvedGrid
			}
		}
		grid.Values[y][x] = 0
	}

	//fmt.Println("sorry...")
	//if no value was sent, then the grid can't be solved.
	return false, grid


}

//helper function to check if int is in list
func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//helper function to suppress the first elem found in a list
func suppressFirstFromSlice(a int, list []int) (filteredList []int) {
	for index, b := range list {
		if b == a {
			list[index] = list[len(list)-1]
			return list[:len(list)-1]
		}
	}
	return list
}

//a helper to add or remove from a slice if elem is or isn't in there
func addOrRemoveFromUniqSlice(originalSlice []int, isFilling bool, newValue int) (treatedValues []int){
	if !isFilling{//we're emptying the array
		originalSlice = suppressFirstFromSlice(newValue, originalSlice)
	}else{
		if !intInSlice(newValue, originalSlice){
			originalSlice = append(originalSlice, newValue)
		}
	}

	return originalSlice
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
		//1 second
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 3, 0, 8, 5},
		//{0, 0, 1, 0, 2, 0, 0, 0, 0},
		//{0, 0, 0, 5, 0, 7, 0, 0, 0},
		//{0, 0, 4, 0, 0, 0, 1, 0, 0},
		//{0, 9, 0, 0, 0, 0, 0, 0, 0},
		//{5, 0, 0, 0, 0, 0, 0, 7, 3},
		//{0, 0, 2, 0, 1, 0, 0, 0, 0},
		//{0, 0, 0, 0, 4, 0, 0, 0, 9},

		//10 - 13 seconds
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 7, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 7, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{1, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 8, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 5, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},

		//nath long grid 12 minutes!!!
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 5, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 1, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},

		//invalid grid from start OK
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 5, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 1, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 1, 0, 0, 0, 0, 0, 0, 0},

		//insolvable simple grid CHECK
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 5, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 5},
		//{0, 0, 0, 0, 0, 0, 1, 0, 0},
		//{0, 0, 0, 0, 0, 5, 0, 0, 0},
		//{0, 5, 0, 0, 0, 0, 0, 0, 0},

		//insolvable originally long grid FALSE!!!
		{7, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 7, 0, 0},
		{0, 0, 0, 0, 1, 2, 0, 0, 0},
		{0, 0, 0, 7, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},

		//hardest in the world : 821 mSec
		//{1, 0, 0, 0, 0, 7, 0, 9, 0},
		//{0, 3, 0, 0, 2, 0, 0, 0, 8},
		//{0, 0, 9, 6, 0, 0, 5, 0, 0},
		//{0, 0, 5, 3, 0, 0, 9, 0, 0},
		//{0, 1, 0, 0, 8, 0, 0, 0, 2},
		//{6, 0, 0, 0, 0, 4, 0, 0, 0},
		//{3, 0, 0, 0, 0, 0, 0, 1, 0},
		//{0, 4, 0, 0, 0, 0, 0, 0, 7},
		//{0, 0, 7, 0, 0, 0, 3, 0, 0},
	}}

	currentGrid.prettyPrint()

	locBegin := time.Now()
	_, solvedGrid := recursivelySolveGrid(currentGrid, false, true)
	fmt.Println(time.Since(locBegin))


	solvedGrid.prettyPrint()
	//currentGrid.fillGrid(5)
	//currentGrid.prettyPrint()
	//_, solvedGrid = recursivelySolveGrid(currentGrid, false, true)
	//solvedGrid.prettyPrint()
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
