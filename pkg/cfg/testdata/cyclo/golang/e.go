package _

import "fmt"

func CycloE(a int) {
  switch a {
  case 1:
    fmt.Println("one")
  case 2:
    fmt.Println("two")
  default:
    fmt.Println("whatever")
  }
}
