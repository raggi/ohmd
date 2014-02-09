// Package ohmd implements a wrapper around OpenHMD
// OpenHMD aims to provide a Free and Open Source API and drivers for immersive
// technology, such as head mounted displays with built in head tracking.

package ohmd

/*
#cgo CFLAGS:-I/usr/local/include
#cgo LDFLAGS:-L/usr/local/lib -lopenhmd
#include "openhmd/openhmd.h"
*/
import "C"
import "runtime"

// String values for #Context.ListGets
const (
	VENDOR = C.ohmd_string_value(iota)
	PRODUCT
	PATH
)

// Float values for #Device.Getf
const (
	// float64[4], get - Absolute rotation of the device, in space, as a quaternion.
	ROTATION_QUAT = C.ohmd_float_value(iota + 1)

	// float64[16], get - A "ready to use" OpenGL style 4x4 matrix with a modelview matrix for the left eye of the HMD.
	LEFT_EYE_GL_MODELVIEW_MATRIX

	// float64[16], get - A "ready to use" OpenGL style 4x4 matrix with a modelview matrix for the right eye of the HMD.
	RIGHT_EYE_GL_MODELVIEW_MATRIX

	// float64[16], get - A "ready to use" OpenGL style 4x4 matrix with a projection matrix for the left eye of the HMD.
	LEFT_EYE_GL_PROJECTION_MATRIX
	// float64[16], get - A "ready to use" OpenGL style 4x4 matrix with a projection matrix for the right eye of the HMD.
	RIGHT_EYE_GL_PROJECTION_MATRIX

	// float64[3], get - A 3-D vector representing the absolute position of the device, in space.
	POSITION_VECTOR

	// float64[1], get - Physical width of the device screen, in centimeters.
	SCREEN_HORIZONTAL_SIZE
	// float64[1], get - Physical height of the device screen, in centimeters.
	SCREEN_VERTICAL_SIZE

	// float64[1], get - Physical speration of the device lenses, in centimeters.
	LENS_HORIZONTAL_SEPARATION
	// float64[1], get - Physical vertical position of the lenses, in centimeters.
	LENS_VERTICAL_POSITION

	// float64[1], get - Physical field of view for the left eye, in degrees.
	LEFT_EYE_FOV
	// float64[1], get - Physical display aspect ratio for the left eye screen.
	LEFT_EYE_ASPECT_RATIO
	// float64[1], get - Physical field of view for the left right, in degrees.
	RIGHT_EYE_FOV
	// float64[1], get Physical display aspect ratio for the right eye screen.
	RIGHT_EYE_ASPECT_RATIO

	// float64[1], get/set - Physical interpupilary distance of the user, in centimeters.
	EYE_IPD

	// float64[1], get/set - Z-far value for the projection matrix calculations, i.e. drawing distance.
	PROJECTION_ZFAR
	// float64[1], get/set - Z-near value for the projection matrix calculations, i.e. close clipping distance.
	PROJECTION_ZNEAR

	// float64[6], get - Device specifc distortion value.
	DISTORTION_K
)

var numFloats = map[C.ohmd_float_value]int{
	ROTATION_QUAT:                  4,
	LEFT_EYE_GL_MODELVIEW_MATRIX:   16,
	RIGHT_EYE_GL_MODELVIEW_MATRIX:  16,
	LEFT_EYE_GL_PROJECTION_MATRIX:  16,
	RIGHT_EYE_GL_PROJECTION_MATRIX: 16,
	POSITION_VECTOR:                3,
	SCREEN_HORIZONTAL_SIZE:         1,
	SCREEN_VERTICAL_SIZE:           1,
	LENS_HORIZONTAL_SEPARATION:     1,
	LENS_VERTICAL_POSITION:         1,
	LEFT_EYE_FOV:                   1,
	LEFT_EYE_ASPECT_RATIO:          1,
	RIGHT_EYE_FOV:                  1,
	RIGHT_EYE_ASPECT_RATIO:         1,
	EYE_IPD:                        1,
	PROJECTION_ZFAR:                1,
	PROJECTION_ZNEAR:               1,
	DISTORTION_K:                   6,
}

// int values for #Device.Geti
const (
	// int[1], get Physical horizontal resolution of the device screen.
	SCREEN_HORIZONTAL_RESOLUTION = C.ohmd_int_value(iota)
	// int[1], get Physical vertical resolution of the device screen.
	SCREEN_VERTICAL_RESOLUTION
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Context struct {
	ctx *C.ohmd_context
}

func ctxFinalizer(o *Context) {
	C.ohmd_ctx_destroy(o.ctx)
}

// Create returns a new Context
func Create() *Context {
	c := &Context{}
	c.ctx = C.ohmd_ctx_create()
	runtime.SetFinalizer(c, ctxFinalizer)
	return c
}

// Create for those that prefer common Go conventions
func New() *Context {
	return Create()
}

func (c *Context) destroy() {
	C.ohmd_ctx_destroy(c.ctx)
}

func (c *Context) getError() error {
	if e := C.ohmd_ctx_get_error(c.ctx); e != nil && *e != 0 {
		return Error(C.GoString(e))
	}
	return nil
}

// Update refreshes all values in the context (and devices opened
// from the context). This performs background event pumping.
// Typically users would call this during rendering or animation frames.
func (c *Context) Update() error {
	C.ohmd_ctx_update(c.ctx)
	return c.getError()
}

// Probe searches for devices and returns the number found.
func (c *Context) Probe() (int, error) {
	return int(C.ohmd_ctx_probe(c.ctx)), c.getError()
}

// ListGets fetches device information from the last probe, by idx.
func (c *Context) ListGets(idx int, t C.ohmd_string_value) (string, error) {
	return C.GoString(C.ohmd_list_gets(c.ctx, C.int(idx), t)), c.getError()
}

// ListOpenDevice returns a new device pointer for the device given by idx.
func (c *Context) ListOpenDevice(idx int) (*Device, error) {
	return &Device{c, C.ohmd_list_open_device(c.ctx, C.int(idx))}, c.getError()
}

type Device struct {
	c *Context
	d *C.ohmd_device
}

// Getf queries the device for current floating point values for the given t.
// returns float64 instead of float32 for convenience with the math package.
func (d *Device) Getf(t C.ohmd_float_value) ([]float64, error) {
	n := numFloats[t]
	var f = make([]float32, n, n)

	if C.ohmd_device_getf(d.d, t, (*C.float)(&f[0])) == 0 {
		r := make([]float64, n, n)
		for i := 0; i < n; i++ {
			r[i] = float64(f[i])
		}
		return r, nil
	}
	return nil, d.c.getError()
}

// Setf sets floating point parameters on the device. See the float
// value constants for which fields may be set. The provided slice must
// contain at least the appropriate number of values to avoid segfault.
func (d *Device) Setf(t C.ohmd_float_value, v []float64) error {
	var fs = make([]float32, len(v), len(v))
	for i := 0; i < len(v); i++ {
		fs[i] = float32(v[i])
	}

	if C.ohmd_device_setf(d.d, t, (*C.float)(&fs[0])) == 0 {
		return nil
	}
	return d.c.getError()
}

// Geti queries the device for current integer values for the given t.
// returns a slice of integers. at time of writing this is always a single
// integer.
func (d *Device) Geti(t C.ohmd_int_value) ([]int, error) {
	var r C.int = 0
	if C.ohmd_device_geti(d.d, t, &r) == 0 {
		return []int{int(r)}, nil
	}
	return nil, d.c.getError()
}
