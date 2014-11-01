package main

import "fmt"
import "reflect"

func main() {
  var ia []int  // slice with 4 integers

  ia = make([]int, 10)
  ia[0] = 1
  ia = append(ia, 2,3,4)

  fmt.Println("type", reflect.TypeOf(ia))
  fmt.Println("val_ia", ia)

  for _, v:= range ia {
    if (v % 2) != 0 {
      fmt.Println("val_v", v)
    }
  }
}
