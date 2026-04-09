package gorules

import "errors"

var (
	ErrToolchainMissing   = errors.New("Go toolchain not installed")
	ErrMissingSrcs        = errors.New("go_library requires non-empty srcs")
	ErrMissingBinarySrcs  = errors.New("go_binary requires non-empty srcs")
	ErrMissingLibraryDeps = errors.New("go_binary requires at least one go_library dep")
)