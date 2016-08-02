package main

/**
This takes look at types - and is based upon https://www.golang-book.com/books/intro/9
*/

import (
	"fmt"
	"math"
)

/**
Data cab be encapulated using a struct. For example a Circle
*/
type Circle struct {
	x float64
	y float64
	r float64
}

/**
You can collapse the fields if they're the same type:
*/
type Circle_collapsed struct {
	x, y, r float64
}

/**
Creating circles - example

Like with other data types, this will create a local Circle variable that is by default set to zero.
For a struct zero means each of the fields is set to their corresponding zero
value (0 for ints, 0.0 for floats, "" for strings, nil for pointers, …) We can also use the new function:
*/
var c Circle

/**
This allocates memory for all the fields, sets each of them to their zero value and returns a pointer. (*Circle)
*/
func createCircle() {

	c := new(Circle)

	fmt.Println("printing the circle to stop the not used error", c)
}

/**
Create a circle that's initialised
*/
func createCircle2() {

	// Set x, y, and r
	c := Circle{x: 0, y: 0, r: 5}

	// Leave the field names if you know the order that they were defined
	c = Circle{0, 0, 5}
	fmt.Println("printing the circle to stop the not used error", c)
}

/**
allocate and reference Circle fields
*/
func fieldAccess() {
	// Create an example
	c := Circle{x: 1, y: 2, r: 5}

	c.x = 10
	c.y = 5
	fmt.Println(c.x, c.y, c.r)
}

/** Calc the area of the circle using args passed by value */
func circleAreaByValue(c Circle) float64 {
	return math.Pi * c.r * c.r
}

/** Calc the area of the circle using args passed by reference */
func circleAreaByRef(c *Circle) float64 {
	return math.Pi * c.r * c.r
}

func DemoCalcAreaOfCircle() {

	// Create the ciccle
	c := Circle{0, 0, 5}

	fmt.Println(circleAreaByValue(c))
	fmt.Println(circleAreaByRef(&c))
}

/**
In between the keyword func and the name of the function we've added a “receiver”.
The receiver is like a parameter – it has a name and a type – but by creating the function
in this way it allows us to call the function using the . operator.

The big idea is to link the function to the struct
*/
func (c *Circle) area() float64 {
	return math.Pi * c.r * c.r
}

func DemoCalcAreaOfCircleUsingMethod() {

	// Create the ciccle
	c := Circle{0, 0, 5}

	// accesses c vi the method
	fmt.Println(c.area)
}
