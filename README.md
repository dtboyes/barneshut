# barneshut
An implementation of a gravity simulator using the Barnes-Hut algorithm in Go

The Barnes-Hut algorithm makes use of a heuristic that avoids explicitly computing the distances between each pair of bodies in a gravity simulator. This heuristic works by using a quadtree datastructure that subdivides the quadrants of the window space, taking note of the bodies within each quadrant. At each point in the quadtree traversal, we compute a statistic that dictates whether we should consider a particular node in the computation of the netforce acting on an object. The drawn space then updates accordingly.

There are 3 separate functionalities.
1. Two galaxies rotating in space
```./BarnesHut galaxy```
2. Two galaxies colliding
```./BarnesHut collision```
3. A simulation of jupiter's gravity
```./BarnesHut jupiter```
