package workspace

func (w *Workspace) Exists(path string) bool {
    _, ok := w.files[path]
    return ok
}

func (w *Workspace) Read(path string) (string, bool) {
    c, ok := w.files[path]
    return c, ok
}

func (w *Workspace) Write(path, contents string) {
    w.files[path] = contents
}

func (w *Workspace) Delete(path string) {
    delete(w.files, path)
}

func (w *Workspace) List() []string {
    out := make([]string, 0, len(w.files))
    for p := range w.files {
        out = append(out, p)
    }
    return out
}