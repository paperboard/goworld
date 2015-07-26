package main

import (
    "fmt"
    "math"
    "math/rand"
    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/geometry"
    "github.com/johanhenriksson/goworld/render"
)

func main() {
    wnd := engine.CreateWindow("voxels", 1280, 800)

    cam := engine.CreateCamera(5,2,5, 1280,800, 65.0, 0.1, 1000.0)

    /* Line material */
    lineProgram := render.CompileVFShader("/assets/shaders/3d_line")
    lineProgram.Use()
    lineProgram.Matrix4f("projection", &cam.Projection[0])
    lineMat := render.CreateMaterial(lineProgram)
    lineMat.AddDescriptor("vertex", gl.FLOAT, 3, 28,  0, false)
    lineMat.AddDescriptor("color",  gl.FLOAT, 4, 28, 12, false)

    /* Tileset Material */
    program := render.CompileVFShader("/assets/shaders/3d_voxel")
    program.Use()
    program.Matrix4f("projection", &cam.Projection[0])
    ttx, _ := render.LoadTexture("/assets/tileset.png")
    tilesetMat := render.CreateMaterial(program)
    tilesetMat.AddDescriptor("vertex", gl.UNSIGNED_BYTE, 3, 8, 0, false)
    tilesetMat.AddDescriptor("normal", gl.BYTE,          3, 8, 3, false)
    tilesetMat.AddDescriptor("tile",   gl.UNSIGNED_BYTE, 2, 8, 6, false)
    tilesetMat.AddTexture(0, ttx)

    tileset := engine.CreateTileset(tilesetMat)

    //gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

    /* Define a voxels */
    grass := &engine.Voxel {
        Xp: tileset.GetId(4, 0),
        Xn: tileset.GetId(4, 0),
        Yp: tileset.GetId(3, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(4, 0),
        Zn: tileset.GetId(4, 0),
    }
    rock := &engine.Voxel {
        Xp: tileset.GetId(2, 0),
        Xn: tileset.GetId(2, 0),
        Yp: tileset.GetId(2, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(2, 0),
        Zn: tileset.GetId(2, 0),
    }

    /* Fill chunk with voxels */
    size := 64
    chk := engine.CreateChunk(size, tileset)
    for z := 0; z < size; z++ {
        for x := 0; x < size; x++ {
            v := rand.Intn(2)
            var vtype *engine.Voxel = nil
            switch v {
            case 0:
                vtype = grass
            case 1:
                vtype = rock
            }
            chk.Set(x,0,z, vtype)
        }
    }

    transf := engine.CreateTransform(0,0,0)

    /* Lines */
    lines := geometry.CreateLines(lineMat)
    /*
    lines.Line(0,3,0, 3,3,0, 1,0,0,1)
    lines.Line(0,3,0, 0,6,0, 0,1,0,1)
    lines.Line(0,3,0, 0,3,3, 0,0,1,1)
    */
    lines.Box(0,0,0,256,256,256,0,1,0,1)
    lines.Box(0,0,0,128,128,128,0,0,1,1)
    lines.Box(0,0,0,64,64,64,1,0,0,1)
    lines.Box(0,0,0,32,32,32,1,1,0,1)
    lines.Box(0,0,0,16,16,16,1,0,1,1)
    lines.Box(0,0,0,8,8,8,0,1,1,1)
    lines.Box(0,0,0,4,4,4,1,0,1,1)
    lines.Box(0,0,0,2,2,2,1,1,1,1)

    lines.Compute()
    lineProgram.Use()
    lineProgram.Matrix4f("model", &transf.Matrix[0])

    /* Compute mesh */
    vmesh := chk.Compute()
    program.Use()
    program.Matrix4f("model", &transf.Matrix[0])
    program.Vec3("lightPos", &mgl.Vec3{ 5,15,-8 })
    program.Float("lightIntensity", 250.0)
    program.Float("ambient", 0.6)

    gl.ClearColor(0,0,0,1)

    /* Render loop */
    wnd.SetRenderCallback(func(wnd *engine.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        program.Use()
        program.Matrix4f("camera", &cam.View[0])
        program.Vec3("cameraPos", &cam.Transform.Position)

        vmesh.Render()

        lineProgram.Use()
        lineProgram.Matrix4f("view", &cam.View[0])
        lines.Render()
    })

    shoot := false
    wnd.SetUpdateCallback(func(dt float32) {
        if engine.KeyDown(engine.KeyF) {
            if !shoot {
                pos := cam.Unproject(1280 / 2, 800 / 2)
                x,y,z := VoxelCoord(cam.Forward, pos)
                fmt.Println(x,y,z)
                chk.Set(x,y,z,grass)
                vmesh = chk.Compute()
                shoot = true
            }
        } else {
            shoot = false
        }
        cam.Update(dt)
    })

    wnd.Loop()
}

func dti(val float32) float32 {
  return float32(math.Abs(float64(val - Round(val))));
}

func Round(f float32) float32 {
    return float32(math.Floor(float64(f + .5)))
}

func VoxelCoord(forward mgl.Vec3, coord mgl.Vec3) (int, int, int) {
    nx := int(coord[0]);
    ny := int(coord[1]);
    nz := int(coord[2]);

    /* find the coordinate that is closer to an integer value */
    /* x < y? */
    if dti(coord[0]) < dti(coord[1]) {
        /* x is less than y */
        /* x < z? */
        if dti(coord[0]) < dti(coord[2]) {
            /* x is closer */
            if forward[0] > 0 {
                /* we are looking to the right */
                fmt.Println("X closest, looking along X+")
                //nx--;
            } else {
                nx++;
                fmt.Println("X closest, looking along X-")
            }
        } else {
            /* z is closer */
            if forward[2] > 0 {
                /* we are looking along z+ */
                fmt.Println("1 Z closest, looking along Z+")
                //nz--
            } else {
                //nz++
                fmt.Println("1 Z closest, looking along Z-")
            }
        }
    } else {
        /* x > y */
        /* is y closer than z? */
        if dti(coord[1]) < dti(coord[2]) {
            /* y is closer! */
            if forward[1] > 0 {
                /* we are looking up */
                fmt.Println("Y closest, looking up")
                ny--
            }
        } else {
            /* z is closer! */
            if forward[2] > 0 {
                /* looking along z+ */
                fmt.Println("1 Z closest, looking along Z+")
                //nz--
            } else {
                //nz++
                fmt.Println("1 Z closest, looking along Z-")
            }
        }
    }
    return nx, ny, nz
}
