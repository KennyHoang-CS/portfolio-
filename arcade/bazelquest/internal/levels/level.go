package levels

type Level struct {
    ID           string
    Title        string
    Description  string
    InitialBUILD string
    GoalHint     string
}

var all = []Level{}

func Register(l Level) {
    all = append(all, l)
}

func All() []Level {
    return all
}

func First() Level {
    if len(all) == 0 {
        return Level{}
    }
    return all[0]
}