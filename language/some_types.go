package language

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

/**
Also consider a similar thing with a Rectangle
*/
type Rectangle struct {
	x1, y1, x2, y2 float64
}

func distance(x1, y1, x2, y2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(a*a + b*b)
}

/** Method for an area of a rectangle */
func (r *Rectangle) area() float64 {
	l := distance(r.x1, r.y1, r.x1, r.y2)
	w := distance(r.x1, r.y1, r.x2, r.y1)
	return l * w
}

func calcAreaOfRectangle() {

	r := Rectangle{0, 0, 10, 10}
	fmt.Println(r.area())
}

var r Rectangle

/**
This bit covers embedded types
*/

/**
Define a Person struct with a Talk() method
*/
type Person struct {
	Name string
}

func (p *Person) Talk() {
	fmt.Println("Hi, my name is", p.Name)
}

/**
We could now create an Android struct that HAS A Person
*/

type Android struct {
	// HAS A person
	Person Person
	Model  string
}

var a1 Android

func AndriodHASATalk() {
	a1 := new(Android)
	a1.Person.Talk()
}

/**
If we wanted an IS A (inheritance) relationship between the two structs - then use embedded types
*/
type Android2 struct {
	Person
	Model string
}

var a2 Android2

/**
Note that we don't need to reference the 'Person'
The is-a relationship works this way intuitively: People can talk,
an android is a person, therefore an android can talk.
*/

func AndriodISATalk() {
	a2 := new(Android2)
	a2.Talk()
}

/**** Interfaces ****/

/** Interfaces - the same as structs, but contain method sets.A method set is a list of methods
that a type must have in order to “implement” the interface. */
type Shape interface {
	area() float64
}

/* In our case both Rectangle and Circle have area methods which return float64s so
both types implement the Shape interface. By itself this wouldn't be particularly useful,
but we can use interface types as arguments to functions: */

func totalArea(shapes ...Shape) float64 {
	var area float64
	for _, s := range shapes {
		area += s.area()
	}
	return area
}

func PrintTotalArea() {

	fmt.Println(totalArea(&c, &r))
}

/** Interfaces can also be used as fields: */
type MultiShape struct {
	shapes []Shape
}

/** We can even turn MultiShape itself into a Shape by giving it an area method: */
func (m *MultiShape) area() float64 {
	var area float64
	for _, s := range m.shapes {
		area += s.area()
	}
	return area
}
