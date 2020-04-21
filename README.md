# Glouglif

Glouglif is what we often call a *proof-of-concept*. It's a project that we do to prove that something work or doesn't with one or more specified technologies.

Here I tried to demonstrate that YES WE CAN USE GOLANG to create a real time 3d application (a game <(^_^)> ! ).

## Build and Run

Before explaining whatever on this cursed land I used and did, just run :

```shell script
go run . SERVER 3000
```

```shell script
go run . CLIENT localhost:3000
```

*maybe you'll have to do some `sudo apt install blah-devel`*

Yes, just specify SERVER or CLIENT because there is a little support for multiplayer.

## Frameworks, Tools and Languages

So to create this simple game, I used :

* golang.org/pkg/image
    * image manipulation
* golang.org/pkg/net
    * network stuff and protocol
* github.com/go-gl/mathgl/mgl32
    * linear algebra (Vec3, Mat4, etc)
* github.com/go-gl/gl
    * OpenGL loader and bindings
* github.com/veandco/go-sdl2
    * input and window system

## Completeness

* [x] First Person Camera
* [x] .obj loader
* [x] Texturing
* [x] Shadows
* [x] Multiplayer
* [x] Extensible
    * [x] Node architecture
    * [x] Controller architecture
* [x] Animated
    * [x] Animation based on pos/rot/scale

## Conclusion

Does it worth ? Does it work ?

### Disadvantages

* the golang standard library (image, net) is kinda slow :/
    * jpeg/png encoding and decoding at least two times slower than libjpeg
    * net has a strange behavior (maybe it's just me) as the ping is slowly increasing over time
* absolutely no object-oriented programming
    * make paradigm like entity-component-system difficult to implement
* golang looks kinda slow on some Windows machine
    * we don't care because gNu/LiNux iS ThE FuTUrE Of GaMinG
* number types
    * `float32(think)`, `int64(other)`, `float64()` are common things in this code

### Advantages

* ez
    * this project has been done in about a week from the ground up
    * about 1.8k loc
    * ultra clean syntax
* binding to c
    * sdl, glfw, sfml, libjpeg, opengl, etc have solid bindings
    * easy to generate other bindings
* dependencies management
    * just run `go run . SERVER 3000` and you're done
    * not the level of npm or cargo but does the job right

### Other criticism

* the golang FAQ itself is saying that the garbage collector is tooooo slow
    > Still, there is room for improvement. The compilers are good but could be better, many libraries need major performance work, and the garbage collector isn't fast enough yet.
    > https://golang.org/doc/faq#Why_does_Go_perform_badly_on_benchmark_x

### Conclusion

Golang is too easy to justify by its syntax the use of a third party language for the implementation of game mechanics (aka scripting). Moreover, the bottleneck produced by Lua&Co on some mainstream engine may make them as fast as a full golang game engine. I insist on the *may*, just to say that the garbage collector, even if it's slow, may not be a worse bottleneck compared to a scripting backend. 

The most successful games all have one thing in common. There are moddable. Players must be able to modify the game and share its modifications. To achieve that with golang, one will implement a lua script system, other will embed the mono runtime.
But one thing is certain, the loss of performance on the "scripting" side will not be fully caught up with by the "native" side. 

## Credits

* some guys on discord who helped a lot
* some guys on github who made awesome repos 
* some guys on the internet who produced high quality tutorials
