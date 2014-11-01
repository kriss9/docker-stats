package main

import "fmt"

type User struct {
  Name string
}

func main() {
  var u int
  u = 42
  fmt.Println("Val = ", u)
  modify(&u)
  fmt.Println("Val = ", u)
}

func modify(u *int)  {
  *u = 1
}



