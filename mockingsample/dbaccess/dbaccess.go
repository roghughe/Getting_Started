/*
 The db access library is the code that you've written to both wrap access your database by whatever driver
 you're using. This may be a Postgres, MySQL or whatever. The reason it's here is to ensure that the thre party driver
 code doesn't leak out over all your other production code.

 In this example - it does nothing...
 */
package dbaccess

import "fmt"

// This is the real DB Connection details.
type RealConnection struct {

	host string // The DB host to connect to
	port int32  // The port
	URL  string // DB access URL / connection string

	// other stuff may go here
}



// This is any old DB read func, you'll have lots of these in your application - probably.
func (r *RealConnection) ReadSomething(arg0, arg1 string) ([]string, error)  {
	fmt.Printf("This is the real database driver - read args: %s -- %s",arg0,arg1)
	return []string{}, nil
}

// This is any old DB insert/ update function.
func (r * RealConnection) WriteSomething(arg0 []string) error {
	fmt.Printf("This is the real database driver - write args: %v",arg0)
	return nil
}


// Group together the methods in one or more interfaces, gathering them along functionality lines and keeping the interface small.
type SomeFunctionalityGroup interface {

	ReadSomething(arg0, arg1 string) ([]string, error)
	WriteSomething(arg0 []string) error
}

/*
There will probably be more functions and maybe other interfaces from here on in.
 */

