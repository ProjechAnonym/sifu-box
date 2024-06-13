package main

import "fmt"
func main() {
	err := Set_value(map[string]interface{}{"b":"aa"},"a")
	fmt.Println(err)
	Test()
	err = Set_value("cc","a")
	fmt.Println(err)
	Test()
}