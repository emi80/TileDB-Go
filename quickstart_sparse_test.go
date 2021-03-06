/*
Copyright (c) 2018 TileDB, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.


This is a part of the TileDB quickstart tutorial:
https://docs.tiledb.io/en/latest/quickstart.html

When run, this program will create a simple 2D sparse array, write some data
to it, and read a slice of the data back, then clean up.
For simplicity this program does not handle errors
*/
package tiledb_test

import (
	"fmt"
	"os"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
)

// Name of array.
var sparseArrayName = "quickstart_sparse"

func createSparseArray() {
	// Create a TileDB context.
	ctx, _ := tiledb.NewContext(nil)

	// The array will be 4x4 with dimensions "rows" and "cols", with domain [1,4].
	domain, _ := tiledb.NewDomain(ctx)
	rowDim, _ := tiledb.NewDimension(ctx, "rows", []int32{1, 4}, int32(4))
	colDim, _ := tiledb.NewDimension(ctx, "cols", []int32{1, 4}, int32(4))
	domain.AddDimensions(rowDim, colDim)

	// The array will be dense.
	schema, _ := tiledb.NewArraySchema(ctx, tiledb.TILEDB_SPARSE)
	schema.SetDomain(domain)
	schema.SetCellOrder(tiledb.TILEDB_ROW_MAJOR)
	schema.SetTileOrder(tiledb.TILEDB_ROW_MAJOR)

	// Add a single attribute "a" so each (i,j) cell can store an integer.
	a, _ := tiledb.NewAttribute(ctx, "a", tiledb.TILEDB_INT32)
	schema.AddAttributes(a)

	// Create the (empty) array on disk.
	array, _ := tiledb.NewArray(ctx, sparseArrayName)
	array.Create(schema)
}

func writeSparseArray() {
	ctx, _ := tiledb.NewContext(nil)

	// Write some simple data to cells (1, 1), (2, 4) and (2, 3).
	coords := []int32{1, 1, 2, 4, 2, 3}
	data := []int32{1, 2, 3}

	// Open the array for writing and create the query.
	array, _ := tiledb.NewArray(ctx, sparseArrayName)
	array.Open(tiledb.TILEDB_WRITE)
	query, _ := tiledb.NewQuery(ctx, array)
	query.SetLayout(tiledb.TILEDB_UNORDERED)
	query.SetBuffer("a", data)
	query.SetCoordinates(coords)

	// Perform the write and close the array.
	query.Submit()
	array.Close()
}

func readSparseArray() {
	ctx, _ := tiledb.NewContext(nil)

	// Prepare the array for reading
	array, _ := tiledb.NewArray(ctx, sparseArrayName)
	array.Open(tiledb.TILEDB_READ)

	// Slice only rows 1, 2 and cols 2, 3, 4
	subArray := []int32{1, 2, 2, 4}

	// Prepare the vector that will hold the results
	// We take the upper bound on the result size as we do not know how large
	// a buffer is needed since the array is sparse
	maxElements, _ := array.MaxBufferElements(subArray)
	data := make([]int32, maxElements["a"][1])
	coords := make([]int32, maxElements[tiledb.TILEDB_COORDS][1])

	// Prepare the query
	query, _ := tiledb.NewQuery(ctx, array)
	query.SetSubArray(subArray)
	query.SetLayout(tiledb.TILEDB_ROW_MAJOR)
	query.SetBuffer("a", data)
	query.SetCoordinates(coords)

	// Submit the query and close the array.
	query.Submit()
	array.Close()

	// Print out the results.
	for r := 0; r < len(data); r++ {
		i := coords[2*r]
		j := coords[2*r+1]
		fmt.Printf("Cell (%d, %d) has data %d\n", i, j, data[r])
	}
}

// ExampleSparseArray shows and example creation, writing and reading of a
// sparse array
func Example_sparseArray() {
	createSparseArray()
	writeSparseArray()
	readSparseArray()

	// Output: Cell (2, 3) has data 3
	// Cell (2, 4) has data 2

	// Cleanup example so unit tests are clean
	if _, err := os.Stat(sparseArrayName); err == nil {
		os.RemoveAll(sparseArrayName)
	}
}
