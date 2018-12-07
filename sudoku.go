package main

import (
	"fmt"
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
	LinkedPossibilities []*Possibility
}

//structure for the grid
type Grid struct{
	Values [9][9]int
	Possibilities []*Possibility
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

//a helper to check if a value is valid at a position with performance in mind
func (this *Grid) isValueValidThere(x int, y int, solving int) (isValid bool){
	for i := 0; i < 9; i++{
		inLineForbidden := this.Values[y][i]
		if inLineForbidden == solving{
			return false
		}

		inColumnForbidden := this.Values[i][x]
		if inColumnForbidden == solving{
			return false
		}
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
			if forbidden == solving{
				return false
			}
		}
	}

	return true
}

//helper to suppress a value from the possibility linked to one whith the value assigned to the possibility
func (this *Grid) suppressFromLinkedPossibility(possibility Possibility, assignedValue int){
	for _, linkedPossibility := range possibility.LinkedPossibilities{
		linkedPossibility.Numbers = suppressFirstFromSlice(assignedValue, linkedPossibility.Numbers)
	}
}

//an adder of number which checks if the number can be added
func (this *Grid) addSolvingNumber(solving int, x int, y int, trustSolving bool) (isValid bool, modifiedGrid Grid){
	//check the box is empty
	if this.Values[y][x]  != 0{
		return false, *this
	}

	if !trustSolving && !this.isValueValidThere(x, y, solving){
		return false, *this
	}

	this.Values[y][x] = solving
	return true, *this
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

	_, solvedGrid := recursivelySolveGrid(*this, true, true, 0)
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

//a helper to add the subPossibility array to possibilities in grid
func (this* Grid) addPossibilitiesSubPossibilities(){
	var futurePossibilities []*Possibility
	for _, possibility := range this.Possibilities{

		futurePossibilities = append(futurePossibilities, possibility)
	}

	this.Possibilities = futurePossibilities
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
					this.Possibilities = append(this.Possibilities, &Possibility{X: x, Y: y, Numbers: legalNumbers})
				}
			}else{
				forbiddenValues := this.getIncludingSets(x, y, false)
				if intInSlice(this.Values[y][x], forbiddenValues){
					return false
				}
			}
		}
	}

	//sort.Sort(byPossibility(this.Possibilities))

	this.addPossibilitiesSubPossibilities()

	fmt.Println(this.Possibilities)
	panic("stop")

	return true
}

//solver for the grid
func recursivelySolveGrid(grid Grid, randomly bool, firstTime bool, index int) (isSolved bool, solvedGrid Grid){
	if firstTime{
		prepareBegin := time.Now()
		validGrid := grid.prepare()
		if !validGrid{
			return false, grid
		}
		fmt.Println("grid prepared in ")
		fmt.Println(time.Since(prepareBegin))
	}

	//if the sudoku is already filled, nothing to do
	if len(grid.Possibilities) == index{
		return true, grid
	}

	nextPossibility := grid.Possibilities[index]

	var arrayIter []int
	x := nextPossibility.X
	y := nextPossibility.Y
	//arrayIter = nextPossibility.Numbers

	//x, y, arrayIter, isSolved := grid.getPossibility(index)

	//it avoids to check if the number is valid in addSolvingNumber...
	//if x != -1{
	//	arrayIter = grid.getLegalNumbersAtPos(x, y)
	//}

	//shuffle the array of numbers to test in case the user wants random
	if randomly{
		shuffle(arrayIter)
	}

	//test each value in the array of valid numbers.
	for _, solvingValue := range arrayIter{
		currentPossibilities := grid.Possibilities
		addSuccessfull, newGrid := grid.addSolvingNumber(solvingValue, x, y, true)
		if addSuccessfull{
			isSolved, solvedGrid := recursivelySolveGrid(newGrid, randomly, false, index + 1)
			if isSolved{
				return true, solvedGrid
			}
		}
		grid.Values[y][x] = 0
		grid.Possibilities = currentPossibilities
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
		//20 - 30 second OR 1 second with sort
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 3, 0, 8, 5},
		//{0, 0, 1, 0, 2, 0, 0, 0, 0},
		//{0, 0, 0, 5, 0, 7, 0, 0, 0},
		//{0, 0, 4, 0, 0, 0, 1, 0, 0},
		//{0, 9, 0, 0, 0, 0, 0, 0, 0},
		//{5, 0, 0, 0, 0, 0, 0, 7, 3},
		//{0, 0, 2, 0, 1, 0, 0, 0, 0},
		//{0, 0, 0, 0, 4, 0, 0, 0, 9},

		//infinite (with or without sort)
		{0, 0, 0, 7, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{7, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 7, 0, 0},
		{0, 0, 0, 0, 1, 2, 0, 0, 0},


		//1 ms
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 7, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 7, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//{1, 0, 0, 0, 0, 0, 0, 0, 0},
		//{0, 0, 8, 0, 0, 0, 0, 0, 0},
		//{0, 0, 0, 0, 0, 0, 5, 0, 0},
		//{0, 0, 0, 0, 0, 0, 0, 0, 0},

		//nath long grid 1 ms
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
		//{6, 5, 7, 8, 9, 2, 0, 3, 4},

		//hardest in the world : 3 mSec
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
	_, solvedGrid := recursivelySolveGrid(currentGrid, false, true, 0)
	fmt.Println(time.Since(locBegin))


	solvedGrid.prettyPrint()
	//currentGrid.fillGrid(5)
	//currentGrid.prettyPrint()
	//_, solvedGrid = recursivelySolveGrid(currentGrid, false, true, 0)
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
