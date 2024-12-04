package main

import (
	. "bufio"
	. "os"
)

// abuse global variables a bit, because passing arguments = extra tokens
var input []string

// each kernel is separated by 2 newlines
// each line of each kernel is separated by a newline
// some lines have trailing spaces to indicate that nothing needs to be there.
// nvim: use :set list to show trailing spaces
const kernels1 = `XMAS

SAMX

X
M
A
S

S
A
M
X

X   
 M  
  A 
   S

   X
  M 
 A  
S   

   S
  A 
 M  
X   

S   
 A  
  M 
   X`

const kernels2 = `M S
 A 
M S

M M
 A 
S S

S M
 A 
S M

S S
 A 
M M`

func main() {
	scanner := NewScanner(Stdin)

	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	// The search function is defined in another file. In that file I also
	// import "strings" which, if you import it unqualified, conflicts with
	// bufio. This is why I put it in a separate file, so that I can import
	// both bufio and strings unqualified. Saves 1 token
	println(Search(kernels1), Search(kernels2))
}
