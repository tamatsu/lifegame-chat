package main

const sizex = 9
const sizey = 9
type Board [sizex][sizey]int

func Transition(b Board) Board {
	c := Board{}
	
	for y, line := range b {
		for x, _ := range line {
			c[x][y] = CellNext(b, x, y)
		}
	}
		
	return c
}

func CellNext(b Board, x int, y int) int {
	neighbor := [8]int{
		get(b, x-1, y-1),
		get(b, x, y-1),
		get(b, x+1, y-1),
		
		get(b, x-1, y),
		get(b, x+1, y),
		
		get(b, x-1, y+1),
		get(b, x, y+1),
		get(b, x+1, y+1),
	}
	
  bucket := make(map[int]int)
	sum := 0
  mode := -1
	for _, n := range neighbor {
    if (n >= 0) {
      sum += 1
      if (mode != -1) {
        mode = (mode + n) / 2 
      } else {
        mode = n
      }
    }

    bucket[n]++
	}

  // mode := -1
  // for key, n := range bucket {
  //   if (n > mode) {
  //     mode = key
  //   }
  // }

	
	if (b[x][y] == -1) {
		if (sum == 3) {
			// Birth
			return mode
		} else {
			// Nothing
			return -1
		}
	} else {
		if (sum <= 1) {
			// Death
			return -1
		} else if (sum == 2 || sum == 3) {
			// Survive
			return mode
		} else if (sum >= 4) {
			// Death
			return -1
		}
	}
	
	return -1 // unexpected state
}
		
func get(b Board, x int, y int) int {
	if (x < 0 || x >= sizex || y < 0 || y >= sizey) {
		return -1
	} else {
		return b[x][y]
	}
}
