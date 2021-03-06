package game

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/render"
)

type PlacementGrid struct {
	ChunkMesh *ChunkMesh
	Color     render.Color

	mesh *geometry.Lines

	/* Current height */
	Y int
}

func NewPlacementGrid(chunk *ChunkMesh) *PlacementGrid {
	pg := &PlacementGrid{
		ChunkMesh: chunk,
		Color:     render.Black,
	}

	// compute grid mesh
	pg.ChunkMesh = chunk
	pg.Y = 9
	pg.Compute()

	return pg
}

func (grid *PlacementGrid) Up() {
	if grid.Y < (grid.ChunkMesh.Sy - 1) {
		fmt.Println("grid up")
		grid.Y++
		grid.Compute()
	}
}

func (grid *PlacementGrid) Down() {
	if grid.Y > 0 {
		fmt.Println("grid down")
		grid.Y--
		grid.Compute()
	}
}

func (grid *PlacementGrid) Update(dt float32) {
	if keys.Pressed(keys.J) {
		grid.Down()
	}
	if keys.Pressed(keys.K) {
		grid.Up()
	}
}

func (grid *PlacementGrid) DrawLines(args engine.DrawArgs) {
	grid.mesh.DrawLines(args)
}

/* Compute grid mesh - draw an empty box for every empty
 * voxel in the current layer */
func (grid *PlacementGrid) Compute() {
	grid.mesh.Clear()

	for x := 0; x < grid.ChunkMesh.Sx; x++ {
		for z := 0; z < grid.ChunkMesh.Sz; z++ {
			if true || grid.ChunkMesh.At(x, grid.Y, z) == EmptyVoxel {
				// place box
				grid.mesh.Box(float32(x), float32(grid.Y)+0.001, float32(z), // position
					1, 0, 1, // size
					grid.Color) // color (RGBA)
			}
		}
	}

	grid.mesh.Compute()
}
