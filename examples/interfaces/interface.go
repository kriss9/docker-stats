package main

import "fmt"
import "math"


// INTERFACE
// interface are a method protocoll, and thus models behaviour
// 'interface' is just another type
type Shape interface {
    area() float64
}

func totalArea(shapes ...Shape) float64 {
    var area float64
    for _, s := range shapes {
        area += s.area()
    }
    return area
}

// RECTANGLE
type Rectangle struct {
  w,h float64
}

func (r *Rectangle) area() float64 {
  return r.w * r.h
}

// CIRCLE
type Circle struct {
  r,x,y float64
}

func (c *Circle) area() float64 {
  return math.Pi * c.r * c.r
}

func main() {
  c := Circle{x:0, y:0, r:2}
  r := Rectangle{w:2, h:3}

  // you don't have to pass the address !
  // whatever you pass needs to match the receiver type.
  // so you could simply pass (c,r) and change the receiver types to not be pointer
  // HOWEVER, passing addresses is more efficient by guaranteeing a fixed 4 bytes value passing.
  fmt.Println(totalArea(&c, &r))
}
