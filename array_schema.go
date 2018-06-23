package tiledb

/*
#cgo LDFLAGS: -ltiledb
#include <tiledb/tiledb.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// ArraySchema is tiledb array_schema
type ArraySchema struct {
	tiledbArraySchema *C.tiledb_array_schema_t
	context           *Context
}

// NewArraySchema alloc a new array_schema
func NewArraySchema(ctx *Context, arrayType ArrayType) (*ArraySchema, error) {
	arraySchema := ArraySchema{context: ctx}
	ret := C.tiledb_array_schema_alloc(arraySchema.context.tiledbContext, C.tiledb_array_type_t(arrayType), &arraySchema.tiledbArraySchema)
	if ret != C.TILEDB_OK {
		return nil, fmt.Errorf("Error creating tiledb arraySchema: %s", arraySchema.context.GetLastError())
	}

	// Set finalizer for free C pointer on gc
	runtime.SetFinalizer(&arraySchema, func(arraySchema *ArraySchema) {
		arraySchema.Free()
	})

	return &arraySchema, nil
}

// Free tiledb_domain_t that was allocated on heap in c
func (a *ArraySchema) Free() {
	if a.tiledbArraySchema != nil {
		C.tiledb_array_schema_free(&a.tiledbArraySchema)
	}
}

// AddAttributes add one or more attributes
func (a *ArraySchema) AddAttributes(attributes ...Attribute) error {
	for _, attribute := range attributes {
		ret := C.tiledb_array_schema_add_attribute(a.context.tiledbContext, a.tiledbArraySchema, attribute.tiledbAttribute)
		if ret != C.TILEDB_OK {
			return fmt.Errorf("Error adding attributes to tiledb arraySchema: %s", a.context.GetLastError())
		}
	}
	return nil
}

// AttributeNum returns the number of attributes
func (a *ArraySchema) AttributeNum() (uint, error) {
	var attrNum C.uint
	ret := C.tiledb_array_schema_get_attribute_num(a.context.tiledbContext, a.tiledbArraySchema, &attrNum)
	if ret != C.TILEDB_OK {
		return 0, fmt.Errorf("Error getting attribute number for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return uint(attrNum), nil
}

// AttributeFromIndex returns the attribute at a given index
func (a *ArraySchema) AttributeFromIndex(index uint) (*Attribute, error) {
	attr := Attribute{context: a.context}
	ret := C.tiledb_array_schema_get_attribute_from_index(a.context.tiledbContext, a.tiledbArraySchema, C.uint(index), &attr.tiledbAttribute)
	if ret != C.TILEDB_OK {
		return nil, fmt.Errorf("Error getting attribute %d for tiledb arraySchema: %s", index, a.context.GetLastError())
	}
	return &attr, nil
}

// AttributeFromName returns the attribute at a given index
func (a *ArraySchema) AttributeFromName(attrName string) (*Attribute, error) {
	cAttrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(cAttrName))
	attr := Attribute{context: a.context}
	ret := C.tiledb_array_schema_get_attribute_from_name(a.context.tiledbContext, a.tiledbArraySchema, cAttrName, &attr.tiledbAttribute)
	if ret != C.TILEDB_OK {
		return nil, fmt.Errorf("Error getting attribute %s for tiledb arraySchema: %s", attrName, a.context.GetLastError())
	}
	return &attr, nil
}

// SetDomain sets a ArraySchema's domain
func (a *ArraySchema) SetDomain(domain *Domain) error {
	ret := C.tiledb_array_schema_set_domain(a.context.tiledbContext, a.tiledbArraySchema, domain.tiledbDomain)
	if ret != C.TILEDB_OK {
		return fmt.Errorf("Error setting domain for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return nil
}

// Domain returns an ArraySchema's domain
func (a *ArraySchema) Domain() (*Domain, error) {
	domain := Domain{context: a.context}
	ret := C.tiledb_array_schema_get_domain(a.context.tiledbContext, a.tiledbArraySchema, &domain.tiledbDomain)
	if ret != C.TILEDB_OK {
		return nil, fmt.Errorf("Error setting domain for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return &domain, nil
}

// SetCapacity sets an array's capacity
func (a *ArraySchema) SetCapacity(capacity uint64) error {
	ret := C.tiledb_array_schema_set_capacity(a.context.tiledbContext, a.tiledbArraySchema, C.uint64_t(capacity))
	if ret != C.TILEDB_OK {
		return fmt.Errorf("Error setting capacity for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return nil
}

// Capacity gets an array's capacity
func (a *ArraySchema) Capacity() (uint64, error) {
	var capacity C.uint64_t
	ret := C.tiledb_array_schema_get_capacity(a.context.tiledbContext, a.tiledbArraySchema, &capacity)
	if ret != C.TILEDB_OK {
		return 0, fmt.Errorf("Error getting capacity for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return uint64(capacity), nil
}

// SetCellOrder sets an array's cell order
func (a *ArraySchema) SetCellOrder(cellOrder Layout) error {
	ret := C.tiledb_array_schema_set_cell_order(a.context.tiledbContext, a.tiledbArraySchema, C.tiledb_layout_t(cellOrder))
	if ret != C.TILEDB_OK {
		return fmt.Errorf("Error setting cell order for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return nil
}

// CellOrder gets an array's capacity
func (a *ArraySchema) CellOrder() (Layout, error) {
	var cellOrder C.tiledb_layout_t
	ret := C.tiledb_array_schema_get_cell_order(a.context.tiledbContext, a.tiledbArraySchema, &cellOrder)
	if ret != C.TILEDB_OK {
		return -1, fmt.Errorf("Error getting cell order for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return Layout(cellOrder), nil
}

// SetTileOrder sets an array's cell order
func (a *ArraySchema) SetTileOrder(tileOrder Layout) error {
	ret := C.tiledb_array_schema_set_tile_order(a.context.tiledbContext, a.tiledbArraySchema, C.tiledb_layout_t(tileOrder))
	if ret != C.TILEDB_OK {
		return fmt.Errorf("Error setting cell order for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return nil
}

// TileOrder gets an array's capacity
func (a *ArraySchema) TileOrder() (Layout, error) {
	var cellOrder C.tiledb_layout_t
	ret := C.tiledb_array_schema_get_tile_order(a.context.tiledbContext, a.tiledbArraySchema, &cellOrder)
	if ret != C.TILEDB_OK {
		return -1, fmt.Errorf("Error getting cell order for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return Layout(cellOrder), nil
}

// SetCoordsCompressor sets the compressor used for coordinates
func (a *ArraySchema) SetCoordsCompressor(compressor Compressor) error {
	ret := C.tiledb_array_schema_set_coords_compressor(a.context.tiledbContext, a.tiledbArraySchema, C.tiledb_compressor_t(compressor.Compressor), C.int(compressor.Level))
	if ret != C.TILEDB_OK {
		return fmt.Errorf("Error setting coordinates compressor for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return nil
}

// CoordsCompressor gets the compressor used for coordinates
func (a *ArraySchema) CoordsCompressor() (*Compressor, error) {
	var compressorT C.tiledb_compressor_t
	var level C.int
	ret := C.tiledb_array_schema_get_coords_compressor(a.context.tiledbContext, a.tiledbArraySchema, &compressorT, &level)
	if ret != C.TILEDB_OK {
		return nil, fmt.Errorf("Error getting coordinates compressor for tiledb arraySchema: %s", a.context.GetLastError())
	}
	compressor := Compressor{Compressor: CompressorType(compressorT), Level: int(level)}
	return &compressor, nil
}

// SetOffsetsCompressor sets the compressor used for coordinates
func (a *ArraySchema) SetOffsetsCompressor(compressor Compressor) error {
	ret := C.tiledb_array_schema_set_offsets_compressor(a.context.tiledbContext, a.tiledbArraySchema, C.tiledb_compressor_t(compressor.Compressor), C.int(compressor.Level))
	if ret != C.TILEDB_OK {
		return fmt.Errorf("Error setting offsets compressor for tiledb arraySchema: %s", a.context.GetLastError())
	}
	return nil
}

// OffsetsCompressor gets the compressor used for coordinates
func (a *ArraySchema) OffsetsCompressor() (*Compressor, error) {
	var compressorT C.tiledb_compressor_t
	var level C.int
	ret := C.tiledb_array_schema_get_offsets_compressor(a.context.tiledbContext, a.tiledbArraySchema, &compressorT, &level)
	if ret != C.TILEDB_OK {
		return nil, fmt.Errorf("Error getting offsets compressor for tiledb arraySchema: %s", a.context.GetLastError())
	}
	compressor := Compressor{Compressor: CompressorType(compressorT), Level: int(level)}
	return &compressor, nil
}

// Check validates an array schema
func (a *ArraySchema) Check() error {
	ret := C.tiledb_array_schema_check(a.context.tiledbContext, a.tiledbArraySchema)
	if ret != C.TILEDB_OK {
		return fmt.Errorf("Error in checking arraySchema: %s", a.context.GetLastError())
	}
	return nil
}

// LoadArraySchema reads a directory for a ArraySchema
func LoadArraySchema(context *Context, path string) (*ArraySchema, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	a := ArraySchema{context: context}
	ret := C.tiledb_array_schema_load(a.context.tiledbContext, cpath, &a.tiledbArraySchema)
	if ret != C.TILEDB_OK {
		return nil, fmt.Errorf("Error in loading arraySchema from %s: %s", path, a.context.GetLastError())
	}
	return &a, nil
}