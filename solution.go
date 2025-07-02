package main

import (
	"bufio"
	. "fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type coordinates struct {
	x, y int
}

type Piano struct {
	automata map[string]coordinates
	numberOfAutomata map[coordinates]int
	obstacles [][2]coordinates
}

type piano *Piano

func newPiano() piano {
	return &Piano{automata: make(map[string]coordinates),
		obstacles: make([][2]coordinates, 0),
		numberOfAutomata: make(map[coordinates]int)}
}

func distance(point1 coordinates, point2 coordinates) int {
	return int(math.Abs(float64(point2.x-point1.x)) + math.Abs(float64(point2.y-point1.y)))
}

func recall(p piano, prefix string, point coordinates) {
	var minDistance int = -1
	var names []string

	for key, pos := range p.automata {
		if strings.HasPrefix(key, prefix) && findPath(p, pos, point, make(map[coordinates]bool)) {
			d := distance(pos, point);
			if d == 0 {
				continue
			}
			if d < minDistance || minDistance == -1 {
				minDistance = d
				names = []string{key}
			} else if d == minDistance {
				names = append(names, key)
			}
		}
	}
	for _, name := range names {
		pos := p.automata[name]
		p.numberOfAutomata[pos] -= 1
		if p.numberOfAutomata[pos] == 0 {
			delete(p.numberOfAutomata, pos)
		}
		p.automata[name] = point
		p.numberOfAutomata[point] += 1
	}
}

func findPath(p piano, start coordinates, goal coordinates, visited map[coordinates]bool) bool {
	visited[start] = true
	if start == goal {
		return true
	}
	
	directions := [2]coordinates {
		findNextX(p, start, goal),
		findNextY(p, start, goal),
	}

	for _, direction := range directions {
		if !visited[direction] {
			if findPath(p, direction, goal, visited) {
				return true
			}
			if direction.x != start.x {
				step := step(direction.x - start.x)
				for x := direction.x - step; x != start.x; x -= step {
					newPos := coordinates{x, start.y}
					if !visited[newPos] {
						if findPath(p, newPos, goal, visited) {
							return true
						}
					}
				}
			} else if direction.y != start.y {
				step := step(direction.y - start.y)
				for y := direction.y - step; y != start.y; y -= step {
					newPos := coordinates{start.x, y}
					if !visited[newPos] {
						if findPath(p, newPos, goal, visited) {
							return true
						}
					}
				}
			}	
		}
	}
	return false
}

func step(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func findNextX(p piano, start, goal coordinates) coordinates {
	step := step(goal.x - start.x)
	target := coordinates{goal.x, start.y}

	for _, obstacle := range p.obstacles {
		if start.y >= obstacle[0].y && start.y <= obstacle[1].y {
			if step > 0 && obstacle[0].x > start.x && obstacle[0].x <= target.x {
				target.x = obstacle[0].x - 1
			} else if step < 0 && obstacle[1].x < start.x && obstacle[1].x >= target.x {
				target.x = obstacle[1].x + 1
			}
		}
	}
	return target
}

func findNextY(p piano, start, goal coordinates) coordinates {
	step := step(goal.y - start.y)
	target := coordinates{start.x, goal.y}

	for _, obstacle := range p.obstacles {
		if start.x >= obstacle[0].x && start.x <= obstacle[1].x {
			if step > 0 && obstacle[0].y > start.y && obstacle[0].y <= target.y {
				target.y = obstacle[0].y - 1
			} else if step < 0 && obstacle[1].y < start.y && obstacle[1].y >= target.y {
				target.y = obstacle[1].y + 1
			}
		}
	}
	return target
}

func existsPath(p piano, point coordinates, name string) {
	if _, exists := p.automata[name]; !exists || pointIsInObstacle(p, point) {
		Println("NO")
	} else if d := distance(p.automata[name], point); d == 0 {
		Println("NO")
	} else if findPath(p, p.automata[name], point, make(map[coordinates]bool)) {
		Println("SI")
	} else {
		Println("NO")
	}
}

func obstacle(p piano, point1 coordinates, point2 coordinates) {
	if !pointIsInAutomaton(p, point1) && !pointIsInAutomaton(p, point2) {
		p.obstacles = append(p.obstacles, [2]coordinates{point1, point2})
	}
}

func automaton(p piano, point coordinates, name string) {
	if pointIsInObstacle(p, point) {
		return
	}
	if _, exists := p.automata[name]; !exists {
		p.automata[name] = point
		p.numberOfAutomata[point] += 1
	}
}

func printPiano(p piano) {
	Println("(")
	for key, pos := range p.automata {
		Print(key, ": ", pos.x, ",", pos.y, "\n")
	}
	Println(")\n[")
	for _, obstacle := range p.obstacles {
		Print("(", obstacle[0].x, ",", obstacle[0].y, ")(", obstacle[1].x, ",", obstacle[1].y, ")\n")
	}
	Println("]")
}

func positions(p piano, prefix string) {
	Println("(")
	for key, pos := range p.automata {
		if strings.HasPrefix(key, prefix) {
			Print(key, ": ", pos.x, ",", pos.y, "\n")
		}
	}
	Println(")")
}

func pointIsInObstacle(p piano, point coordinates) bool {
	for _, obstacle := range p.obstacles {
		if point.x >= obstacle[0].x && point.x <= obstacle[1].x &&
			point.y >= obstacle[0].y && point.y <= obstacle[1].y {
			return true
		}
	}
	return false
}

func pointIsInAutomaton(p piano, point coordinates) bool {
	return p.numberOfAutomata[point] > 0
}

func state(p piano, point coordinates) {
	if pointIsInObstacle(p, point) {
		Println("O")
		return
	}
	if pointIsInAutomaton(p, point) {
		Println("A")
		return
	}
	Println("E")
}

func esegui(p piano, s string) {
	switch {
	case s == "c":
		p = newPiano()
	case s == "S":
		printPiano(p)
	case s == "f":
		os.Exit(0)
	default:
		splitted := strings.Split(s, " ")
		switch splitted[0] {
		case "s":
			n1, _ := strconv.Atoi(splitted[1])
			n2, _ := strconv.Atoi(splitted[2])
			state(p, coordinates{n1, n2})
		case "a":
			n1, _ := strconv.Atoi(splitted[1])
			n2, _ := strconv.Atoi(splitted[2])
			automaton(p, coordinates{n1, n2}, splitted[3])
		case "o":
			n1, _ := strconv.Atoi(splitted[1])
			n2, _ := strconv.Atoi(splitted[2])
			n3, _ := strconv.Atoi(splitted[3])
			n4, _ := strconv.Atoi(splitted[4])
			obstacle(p, coordinates{n1, n2}, coordinates{n3, n4})
		case "p":
			positions(p, splitted[1])
		case "e":
			n1, _ := strconv.Atoi(splitted[1])
			n2, _ := strconv.Atoi(splitted[2])
			existsPath(p, coordinates{n1, n2}, splitted[3])
		case "r":
			n1, _ := strconv.Atoi(splitted[1])
			n2, _ := strconv.Atoi(splitted[2])
			recall(p, splitted[3], coordinates{n1, n2})
		}
	}
}

func main() {
	var p = newPiano()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s := scanner.Text()
		esegui(p, s)
	}
}
