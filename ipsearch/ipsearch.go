package main

/*
Demonstrate how to quickly find IP geolocation information for a given IP whose details may be held in a range.
EG

IP "213.105.9.78" is in the range: "213.105.9.0" to "213.105.9.255"

IP location details are typically held in ranges - it uses fewer resources.

*/

import (
	"encoding/binary"
	"fmt"
	"net"
)

type IPLocation struct {
	StartIP net.IP // The start IP address of the range
	EndIP   net.IP // The end IP of the range
	City    string // The city location
	Country string // The country for the IP
	// Add more details if required
}

// The key is the start IP of the range as a uint32. The value is the IP location struct.
type LocCache map[uint32]IPLocation

const (
	MaxUint32 = ^uint32(0)
	Bits      = 32
)

// This is our cache - in this case it's a simple map
var cache LocCache

// This populates the cache, adding two simple IP ranges of 255. These will be used for a cache search of any given IP
func init() {

	cache = make(map[uint32]IPLocation)
	addToCache("213.105.9.0", "213.105.9.255", "London", "United Kingdom")
	addToCache("213.186.33.0", "213.186.33.255", "Roubaix", "France")
	addToCache("213.186.10.0", "213.186.11.128", "Saint Neots", "United Kingdom")

	fmt.Printf("Cache setup  -- \n%+v\n",cache.String())

	// In the real world, this cache will contain many, many entries - maybe even the whole of the internet...
}

// Convenience function that populates the cache
func addToCache(s, e, city, country string) {

	ipl := IPLocation{
		StartIP: net.ParseIP(s),
		EndIP:   net.ParseIP(e),
		City:    city,
		Country: country,
	}

	// convert the IP location start IP to a uint32 cache key.
	ipAsInt := ip2int(ipl.StartIP)
	cache[ipAsInt] = ipl
}

// Convert an IP to an int
func ip2int(ip net.IP) uint32 {

	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// Convert to a string
func (ipl *IPLocation) String() string {

	return fmt.Sprintf("LOC. Start: %s\t(%8x)\tEnd: %s (%8x) -- City %s, Country %s\n", ipl.StartIP.String(), ip2int(ipl.StartIP),
		ipl.EndIP.String(), ip2int(ipl.EndIP), ipl.City, ipl.Country)
}

// A quick and naive way of formatting the cache into string
func (lc LocCache) String() string {

	ret := "Cache Contents\n"
	for key, ipl := range lc {
		ret += fmt.Sprintf("Key: %8x, - value: %s",key, ipl.String())
	}

    return ret
}

// Main func, try a couple of cache checks
func main() {

	fmt.Println("So far so good")

	fmt.Println("\nDemonstrate a cache match")
	ip := "213.105.9.34"
	if result, ok := searchCache(ip); ok {
		fmt.Printf("Found: %s with result: %s\n", ip, result.String())
	} else {
		fmt.Printf("Not found: %s\n", ip)
	}

	fmt.Println("\nDemonstrate a cache miss")
	ip = "1.2.9.54"
	if result, ok := searchCache(ip); ok {
		fmt.Printf("Found: %s with result: %s\n", ip, result.String())
	} else {
		fmt.Printf("Not found: %s\n", ip)
	}

	fmt.Println("\nDemonstrate another cache match - one over the lower boundary")
	ip = "213.186.10.1"
	if result, ok := searchCache(ip); ok {
		fmt.Printf("Found: %s with result: %s\n", ip, result.String())
	} else {
		fmt.Printf("Not found: %s\n", ip)
	}

	fmt.Println("\nDemonstrate another cache match - one under the upper boundary")
	ip = "213.186.11.127"
	if result, ok := searchCache(ip); ok {
		fmt.Printf("Found: %s with result: %s\n", ip, result.String())
	} else {
		fmt.Printf("Not found: %s\n", ip)
	}

	fmt.Println("\nDemonstrate another cache match - one on the lower boundary")
	ip = "213.186.10.0"
	if result, ok := searchCache(ip); ok {
		fmt.Printf("Found: %s with result: %s\n", ip, result.String())
	} else {
		fmt.Printf("Not found: %s\n", ip)
	}
}

// Search the cache. Given an IP as a string, figure out which cache element matches.
func searchCache(ip string) (IPLocation, bool) {

	p := net.ParseIP(ip)
	v := ip2int(p)
	fmt.Printf("Search for %v -- %x\n",p,v)

	mask := MaxUint32 // This will be all bits: 0xffffffff

	for i := 0; i < Bits; i++ {

		ipKey := v & mask // clear the bits in the mask
		fmt.Printf("Cache hit attempt %2d for %8x - ipKey: %8x - mask %x\n",i,v,ipKey,mask)

		if result, ok := cache[ipKey]; ok {
			// If the data in the cache isn't a set of contiguous ranges (this stage could be optional)
			// then check that the IP is the in the range
			endValue := ip2int(result.EndIP)
			if v <= endValue {
				return result, true
			}
		}

		if ipKey == 0 {
			// An optimisation, no need ton continue after this point
			break
		}

		mask = mask << 1 // bit shift and clear the next lower bit
	}

	return IPLocation{}, false
}
