package gorules

type GoSourceProvider struct {
	Srcs []string
}

func (p *GoSourceProvider) Kind() string { return "go_source" }

type GoCompileProvider struct {
	PackageName string
	Objects     []string
}

func (p *GoCompileProvider) Kind() string { return "go_compile" }

type GoBinaryProvider struct {
	Name   string
	Binary string
}

func (p *GoBinaryProvider) Kind() string { return "go_binary" }