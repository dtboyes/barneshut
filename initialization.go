package main

import (
	"math"
	"math/rand"
)

// InitializeUniverse() sets an initial universe given a collection of galaxies and a width.
// It returns a pointer to the resulting universe.
func InitializeUniverse(galaxies []Galaxy, w float64) *Universe {
	var u Universe
	u.width = w
	u.stars = make([]*Star, 0, len(galaxies)*len(galaxies[0]))
	for i := range galaxies {
		for _, b := range galaxies[i] {
			u.stars = append(u.stars, b)
		}
	}
	return &u
}

// InitializeGalaxy takes number of stars in the galaxy, radius of the galaxy to be constructed,
// and center of galaxy to be constructed. Returns a spinning Galaxy object -- which is just a slice of Star pointers
func InitializeGalaxy(numOfStars int, r, x, y float64) Galaxy {
	g := make(Galaxy, numOfStars)

	for i := range g {
		var s Star

		// First choose distance to center of galaxy
		dist := (rand.Float64() + 1.0) / 2.0

		// multiply by factor of r
		dist *= r

		// Next choose the angle in radians to represent the rotation
		angle := rand.Float64() * 2 * math.Pi

		// convert polar coordinates to Cartesian
		s.position.x = x + dist*math.Cos(angle)
		s.position.y = y + dist*math.Sin(angle)

		// set the mass = mass of sun by default
		s.mass = solarMass

		// set the radius equal to radius of sun in m
		s.radius = 696340000

		//set the colors
		s.red = 255
		s.green = 255
		s.blue = 255

		// now spin the galaxy

		// the following is orbital velocity equation
		//dist := Distance(pos, g[i].position)
		speed := 0.5 * math.Sqrt(G*blackHoleMass/dist) // approximation of orbital velocity equation: half of true speed to prevent instability

		s.velocity.x = speed * math.Cos(angle+math.Pi/2.0)
		s.velocity.y = speed * math.Sin(angle+math.Pi/2.0)

		//point g[i] at s
		g[i] = &s

	}

	//add a blackhole to the center of the galaxy

	var blackhole Star
	blackhole.mass = blackHoleMass
	blackhole.position.x = x
	blackhole.position.y = y
	blackhole.blue = 255
	blackhole.radius = 6963400000 // ten times that of a normal star (to make it visible as large)

	g = append(g, &blackhole)

	return g
}

func InitializeJupiterGalaxy(s1, s2, s3, s4, s5 *Star) []*Star {
	s1.red, s1.green, s1.blue = 223, 227, 202
	s2.red, s2.green, s2.blue = 249, 249, 165
	s3.red, s3.green, s3.blue = 132, 83, 52
	s4.red, s4.green, s4.blue = 76, 0, 153
	s5.red, s5.green, s5.blue = 0, 153, 76

	s1.mass = 1.898 * math.Pow(10, 27)
	s2.mass = 8.9319 * math.Pow(10, 22)
	s3.mass = 4.7998 * math.Pow(10, 22)
	s4.mass = 1.4819 * math.Pow(10, 23)
	s5.mass = 1.0759 * math.Pow(10, 23)

	s1.radius = 71000000
	s2.radius = 18210000
	s3.radius = 15690000
	s4.radius = 26310000
	s5.radius = 24100000

	s1.position.x, s1.position.y = 2000000000, 2000000000
	s2.position.x, s2.position.y = 2000000000-421600000, 2000000000
	s3.position.x, s3.position.y = 2000000000, 2000000000+670900000
	s4.position.x, s4.position.y = 2000000000+1070400000, 2000000000
	s5.position.x, s5.position.y = 2000000000, 2000000000-1882700000

	s1.velocity.x, s1.velocity.y = 0, 0
	s2.velocity.x, s2.velocity.y = 0, -17320
	s3.velocity.x, s3.velocity.y = -13740, 0
	s4.velocity.x, s4.velocity.y = 0, 10870
	s5.velocity.x, s5.velocity.y = 8200, 0
	galaxy := []*Star{s1, s2, s3, s4, s5}
	return galaxy
}
