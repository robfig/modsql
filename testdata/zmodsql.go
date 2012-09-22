// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

package main

type types struct {
	t_int     int
	t_int8    int8
	t_int16   int16
	t_int32   int32
	t_int64   int64
	t_float32 float32
	t_float64 float64
	t_string  string
	t_binary  []byte
	t_byte    byte
	t_rune    rune
	t_bool    bool
}

type default_value struct {
	id        int
	d_int8    int8
	d_float32 float32
	d_string  string
	d_binary  []byte
	d_byte    byte
	d_rune    rune
	d_bool    bool
	d_findex  int
}

type times struct {
	typeId     int
	t_duration time.Duration
	t_datetime time.Time
}