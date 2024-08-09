package main

import (
	"fmt"
	"gifhelper"
	"os"
)

func main() {
	mode := os.Args[1]

	frequency := 1000
	scalingFactor := 1e11 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
	// the following sample parameters may be helpful for the "collide" command
	// all units are in SI (meters, kg, etc.)
	// but feel free to change the positions of the galaxies.

	g0 := InitializeGalaxy(500, 4e21, 7e22, 2e22)
	g1 := InitializeGalaxy(500, 4e21, 2e22, 7e22)

	// you probably want to apply a "push" function at this point to these galaxies to move
	// them toward each other to collide.
	// be careful: if you push them too fast, they'll just fly through each other.
	// too slow and the black holes at the center collide and hilarity ensues.
	width := 1.0e23
	galaxies := []Galaxy{g0, g1}
	var initialUniverse *Universe
	var numGens int
	time := 2e14
	theta := 0.5

	if mode == "galaxy" {
		numGens = 200000
		initialUniverse = InitializeUniverse(galaxies, width)
	} else if mode == "jupiter" {
		width = 4000000000
		numGens = 1000000
		frequency = 10000
		scalingFactor = 1.0
		time = 1.0
		var s1, s2, s3, s4, s5 *Star
		s1 = &Star{}
		s2 = &Star{}
		s3 = &Star{}
		s4 = &Star{}
		s5 = &Star{}
		initialJupiterGalaxy := InitializeJupiterGalaxy(s1, s2, s3, s4, s5)
		jupiterGalaxy := []Galaxy{initialJupiterGalaxy}
		initialUniverse = InitializeUniverse(jupiterGalaxy, width)
	} else if mode == "collision" {
		initialUniverse = InitializeUniverse(galaxies, width)
		numGens = 200000
		PushGalaxy(g0, -1e3, 1e3)
		PushGalaxy(g1, 1e3, -1e3)
	}
	// now evolve the universe: feel free to adjust the following parameters.

	timePoints := BarnesHut(initialUniverse, numGens, time, theta)

	fmt.Println("Simulation run. Now drawing images.")
	canvasWidth := 1000
	imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

	fmt.Println("Images drawn. Now generating GIF.")
	gifhelper.ImagesToGIF(imageList, mode)
	fmt.Println("GIF drawn.")
}
