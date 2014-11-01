package main

import "fmt"
import "reflect"

func main() {
  var ia [4]int  // array with 4 integers

  ia[0],ia[1] = 1,2

  fmt.Println(ia)

  str := [2]string{"abc", "def"}
  str2 := [...]string{"this", "is", "a", "variable", "amount", "of", "string"}


  fmt.Println("type", reflect.TypeOf(str))
  fmt.Println("val", str)
  fmt.Println("type", reflect.TypeOf(str2))
  fmt.Println("val", str2)

}
