package ohmd

import (
	"fmt"
	"os"
	"testing"
)

var e func(...interface{}) = func(args ...interface{}) {
	if e, ok := args[0].(error); ok {
		fmt.Fprintf(os.Stderr, e.Error())
	} else {
		fmt.Fprintf(os.Stderr, "Unknown error: %#v", e)
	}
}

func ExampleGeneralUse() {
	c := New()
	c.destroy()

	c = New()

	if err := c.Update(); err != nil {
		e(err)
	}

	n, err := c.Probe()
	if err != nil {
		e(err)
		return
	}

	fmt.Printf("num devices: %d\n", n)

	for i := 0; i < n; i++ {
		if s, err := c.ListGets(i, VENDOR); err != nil {
			e(err)
		} else {
			fmt.Printf("vendor: %s\n", s)
		}
		if s, err := c.ListGets(i, PRODUCT); err != nil {
			e(err)
		} else {
			fmt.Printf("product: %s\n", s)
		}
		if s, err := c.ListGets(i, PATH); err != nil {
			e(err)
		} else {
			fmt.Printf("path: %s\n", s)
		}
	}

	// default device is normally 0
	d, err := c.ListOpenDevice(0)
	if err != nil {
		e(err)
		return
	}

	if i, err := d.Geti(SCREEN_HORIZONTAL_RESOLUTION); err != nil {
		e(err)
	} else {
		fmt.Printf("horizontal resolution: %d\n", i[0])
	}
	if i, err := d.Geti(SCREEN_VERTICAL_RESOLUTION); err != nil {
		e(err)
	} else {
		fmt.Printf("horizontal resolution: %d\n", i[0])
	}

	if f, err := d.Getf(ROTATION_QUAT); err != nil {
		e(err)
	} else {
		fmt.Printf("rotation quaternion: %v\n", f)
	}
	if f, err := d.Getf(LEFT_EYE_GL_MODELVIEW_MATRIX); err != nil {
		e(err)
	} else {
		fmt.Printf("left eye model matrix: %v\n", f)
	}
	if f, err := d.Getf(RIGHT_EYE_GL_MODELVIEW_MATRIX); err != nil {
		e(err)
	} else {
		fmt.Printf("right eye model matrix: %v\n", f)
	}
	if f, err := d.Getf(LEFT_EYE_GL_PROJECTION_MATRIX); err != nil {
		e(err)
	} else {
		fmt.Printf("left eye gl projection matrix: %v\n", f)
	}
	if f, err := d.Getf(RIGHT_EYE_GL_PROJECTION_MATRIX); err != nil {
		e(err)
	} else {
		fmt.Printf("right eye gl projection matrix: %v\n", f)
	}
	if f, err := d.Getf(POSITION_VECTOR); err != nil {
		e(err)
	} else {
		fmt.Printf("position vector: %v\n", f)
	}
	if f, err := d.Getf(SCREEN_HORIZONTAL_SIZE); err != nil {
		e(err)
	} else {
		fmt.Printf("screen horizontal size: %v\n", f)
	}
	if f, err := d.Getf(SCREEN_VERTICAL_SIZE); err != nil {
		e(err)
	} else {
		fmt.Printf("screen vertical size: %v\n", f)
	}
	if f, err := d.Getf(LENS_HORIZONTAL_SEPARATION); err != nil {
		e(err)
	} else {
		fmt.Printf("lens horizontal separation: %v\n", f)
	}
	if f, err := d.Getf(LENS_VERTICAL_POSITION); err != nil {
		e(err)
	} else {
		fmt.Printf("lens veritcal position: %v\n", f)
	}
	if f, err := d.Getf(LEFT_EYE_FOV); err != nil {
		e(err)
	} else {
		fmt.Printf("left eye fov: %v\n", f)
	}
	if f, err := d.Getf(LEFT_EYE_ASPECT_RATIO); err != nil {
		e(err)
	} else {
		fmt.Printf("left eye aspect ration: %v\n", f)
	}
	if f, err := d.Getf(RIGHT_EYE_FOV); err != nil {
		e(err)
	} else {
		fmt.Printf("right eye fov: %v\n", f)
	}
	if f, err := d.Getf(RIGHT_EYE_ASPECT_RATIO); err != nil {
		e(err)
	} else {
		fmt.Printf("right eye aspect ratio: %v\n", f)
	}
	if f, err := d.Getf(EYE_IPD); err != nil {
		e(err)
	} else {
		fmt.Printf("eye ipd: %v\n", f)
	}
	if f, err := d.Getf(PROJECTION_ZFAR); err != nil {
		e(err)
	} else {
		fmt.Printf("projection zfar: %v\n", f)
	}
	if f, err := d.Getf(PROJECTION_ZNEAR); err != nil {
		e(err)
	} else {
		fmt.Printf("projection znear: %v\n", f)
	}
	if f, err := d.Getf(DISTORTION_K); err != nil {
		e(err)
	} else {
		fmt.Printf("distorion k: %v\n", f)
	}

}

func TestGeneral(t *testing.T) {
	e = t.Error
	ExampleGeneralUse()
}
