package main

import (
    "fmt"
    "github.com/johanhenriksson/goworld/game"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"

    "github.com/ianremmler/ode"
    mgl "github.com/go-gl/mathgl/mgl32"

    opensimplex "github.com/ojrac/opensimplex-go"
)

const (
    WIDTH = 1280
    HEIGHT = 800
)

func fromOdeVec3(vec ode.Vector3) mgl.Vec3 {
    return mgl.Vec3 {
        float32(vec[0]),
        float32(vec[1]),
        float32(vec[2]),
    }
}

func toOdeVec3(vec mgl.Vec3) ode.Vector3 {
    return ode.Vector3 {
        float64(vec[0]),
        float64(vec[1]),
        float64(vec[2]),
    }
}

func main() {
    app := engine.NewApplication("voxels", WIDTH, HEIGHT)

    geomPass := app.Render.Passes[0].(*engine.GeometryPass)

    /* create a camera */
    app.Scene.Camera = engine.CreateCamera(-3,10,-3, WIDTH, HEIGHT, 65.0, 0.1, 500.0)
    app.Scene.Camera.Transform.Rotation[1] = 130.0

    /* test voxel chunk */
    tileset := game.CreateTileset()
    chk := generateChunk(1, tileset)
    chk.Compute()
    geomPass.Material.SetupVertexPointers()

    obj := engine.NewObject(0,0,0)
    obj.Attach(chk)

    /* attach to scene */
    app.Scene.Add(obj)

    // physics
    ode.Init(0, ode.AllAFlag)
    side := 0.5
    world := ode.NewWorld()
    space := ode.NilSpace().NewHashSpace()
    box1 := world.NewBody()
    box1.SetPosition(ode.V3(0, 20, 0))
    mass := ode.NewMass()
    mass.SetBox(1, ode.V3(side, side, side))
    mass.Adjust(1)
    box1.SetMass(mass)
    box1_col := space.NewBox(ode.V3(side, side, side))
    box1_col.SetBody(box1)
    ctGrp := ode.NewJointGroup(1000)

    world.SetGravity(ode.V3(0,-0.1,0))
    space.NewPlane(ode.V4(0,1,0,0))

    cam_ray := space.NewRay(10)


    // buffer display window
    bufferWindow := func(title string, texture *render.Texture, x, y float32) {
        win_color := render.Color{0.15, 0.15, 0.15, 0.8}
        text_color := render.Color{1, 1, 1, 1}

        win   := app.UI.NewRect(win_color, x, y, 250, 280, -10)
        label := app.UI.NewText(title, text_color, 0, 0, -21)
        img   := app.UI.NewImage(texture, 0, 30, 250, 250, -20)
        img.Quad.FlipY()

        win.Append(img)
        win.Append(label)

        /* attach UI element */
        app.UI.Append(win)
    }

    bufferWindow("Diffuse", geomPass.Buffer.Diffuse, 30, 30)
    bufferWindow("Normal", geomPass.Buffer.Normal, 30, 340)

    /* Render loop */
    app.Window.SetRenderCallback(func(wnd *engine.Window, dt float32) {
        /* render scene */
        app.Render.Draw()

        /* draw user interface */
        app.UI.Draw()

        // update position
        fmt.Println(box1.Position())
        obj.Transform.Position = fromOdeVec3(box1.Position())

        cam_ray.SetPosDir(toOdeVec3(app.Scene.Camera.Position), toOdeVec3(app.Scene.Camera.Forward))

        space.Collide(0, func(data interface{}, obj1, obj2 ode.Geom) {
            contact := ode.NewContact()
            body1, body2 := obj1.Body(), obj2.Body()
            if body1 != 0 && body2 != 0 && body1.Connected(body2) {
                return
            }
            contact.Surface.Mode = 0
            contact.Surface.Mu = 0.1
            contact.Surface.Mu2 = 0
            cts := obj1.Collide(obj2, 1, 0)
            if len(cts) > 0 {
                if obj1 == cam_ray || obj2 == cam_ray {
                    fmt.Println("ray collision")
                    return // dont attach anything
                }

                contact.Geom = cts[0]
                ct := world.NewContactJoint(ctGrp, contact)
                ct.Attach(body1, body2)
            }
        })
        world.QuickStep(0.05)
        ctGrp.Empty()
    })

    app.Run()
}

func generateChunk(size int, tileset *game.Tileset) *game.Chunk {
    /* Define voxels */
    grass := &game.Voxel{
        Xp: tileset.GetId(4, 0),
        Xn: tileset.GetId(4, 0),
        Yp: tileset.GetId(3, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(4, 0),
        Zn: tileset.GetId(4, 0),
    }
    rock := &game.Voxel{
        Xp: tileset.GetId(2, 0),
        Xn: tileset.GetId(2, 0),
        Yp: tileset.GetId(2, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(2, 0),
        Zn: tileset.GetId(2, 0),
    }

    /* Fill chunk with voxels */
    f := 1.0 / 5
    chk := game.CreateChunk(size, tileset)
    simplex := opensimplex.NewWithSeed(1000)
    for z := 0; z < size; z++ {
        for y := 0; y < size; y++ {
            for x := 0; x < size; x++ {
                fx, fy, fz := float64(x) * f, float64(y) * f, float64(z) * f
                v := simplex.Eval3(fx, fy, fz)
                var vtype *game.Voxel = nil
                if y <= size / 2 {
                    vtype = grass
                }
                if v > 0.0 {
                    vtype = rock
                }
                chk.Set(x, y, z, vtype)
            }
        }
    }

    return chk
}