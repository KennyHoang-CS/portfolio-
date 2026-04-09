package engine

type Result struct {
	Success bool
	Logs    []string
}

func NewResult() *Result {
	return &Result{
		Success: true,
		Logs:    []string{},
	}
}

func (r *Result) Log(msg string) {
	r.Logs = append(r.Logs, msg)
}

func (r *Result) Fail(msg string) {
	r.Success = false
	r.Logs = append(r.Logs, msg)
}