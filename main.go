package main

import (
	"fmt"
	"os"
	"runtime"
	"math"
	"math/rand"
	"time"
	"strconv"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/rhencke/glut"
)

var (
	winWidth    	= 600
	winHeight   	= 620

	innerPipeSide  	= 58
	pipeSide    	= 60

	flappyMass 		= 1
	score 			= 0
	noOfPipesAdded 	= 0

	pipeVelX 		= float32(-200)
	pipeVelY		= float32(0)

	space       	*chipmunk.Space
	pipe     		[]*chipmunk.Shape
	flappyBirds 	[]*chipmunk.Shape

	birdCollided 	bool
	justStarted		bool
)

type collisionHandlers struct {}

func setWindowHints() {
	// glfw.WindowHint(glfw.Samples, 4)
	// // Create a context to specify the version of OpenGl to 3.3
	// glfw.WindowHint(glfw.ContextVersionMajor, 3)
	// glfw.WindowHint(glfw.ContextVersionMinor, 3)
	// Disable window resize
	glfw.WindowHint(glfw.Resizable, 0)
	// Following two are needed for mac since it by default uses OpenGl 2.2
	// glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile) // remove any deprecated code from older version
	// glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True) 
}

func initOpenGl(window *glfw.Window, w, h int) {
	w, h = window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	//gl.ClearColor(.52, .8, .98, 1)	// lightskyblue
	gl.ClearColor(.25, .88, .83, 1)	// turquoise
}

func initGlut() {
	glut.InitDisplayMode(glut.SINGLE | glut.RGB)
}

// initPhysics sets up the chipmunk space and other physics properties
func initPhysics() {
	space = chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, -900}
}

func drawSquare(colorRed, colorGreen, colorBlue, alpha float32) {
	// first draw the first dark layer
	gl.Color4f(0, 0, 0, 1)
	gl.Begin(gl.POLYGON)	
	gl.Vertex2d(float64(pipeSide/2), float64(pipeSide/2))
	gl.Vertex2d(float64(-pipeSide/2), float64(pipeSide/2))
	gl.Vertex2d(float64(-pipeSide/2), float64(-pipeSide/2))
	gl.Vertex2d(float64(pipeSide/2), float64(-pipeSide/2))	
	gl.End()

	// then draw the actual color layer
	gl.Color4f(colorRed, colorGreen, colorBlue, alpha)
	gl.Begin(gl.POLYGON)	
	gl.Vertex2d(float64(innerPipeSide/2), float64(innerPipeSide/2))
	gl.Vertex2d(float64(-innerPipeSide/2), float64(innerPipeSide/2))
	gl.Vertex2d(float64(-innerPipeSide/2), float64(-innerPipeSide/2))
	gl.Vertex2d(float64(innerPipeSide/2), float64(-innerPipeSide/2))	
	gl.End()

	gl.Vertex3f(0, 0, 0)
}

// add a pipe box to the space
func addOnePipeBox(pos vect.Vect) {
	pipeBox := chipmunk.NewBox(vect.Vector_Zero, vect.Float(pipeSide), vect.Float(pipeSide))
	pipeBox.SetElasticity(0.6)

	body := chipmunk.NewBody(chipmunk.Inf, chipmunk.Inf)
	body.SetPosition(pos)
	body.SetVelocity(pipeVelX, pipeVelY)
	body.IgnoreGravity = true

	body.AddShape(pipeBox)
	space.AddBody(body)
	pipe = append(pipe, pipeBox)
}

// add a row of 7 pipe boxes to the space
func addPipe() {
	// pick a random position for hole in the pipe
	hole := int(math.Floor(rand.Float64() * 6))+ 1
	// add pipe boxes
	for i := 0; i < 9; i++ {
		if (i != hole && i != hole + 1) {
			addOnePipeBox(vect.Vect{vect.Float(winWidth), vect.Float(i * 60 + 30 + i * 10)})
		}		
	}
}

func addFlappy() {
	flappyBird := chipmunk.NewBox(vect.Vector_Zero, vect.Float(pipeSide), vect.Float(pipeSide))
	flappyBird.SetElasticity(0.95)

	body := chipmunk.NewBody(vect.Float(flappyMass), vect.Float(flappyMass))
	body.SetPosition(vect.Vect{100, vect.Float(winHeight)})
	body.SetAngularVelocity(0.2)
	
	// hook collision events
	handlers := collisionHandlers{}
	body.CallbackHandler = handlers

	body.AddShape(flappyBird)
	space.AddBody(body)
	flappyBirds = append(flappyBirds, flappyBird)
}

// renders the display on each update
func render() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()

	// draw pipes
	for _, pipeBox := range pipe {
		gl.PushMatrix()
		pos := pipeBox.Body.Position()
		gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
		drawSquare(.19, .8, .19, 1)	// limegreen
		gl.PopMatrix()
	}

	// draw flappy
	for _, flappyBird := range flappyBirds {
		gl.PushMatrix()
		pos := flappyBird.Body.Position()
		gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
		drawSquare(1, .84, 0, 1)	// gold
		gl.PopMatrix()
	}

	gl.Color4f(1, 0, 1, 1)
	scoreStr := "[ Score: " + strconv.Itoa(score) + " ]"
	// draw score
	drawScore(scoreStr)
}

// step advances the physics engine and cleans up flappy or any pipes that are off-screen
func step(dt float32) {
	space.Step(vect.Float(dt))

	// clean up flappy 
	for i := 0; i < len(flappyBirds); i++ {
		p := flappyBirds[i].Body.Position()
		if p.Y < vect.Float(-pipeSide / 2) || p.Y > vect.Float( winHeight + pipeSide / 2) {
			restartGame()
		}
	}

	// clean up any off-screen pipe
	for i := 0; i < len(pipe); i++ {
		p := pipe[i].Body.Position()
		if p.X < vect.Float(-pipeSide / 2) {
			space.RemoveBody(pipe[i].Body)
			pipe[i] = nil
			pipe = append(pipe[:i], pipe[i+1:]...)
			i-- // consider same index again
		}
	}
}

// output the bitmap onto the screen using glut
func bitmap_output(x, y float32, str string, font glut.BitmapFont) {
	gl.RasterPos2f(x, y)
	for _, ch := range str {
		font.Character(ch)
	}
}

// draw score of the game
func drawScore(score string) {
    bitmap_output(250, 600, score, glut.BITMAP_TIMES_ROMAN_24)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if !glfw.Init() {
		fmt.Fprintf(os.Stderr, "Can't open GLFW")
		return
	}
	defer glfw.Terminate()	

	setWindowHints()

	window, err := glfw.CreateWindow(winWidth, winHeight, "Flappy Bird", nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	window.MakeContextCurrent()

	// set up physics
	initPhysics()
	defer space.Destroy()
	// create the flappy bird
	addFlappy()

	runtime.LockOSThread()
	glfw.SwapInterval(1)

	// set up opengl context
	initOpenGl(window, winWidth, winHeight)

	// init glut
	initGlut()

	// Hook mouse and key events
	window.SetMouseButtonCallback(onMouseBtn)
	window.SetKeyCallback(onKey)
	window.SetCloseCallback(onClose)

	ticksToNextPipe := 10
	ticker := time.NewTicker(time.Second / 60)
	// keep updating till we die ..
	for !window.ShouldClose() {
		// add pipe every 1.5 sec
		ticksToNextPipe--
		if ticksToNextPipe == 0 {
			if !birdCollided {
				ticksToNextPipe = 90
				addPipe()
				noOfPipesAdded++
				// increment score
				if noOfPipesAdded > 2 {
					score++
				}
			} else {
				ticksToNextPipe = 10
			}			
		}
		// render the display
		render()
		step(1.0 / 60.0)
		window.SwapBuffers()
		glfw.PollEvents()

		<-ticker.C // wait up to 1/60th of a second
	}
}

func restartGame() {
	noOfPipesAdded = 0
	score = 0
	birdCollided = false

	cleanFlappy()
	cleanPipes()
	addFlappy()	

	/** TODO
		Just after a new bird takes birth, a collision callback handler
		receives a collision event even when pipes are no where near
		the bird. Set a flag to ignore that event.
	*/
	justStarted = true
}

func jump() {
	for _, flappyBird := range flappyBirds {
		flappyBird.Body.UpdateVelocity(space.Gravity, vect.Float(-.1), vect.Float(-.3))
	}
}

func cleanFlappy() {
	// clean up all pipes
	for i := 0; i < len(flappyBirds); i++ {
		space.RemoveBody(flappyBirds[i].Body)
		flappyBirds[i] = nil
		flappyBirds = append(flappyBirds[:i], flappyBirds[i+1:]...)
		i-- // consider same index again
	}
}

func cleanPipes() {
	// clean up all pipes
	for i := 0; i < len(pipe); i++ {
		space.RemoveBody(pipe[i].Body)
		pipe[i] = nil
		pipe = append(pipe[:i], pipe[i+1:]...)
		i-- // consider same index again
	}
}

func stopPipes() {
	for _, pipeBox := range pipe {
		pipeBox.Body.SetVelocity(0, 0)
	}
}

func sensorizeFlappy() {
	for _, flappyBird := range flappyBirds {
		flappyBird.IsSensor = true
	}
}

func onKey(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
    if action != glfw.Press {
        return
    }

    // disable if event handlers are flagged off
    if birdCollided && (k != glfw.KeyEscape) {
    	return
    }

    switch glfw.Key(k) {
        case glfw.KeyEscape:
            window.SetShouldClose(true)
        case glfw.KeySpace :
            jump()
        default:
            return
    }
}

func onMouseBtn(window *glfw.Window, b glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
        return
    }

    // disable if event handlers are flagged off
    if birdCollided {
    	return
    }

    switch glfw.MouseButton(b) {
        case glfw.MouseButtonLeft :
            jump()
        default:
            return
    }
}

func onClose(window *glfw.Window) {
	window.SetShouldClose(true)
}


func (c collisionHandlers) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	// TODO	investigate the false collision event
	if justStarted {
		justStarted = false
	} else {
		birdCollided = true
		sensorizeFlappy()
		stopPipes()
	}
	
    return true
}

func (c collisionHandlers) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool {
    return true
}

func (c collisionHandlers) CollisionPostSolve(arbiter *chipmunk.Arbiter) {
	return
}

func (c collisionHandlers) CollisionExit(arbiter *chipmunk.Arbiter) {
	return
}