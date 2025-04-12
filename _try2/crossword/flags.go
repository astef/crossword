package main

import "golang.org/x/exp/constraints"

func Set[T constraints.Integer](b, flags T) T    { return b | flags }
func Clear[T constraints.Integer](b, flags T) T  { return b &^ flags }
func Toggle[T constraints.Integer](b, flags T) T { return b ^ flags }
func Has[T constraints.Integer](b, flags T) bool { return b&flags == flags }
