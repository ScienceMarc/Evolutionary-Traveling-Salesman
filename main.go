package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	numberOfLocations = 91
	generationSize    = 5000
	alphabet          = "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	mutationRate      = 0.0005
)

var (
	pointMap map[string]int
)

type point struct {
	x, y  int32
	label string
}

func main() {
	fitnessPoints := make([]float64, 0)
	var highestFitness float32
	var lowestFitness float32
	points := make([]point, numberOfLocations)
	rsource := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(rsource)
	for i := range points {
		points[i] = point{randomizer.Int31n(1280), randomizer.Int31n(640), string(alphabet[i])}
	}
	//fmt.Println(points)
	pointMap = make(map[string]int, numberOfLocations*numberOfLocations)
	for i, a := range points {
		for j, b := range points {
			dx := a.x - b.x
			dy := a.y - b.y
			pointMap[string(alphabet[i])+string(alphabet[j])] = int(dx*dx + dy*dy)
		}
	}
	fmt.Println(pointMap)
	var pop Population
	pop.createPopulation(randomizer)

	rl.InitWindow(1280, 720, "Travelling Salesman")
	//rl.SetTargetFPS(1)

	for !rl.WindowShouldClose() {
		bestPath, bestFitness, AvgFitness := pop.evaluate()
		bestPaths := pop.selection()
		pop.crossover(randomizer)
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		for _, p := range points {
			rl.DrawCircle(p.x, p.y, 5, rl.Red)
		}
		for x := 0; x < len(bestPaths); x++ {
			for i := 0; i < len(bestPaths[x].path)-1; i++ {

				start, end := rl.Vector2{}, rl.Vector2{}
				for _, p := range points {
					if p.label == bestPaths[x].path[i] {
						start = rl.Vector2{X: float32(p.x), Y: float32(p.y)}
					} else if p.label == bestPaths[x].path[i+1] {
						end = rl.Vector2{X: float32(p.x), Y: float32(p.y)}
					}
				}
				rl.DrawLineV(start, end, rl.LightGray)
			}
		}
		for i := 0; i < len(bestPath)-1; i++ {
			start, end := rl.Vector2{}, rl.Vector2{}
			for _, p := range points {
				if p.label == bestPath[i] {
					start = rl.Vector2{X: float32(p.x), Y: float32(p.y)}
				} else if p.label == bestPath[i+1] {
					end = rl.Vector2{X: float32(p.x), Y: float32(p.y)}
				}
			}
			rl.DrawLineEx(start, end, 3.5, rl.Blue)
		}

		rl.DrawText("Best:      "+fmt.Sprint(bestFitness), 5, int32(rl.GetScreenHeight())-45, 20, rl.Gray)
		rl.DrawText("Average: "+fmt.Sprint(AvgFitness), 5, int32(rl.GetScreenHeight())-25, 20, rl.Gray)

		mapFloat := func(x, in_min, in_max, out_min, out_max float32) float32 {
			return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
		}
		highestFitness = float32(math.Max(float64(highestFitness), bestFitness))
		lowestFitness = float32(math.Min(float64(lowestFitness), bestFitness))
		fitnessPoints = append(fitnessPoints, bestFitness)
		graph := make([]rl.Vector2, 0)
		for i, v := range fitnessPoints {
			x := mapFloat(float32(i), 0, float32(len(fitnessPoints)), float32(rl.GetScreenWidth())/2, float32(rl.GetScreenWidth()))
			y := mapFloat(float32(v), lowestFitness, highestFitness, float32(rl.GetScreenHeight()-5), float32(rl.GetScreenHeight()-65))
			graph = append(graph, rl.Vector2{X: x, Y: y})
		}
		for i := 0; i < len(graph)-1; i++ {
			rl.DrawLineV(graph[i], graph[i+1], rl.Gray)
		}

		rl.EndDrawing()
		//time.Sleep(1 * time.Second)
	}
	rl.CloseWindow()
}
