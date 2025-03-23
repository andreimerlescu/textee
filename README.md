# Textee

This project provides a Go package designed to analyze and categorize strings by creating substring groupings and 
quantifying the occurrences of these substrings. Its primary purpose is to offer a robust and efficient means of 
textual analysis, useful in various applications like natural language processing, data mining, and pattern recognition. 
For instance, it can help identify common phrases within large texts, analyze code for repeated patterns, or even assist 
in cryptographic analysis.

A Go developer might use this package to leverage its performance and concurrency capabilities, inherent in Go, to 
handle large datasets or real-time processing with minimal latency. It's particularly beneficial in scenarios requiring 
quick, repetitive analysis of text, such as processing logs, examining user input, or preparing data for machine 
learning models. By simplifying complex tasks into a manageable and efficient process, this package aims to save time 
and resources while providing accurate and valuable insights into text data.

## Results

This package is designed to accomplish the following type of string manipulation: 

### Input 

Some text will be provided into the `.ParseString()` method. For example: 

```txt
All right let's move from this point on 16 March 84, let's move in time to our second location which is a specific building near where you are now. Are you ready? Just a minute. All. right, I will wait. All right, move now from this area to the front ground level of the building known as the Menara Building, to the front of, on the ground, the Menara Building.
```

### Output

The output of the `.ParseString()` method is `*Textee` where the `.Substrings` property is a `map[string]*atomic.Int32` 
which holds the substring and its quantity of occurrences in the string. The larger the number, the more occurrences
of the matched substring. All strings are converted to lowercase to prevent case-sensitive data duplication.

```txt
all: 1
all right: 1
all right lets: 1
right: 1
right lets: 1
right lets move: 1
lets: 1
lets move: 1
lets move from: 1
move: 1
move from: 1
move from this: 1
from: 1
from this: 1
from this point: 1
this: 1
this point: 1
this point on: 1
point: 1
point on: 1
point on 16: 1
on: 1
on 16: 1
on 16 March: 1
16: 1
16 March: 1
16 March 84: 1
March: 1
March 84: 1
March 84 lets: 1
84: 1
84 lets: 1
84 lets move: 1
lets: 2
lets move: 2
lets move in: 1
move: 2
move in: 1
move in time: 1
in: 1
in time: 1
in time to: 1
time: 1
time to: 1
time to our: 1
...
```

The `.SortedSubstrings()` function of the `*Textee` package / `ITextee` interface returns a slice of a struct called: 

```go
type SubstringQuantity struct {
	Substring string
	Quantity  int
}
```

This function resolves the `*atomic.Int32` by storing its basic `int` value at `.Load()` when its executed and the
final results are then sorted using the [sort.Interface](https://pkg.go.dev/sort). 

## Usage

You'll need to include this project in your workspace.

```bash
go get -u github.com/andreimerlescu/textee
```

## Dependencies

This project depends only on the [go-gematria](https://github.com/andreimerlescu/go-gematria) and 
[go-sema](github.com/andreimerlescu/go-sema) modules. They are both by the same author and are both
licensed under Apache 2.0.

## Example

```go
package main

import (
    "fmt"
    "os"
    "errors"

    "github.com/andreimerlescu/textee"
)

const inputString = `All right let's move from this point on 16 March 84, let's move in time to our second location which is a specific building near where you are now. Are you ready? Just a minute. All. right, I will wait. All right, move now from this area to the front ground level of the building known as the Menara Building, to the front of, on the ground, the Menara Building.`

func main() {
    // just calculate the substrings (no Gematria)
    tt1, err := textee.NewTextee(inputString)
    if err != nil {
        _,_ = fmt.Fprintf(os.Stderr, "%v", errors.Join(textee.ErrBadParsing, err))
    }

    // calculate a new substring with Gematria 
    var err2 error
    tt1, err2 = tt1.CalculateGematria()
    if err2 != nil {
        _,_ = fmt.Fprintf(os.Stderr, "%v", errors.Join(textee.ErrBadParsing, err2))
    }
    for substring, quantity := range tt1.Substrings {
        fmt.Printf("substring '%v' has %d occurrences\n", substring, quantity.Load())
    }

    // combine them together
    tt2, err3 := textee.NewTextee(inputString)
    if err3 != nil {
        _,_ = fmt.Fprintf(os.Stderr, "%v", errors.Join(textee.ErrBadParsing, err3))
    }
    var err4 error
    tt2, err4 = tt2.CalculateGematria()
    if err4 != nil {
        _,_ = fmt.Fprintf(os.Stderr, "%v", errors.Join(textee.ErrBadParsing, err4))
    }
    fmt.Println(tt2)

    // sort the substrings by the quantity of occurrences in the original string, most common are first
    sortedSubstrings := tt2.SortedSubstrings()
    for idx, substring := range sortedSubstrings {
        fmt.Printf("%d: substring '%v' has %d occurrences\n", idx, substring.Substring, substring.Quantity)
    }
}

```

## License

This project is Open Source under the Apache 2.0 license. Feel free to use it where you see a need for thing kind of 
string manipulation.
