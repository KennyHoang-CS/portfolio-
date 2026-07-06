package workspace

type Workspace struct {
    files map[string]string // path → contents
}

func New() *Workspace {
    return &Workspace{
        files: make(map[string]string),
    }
}