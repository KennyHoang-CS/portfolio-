package game

import (
    "bytes"
    "io"
    "net/http"

    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type PreloadManager struct {
    stage1 []string
    stage2 []string
    index  int
    stage  int // 1 or 2
    done1  bool
    done2  bool
    ctx    *audio.Context
}

func NewPreloadManager() *PreloadManager {
    // First 3 gameplay tracks
    stage1Gameplay := playlist[:3]

    return &PreloadManager{
        stage1: append([]string{
            "assets/ost/nikke_title_ost.mp3",
            "assets/ost/nikke_retry_ost.mp3",
        }, stage1Gameplay...),

        stage2: playlist[3:], // remaining gameplay tracks

        index: 0,
        stage: 1,
        ctx:   audioContext,
    }
}

func (pm *PreloadManager) Update() {
    if pm.done2 {
        return
    }

    var list []string
    if pm.stage == 1 {
        list = pm.stage1
    } else {
        list = pm.stage2
    }

    if pm.index >= len(list) {
        if pm.stage == 1 {
            pm.done1 = true
            pm.stage = 2
            pm.index = 0
            return
        }
        pm.done2 = true
        return
    }

    path := list[pm.index]

    resp, err := http.Get(path)
    if err != nil {
        pm.index++
        return
    }
    data, _ := io.ReadAll(resp.Body)
    resp.Body.Close()

    decoded, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(data))
    if err != nil {
        pm.index++
        return
    }

    pcm, _ := io.ReadAll(decoded)
    player := pm.ctx.NewPlayerFromBytes(pcm)

    bgmCache[path] = player

    pm.index++
}

func (pm *PreloadManager) Stage1Done() bool { return pm.done1 }
func (pm *PreloadManager) Stage2Done() bool { return pm.done2 }

func (pm *PreloadManager) Progress() float64 {
    if pm.stage == 1 {
        return float64(pm.index) / float64(len(pm.stage1))
    }
    return 1.0
}