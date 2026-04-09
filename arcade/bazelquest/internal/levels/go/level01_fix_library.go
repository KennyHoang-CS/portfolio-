package golevels

import "github.com/KennyHoang-CS/portfolio/bazelquest/internal/levels"

func init() {
    levels.Register(levels.Level{
        ID:    "go_01_fix_library",
        Title: "Fix the Go Library",
        Description: `A go_library represents a Go package in Bazel.
Every library must list its .go source files in srcs.
This one has none, so Bazel cannot build it.`,
        InitialBUILD: `
go_library(
    name = "util",
    srcs = [],
)
`,
        GoalHint: `Add "util.go" to srcs.`,
    })
}