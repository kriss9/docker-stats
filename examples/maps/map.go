package main

import "fmt"


func main() {
  m := make(map[string]int)
  m_empty := map[int]float32{}   // same as make

  m["a"], m["b"] = 1,2

  fmt.Println(m)
  fmt.Println(m["a"])
  fmt.Println(m["b"])
  fmt.Println(m_empty)


  // nested maps
  hits := make(map[string]map[string]int)
  hits["one"] = map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
  }
  hits["two"] = map[string]int{
    "a": 21,
    "b": 22,
    "c": 23,
  }

  fmt.Println(hits)
  fmt.Println(hits["one"]["c"])
}
