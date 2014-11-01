package inheritance

import "fmt"

type Person struct {
  Name string
}

// method
// receiver : (p *Person)
func (p *Person) Talk() {
	fmt.Println("Hi, my name is", p.Name)
}

// is-a relationship
type Android struct {
	Person   // an Android is a Person
  Model string
}

// func main() {
//  a := new(Android)
//  a.name = "Android 1"
//  a.talk()
// }
