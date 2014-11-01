package main

import "fmt"
import "math"


type Circle struct {
  r,x,y float64
}

func (c *Circle) area() float64 {
  return math.Pi * c.r * c.r
}

func main() {
 c1 := Circle{x:0, y:0, r:2}
 cp := new(Circle)
 c2 := Circle{x: 0, y: 0, r: 5}
 var c *Circle = cp
 
 fmt.Println("c1:", c1)
 fmt.Println("cp:", cp)
 fmt.Println("c2:", c2)
 fmt.Println("c:", c)
 c.x = 1
 fmt.Println("c.x = 1:", c)
 fmt.Println("cp.x = 1:", cp)

 fmt.Println("Circle:", c1)
 fmt.Println("Area:", c1.area())
}
