package game

type World struct {
    entities []*Entity
}

func NewWorld() *World {
    return &World{
        entities: []*Entity{},
    }
}

func (w *World) NewEntity() *Entity {
    e := &Entity{}
    w.entities = append(w.entities, e)
    return e
}

func (w *World) Entities() []*Entity {
    return w.entities
}

func (w *World) Cleanup() {
    alive := w.entities[:0]
    for _, e := range w.entities {
        if !e.Destroy {
            alive = append(alive, e)
        }
    }
    w.entities = alive
}