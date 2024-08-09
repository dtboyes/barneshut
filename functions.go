package main

import (
	"fmt"
	"math"
)

// BarnesHut is our highest level engine function.
// Input: initial Universe object, a number of generations, an interval time and a parameter theata.
// Output: collection of Universe objects corresponding to updating the system
// over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)
	timePoints[0] = initialUniverse
	//range over the number of generations and set the i-th Universe equal to updating the (i-1)-th Universe
	for i := 1; i < len(timePoints); i++ {
		q := GenerateQuadTree(timePoints[i-1])
		timePoints[i] = UpdateUniverse(timePoints[i-1], q, time, theta)
		fmt.Println("Generation", i)
	}
	return timePoints
}

// UpdateUniverse takes as input a Universe object and a time parameter.
// It returns a new Universe object corresponding to updating the force of gravity on the objects in the given universe, with a time interval of the time parameter in seconds.
func UpdateUniverse(currentUniverse *Universe, q QuadTree, time, theta float64) *Universe {
	//range over the stars in the universe, and update their position/velocity/acceleration
	// made a copy of the universe to avoid modifying the original universe
	newUniverse := &Universe{
		stars: make([]*Star, len(currentUniverse.stars)),
		width: currentUniverse.width,
	}
	// made copies of each star to avoid modifying original stars
	for i := range newUniverse.stars {
		if currentUniverse.stars[i] == nil {
			continue
		}
		newUniverse.stars[i] = &Star{
			position:     currentUniverse.stars[i].position,
			velocity:     currentUniverse.stars[i].velocity,
			acceleration: currentUniverse.stars[i].acceleration,
			mass:         currentUniverse.stars[i].mass,
			radius:       currentUniverse.stars[i].radius,
			red:          currentUniverse.stars[i].red,
			blue:         currentUniverse.stars[i].blue,
			green:        currentUniverse.stars[i].green,
		}

	}
	// update acceleration, velocity, and position of each star
	for i := range newUniverse.stars {
		// don't consider stars that have been set to nil
		if newUniverse.stars[i] == nil {
			continue
		}
		newUniverse.stars[i].acceleration = UpdateAcceleration(q.root, newUniverse.stars[i], theta)
		newUniverse.stars[i].velocity = UpdateVelocity(newUniverse.stars[i], time)
		position := UpdatePosition(newUniverse.stars[i], time)
		// if the position of a star is outside of the bounds of the universe, set it to nil
		if position.x < 0 || position.x > newUniverse.width || position.y < 0 || position.y > newUniverse.width {
			newUniverse.stars[i] = nil
		} else {
			newUniverse.stars[i].position = position
		}
	}

	return newUniverse
}

// GenerateQuadtree takes as input a pointer to a Universe object, which then outputs a QuadTree structure that represents
// the spatial relationship between the different stars that are present in this universe
func GenerateQuadTree(currentUniverse *Universe) QuadTree {
	var q QuadTree
	for _, star := range currentUniverse.stars {
		// create a starting quadrant for the star, marks the beginning of the quadrant calculations, will recursively divide into appropriate quadrant as we traverse down the tree
		quadrant := InitializeStartingQuadrant(currentUniverse)
		// insert the star into the quadtree using a recursive traversal
		q.root = InsertIntoQuadTree(q.root, star, quadrant)
	}
	return q
}

// InitializeStartingQuadrant takes as input a Universe object and returns a Quadrant object that represents the starting quadrant (the entire universe)
func InitializeStartingQuadrant(u *Universe) *Quadrant {
	return &Quadrant{
		x:     0.0,
		y:     0.0,
		width: u.width,
	}
}

// InsertIntoQuadTree takes as input a Node object, a Star object, and a Quadrant object
// it is a recursive implementation of inserting a star into a quadtree
func InsertIntoQuadTree(n *Node, s *Star, quadrant *Quadrant) *Node {
	// if a node is nil, then we make a copy of a new star and assign it to n, additionally assign the quadrant to the node
	if n == nil {
		n = &Node{
			star: &Star{
				position:     s.position,
				velocity:     s.velocity,
				acceleration: s.acceleration,
				mass:         s.mass,
				radius:       s.radius,
				red:          s.red,
				blue:         s.blue,
				green:        s.green,
			},
			sector: *quadrant,
		}
	} else if !IsLeaf(n) { // if we are at an internal node, then we need to determine which quadrant to insert the star into
		if n.star != nil && s != nil { // first check if the star is nil, then check if the stars positions are equal, if they are then merge them
			if n.star.position.x == s.position.x && n.star.position.y == s.position.y {
				n = MergeStars(n, s)
				return n
			}
			// if we do not have the same star, then we want to update the center of mass of the internal node, then determine the next quadrant (and thus the next child node) to traverse to
			n = UpdateCenterOfMass(n, s)
			quadrant_s, index_s := DetermineQuadrant(n, s)
			n.children[index_s] = InsertIntoQuadTree(n.children[index_s], s, quadrant_s)
		}
	} else { // if the current node is a leaf, then we make a new node and call a separate recursive function that handles leaf insertion
		n.children = make([]*Node, 4)
		if n.star.position.x == s.position.x && n.star.position.y == s.position.y { // like before, if the stars are the same, then merge them
			n = MergeStars(n, s)
			return n
		}
		n = InsertAtLeaf(n, s, &n.sector)
	}

	return n
}

// InsertAtLeaf takes as input a Node object, a Star object, and a Quadrant object
// InsertAtLeaf is a recursive implementation of leaf insertion into a quadtree
// it continues to subdivide the leaf node until we reach a state where the leaf node and the star are in different quadrants
func InsertAtLeaf(n *Node, s *Star, quadrant *Quadrant) *Node {
	newNode := &Node{
		children: make([]*Node, 4),
		star: &Star{
			position:     n.star.position,
			velocity:     n.star.velocity,
			acceleration: n.star.acceleration,
			mass:         n.star.mass,
			radius:       n.star.radius,
			red:          n.star.red,
			blue:         n.star.blue,
			green:        n.star.green,
		},
		sector: *quadrant,
	}
	// continue to determine the appropriate quadrants for the current node and star
	quadrant_n, index_n := DetermineQuadrant(newNode, n.star)

	quadrant_s, index_s := DetermineQuadrant(newNode, s)
	if index_n != index_s { // compare the two quadrants, if they are not the same, then the job is simple, just insert into their respective positions
		newNode.children[index_n] = InsertIntoQuadTree(newNode.children[index_n], n.star, quadrant_n)
		newNode.children[index_s] = InsertIntoQuadTree(newNode.children[index_s], s, quadrant_s)

	} else { // if quadrants are the same, repeat step of creating new node, then determining new quadrant using recursion
		newNode.children[index_n] = InsertIntoQuadTree(newNode.children[index_n], n.star, quadrant_n)
		newNode.children[index_s] = InsertAtLeaf(newNode.children[index_s], s, quadrant_s)
	}
	newNode = UpdateCenterOfMass(newNode, s)
	return newNode
}

// DetermineQuadrant takes as input a Node object and a Star object
// it returns a Quadrant object and an integer that represents the index of the child node that the star should be inserted into
// does so by subdividing the region in half and determining which quadrant the star is in
func DetermineQuadrant(n *Node, s *Star) (*Quadrant, int) {
	var quadrant *Quadrant
	var index int = -1
	// if the star is in the left half of the quadrant, then we know that it is in either the NW or SW quadrant
	if s.position.x <= n.sector.x+n.sector.width/2 { // NW or SW Quadrant
		if s.position.y >= n.sector.y+n.sector.width/2 { // NW Quadrant
			quadrant = &Quadrant{
				x:     n.sector.x,
				y:     n.sector.y + n.sector.width/2,
				width: n.sector.width / 2,
			}
			index = 0
		} else { // SW Quadrant
			quadrant = &Quadrant{
				x:     n.sector.x,
				y:     n.sector.y,
				width: n.sector.width / 2,
			}
			index = 2
		}
	} else { // if the star is in the right half of the quadrant, then we know that it is in either the NE or SE quadrant
		if s.position.y >= n.sector.y+n.sector.width/2 { // NE Quadrant
			quadrant = &Quadrant{
				x:     n.sector.x + n.sector.width/2,
				y:     n.sector.y + n.sector.width/2,
				width: n.sector.width / 2,
			}
			index = 1
		} else { // SE Quadrant
			quadrant = &Quadrant{
				x:     n.sector.x + n.sector.width/2,
				y:     n.sector.y,
				width: n.sector.width / 2,
			}
			index = 3
		}
	}
	return quadrant, index
}

// UpdateAcceleration takes as input a Universe object and a Star (in that universe).
// It returns the net acceleration due to the force of gravity of the Star (in components) computed over all other stars in the Universe.
func UpdateAcceleration(n *Node, s *Star, theta float64) OrderedPair {
	var accel OrderedPair
	force := CalculateNetForce(n, s, theta)
	// split acceleration into separate components
	accel.x = force.x / s.mass
	accel.y = force.y / s.mass
	return accel
}

// CalculateNetForce takes as input a Universe object and a Star b.
// It returns the net force (due to gravity) acting on s by all other objects in the given Universe.
// It does so by traversing the quadtree and calculating the force of gravity between the star and the center of mass of the node
func CalculateNetForce(n *Node, s *Star, theta float64) OrderedPair {
	var netForce OrderedPair
	if n != nil { // checks if the star is nil before doing any calculations
		dist := Distance(s.position, n.star.position)
		if dist != 0 { // checks if the stars are the same or not
			if IsLeaf(n) {
				netForce = AddOrderedPairs(netForce, ComputeForce(s, n.star))
			} else {
				if n.sector.width/dist <= theta { // if the width of the quadrant divided by the distance between the two stars is less than the heuristic theta, then we do not need to traverse anymore, our force calculation can be done using the center of mass of the node
					netForce = AddOrderedPairs(netForce, ComputeForce(s, n.star))
				} else { // otherwise, we have to continue traversing through the quadtree to find a node where n.sector.width/dist <= theta
					for i := range n.children {
						netForce = AddOrderedPairs(netForce, CalculateNetForce(n.children[i], s, theta))
					}
				}
			}
		}
	}
	return netForce
}

// ComputeForce takes as input two Star objects b1 and b2 and returns
// an OrderedPair corresponding to the components of a force vector corresponding to the force of gravity of b2 acting on b1.
func ComputeForce(b1, b2 *Star) OrderedPair {
	var force OrderedPair

	//now we do some physics and apply formula
	dist := Distance(b1.position, b2.position)

	//compute magnitude of force
	F := G * b1.mass * b2.mass / (dist * dist)

	//then, split this into components
	deltaX := b2.position.x - b1.position.x
	deltaY := b2.position.y - b1.position.y

	force.x = F * (deltaX / dist)
	force.y = F * (deltaY / dist)

	return force
}

// UpdateVelocity takes as input a Star object and a float time.
// It uses the Newton dynamics equations to compute the updated velocity (in components) of that Star estimated over time seconds.
func UpdateVelocity(s *Star, time float64) OrderedPair {
	var vel OrderedPair
	vel.x = s.velocity.x + s.acceleration.x*time
	vel.y = s.velocity.y + s.acceleration.y*time

	return vel
}

// UpdatePosition takes as input a Star object and a float time.
// It uses the Newton dynamics equations to compute the updated position (in coordinates) of that Star estimated over time seconds.
func UpdatePosition(s *Star, time float64) OrderedPair {
	var pos OrderedPair
	pos.x = s.position.x + s.velocity.x*time + 0.5*s.acceleration.x*time*time
	pos.y = s.position.y + s.velocity.y*time + 0.5*s.acceleration.y*time*time

	return pos
}

// IsLeaf takes as input a Node object and returns a boolean value indicating whether or not the Node is a leaf.
func IsLeaf(n *Node) bool {
	if n == nil { // if the node is nil, then it is a leaf
		return true
	} else if len(n.children) == 0 { // if the node has no children, then it is a leaf
		return true
	} else if n.children[0] == nil && n.children[1] == nil && n.children[2] == nil && n.children[3] == nil { // if the node has children, but they are all nil, then it is a leaf
		return true
	}
	return false
}

// UpdateCenterOfMass takes as input a Node object and a Star object.
func UpdateCenterOfMass(n *Node, s *Star) *Node {
	n.star.position.x = (n.star.position.x*n.star.mass + s.position.x*s.mass) / (n.star.mass + s.mass)
	n.star.position.y = (n.star.position.y*n.star.mass + s.position.y*s.mass) / (n.star.mass + s.mass)
	n.star.mass += s.mass // add the mass of the star we are inserting into the quadtree to the mass of a dummy star whenever it is encountered
	return n
}

// AddOrderedPairs takes two OrderedPair objects and returns the sum of the two OrderedPair objects.
func AddOrderedPairs(p1, p2 OrderedPair) OrderedPair {
	p1.x += p2.x
	p1.y += p2.y
	return p1
}

// Distance takes two position ordered pairs and it returns the distance between these two points in 2-D space.
func Distance(p1, p2 OrderedPair) float64 {
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

// MergeStars takes as input a Node object and a Star object.
// returns a *Node object that contains the star with updated velocity and mass
func MergeStars(n *Node, s *Star) *Node {
	n.star.mass += s.mass
	n.star.velocity.x = (n.star.velocity.x + s.velocity.x) / 2.0
	n.star.velocity.y = (n.star.velocity.y + s.velocity.y) / 2.0
	return n
}

// PushGalaxy takes as input a Galaxy object and two floats x and y.
// It returns a Galaxy object where each star has been pushed using the velocity vector (x, y).
func PushGalaxy(g Galaxy, x, y float64) Galaxy {
	for i := range g {
		g[i].velocity.x += x
		g[i].velocity.y += y
	}
	return g
}
