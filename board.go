package main

const sizex = 9
const sizey = 9

type Board [sizey][sizex]int

func Transition(b Board) Board {
	c := Board{}

	for y, line := range b {
		for x, _ := range line {
			c[y][x] = getNextCell(b, x, y)
		}
	}

	return c
}

func getNextCell(b Board, x int, y int) int {
	neighbours := [8]int{
		get(b, x-1, y-1),
		get(b, x, y-1),
		get(b, x+1, y-1),

		get(b, x-1, y),
		get(b, x+1, y),

		get(b, x-1, y+1),
		get(b, x, y+1),
		get(b, x+1, y+1),
	}

	num := 0
	mode := -1
	for _, n := range neighbours {
		if n >= 0 {
			num += 1
			if mode != -1 {
				mode = (mode + n) / 2
			} else {
				mode = n
			}
		}
	}

	cellIsDead := (b[y][x] == -1)
	if cellIsDead {
		if num == 3 {
			// Birth
			return mode
		} else {
			// Death
			return -1
		}
	} else {
		if num <= 1 {
			// Death
			return -1
		} else if num == 2 || num == 3 {
			// Survive
			return mode
		} else if num >= 4 {
			// Death
			return -1
		}
	}

	return -1 // unexpected state
}

func get(b Board, x int, y int) int {
	if x < 0 || x >= sizex || y < 0 || y >= sizey {
		return -1
	} else {
		return b[y][x]
	}
}
