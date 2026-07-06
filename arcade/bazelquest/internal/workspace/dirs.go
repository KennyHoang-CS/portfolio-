package workspace

import "strings"

func (w *Workspace) ListDir(dir string) []string {
    prefix := dir
    if prefix != "" && prefix[len(prefix)-1] != '/' {
        prefix += "/"
    }

    out := []string{}
    for p := range w.files {
        if strings.HasPrefix(p, prefix) {
            rest := strings.TrimPrefix(p, prefix)
            if !strings.Contains(rest, "/") {
                out = append(out, rest)
            }
        }
    }
    return out
}