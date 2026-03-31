package game

import (
    "bytes"
    "embed"
    "fmt"
    "math/rand"

    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed assets/ost/nikke_title_ost.mp3
var titleOST []byte

//go:embed assets/ost/*.mp3
var playlistFS embed.FS

//go:embed assets/ost/nikke_retry_ost.mp3
var retryOST []byte

var audioContext = audio.NewContext(44100)

var playlist [][]byte

func init() {
    files := []string{
        "assets/ost/encounter.mp3",
        "assets/ost/nikke_gov_afterglow.mp3",
        "assets/ost/stellar_blade_scarlet_theme.mp3",
        "assets/ost/nikke_in_the_mirror.mp3",
        "assets/ost/nikke_providence_remix.mp3",
        "assets/ost/Raven.mp3",
        "assets/ost/Where the Horizon Meets.mp3",
        "assets/ost/ABSOLUTE _ Tactical [GODDESS OF VICTORY _ NIKKE OST].mp3",
        "assets/ost/Emergency Engage [GODDESS OF VICTORY _ NIKKE OST].mp3",
        "assets/ost/The Interceptor.mp3",
        "assets/ost/Unbreakable Sphere __ Endless Blue [GODDESS OF VICTORY _ NIKKE OST].mp3",
        "assets/ost/ZIZ [GODDESS OF VICTORY _ NIKKE OST].mp3",
        "assets/ost/nikke_what_is_luv.mp3",
        "assets/ost/Chapter 42 _ Arkis GODDESS OF VICTORY_ NIKKE OST.mp3",
        "assets/ost/Stellar Blade OST - Eidos 7 Silent Street Combat.mp3",
        "assets/ost/Stellar Blade OST - Everglow.mp3",
        "assets/ost/Stellar Blade OST - Gigas.mp3",
    }

    for _, f := range files {
        data, err := playlistFS.ReadFile(f)
        if err == nil {
            playlist = append(playlist, data)
        }
    }
}

func loadBGM() (*audio.Player, error) {
    d, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(titleOST))
    if err != nil {
        return nil, err
    }

    p, err := audioContext.NewPlayer(d)
    if err != nil {
        return nil, err
    }

    return p, nil
}

func loadRetryBGM() (*audio.Player, error) {
    d, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(retryOST))
    if err != nil {
        return nil, err
    }

    p, err := audioContext.NewPlayer(d)
    if err != nil {
        return nil, err
    }

    return p, nil
}

func loadRandomBGM() (*audio.Player, error) {
    if len(playlist) == 0 {
        return nil, fmt.Errorf("no audio tracks loaded")
    }

    // pick random track
    idx := rand.Intn(len(playlist))
    data := playlist[idx]

    // decode MP3
    d, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(data))
    if err != nil {
        return nil, err
    }

    // create player
    p, err := audioContext.NewPlayer(d)
    if err != nil {
        return nil, err
    }

    return p, nil
}