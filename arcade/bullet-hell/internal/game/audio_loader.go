package game

import (
    "bytes"
    "fmt"
    "io"
    "math/rand"
    "net/http"

    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var audioContext = audio.NewContext(44100)

// Cache of fully decoded PCM players
var bgmCache = map[string]*audio.Player{}

var playlist = []string{
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

// Fully WASM-safe, stutter-proof loader
func loadMP3(path string) (*audio.Player, error) {
    // Fetch file via HTTP (works in WASM + desktop)
    resp, err := http.Get(path)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch %s: %w", path, err)
    }
    defer resp.Body.Close()

    // Read entire MP3 into memory
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read %s: %w", path, err)
    }

    // Decode MP3 into PCM
    decoded, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(data))
    if err != nil {
        return nil, fmt.Errorf("decode failed: %w", err)
    }

    // Read all PCM into memory (critical for WASM performance)
    pcm, err := io.ReadAll(decoded)
    if err != nil {
        return nil, fmt.Errorf("pcm read failed: %w", err)
    }

    // Create player from raw PCM bytes (zero-lag playback)
    player := audioContext.NewPlayerFromBytes(pcm)
    return player, nil
}

// Preload all BGM at startup
func PreloadAllBGM() error {
    // Preload playlist
    for _, path := range playlist {
        p, err := loadMP3(path)
        if err != nil {
            return fmt.Errorf("failed to preload %s: %w", path, err)
        }
        bgmCache[path] = p
    }

    // Preload title + retry
    title, err := loadMP3("assets/ost/nikke_title_ost.mp3")
    if err != nil {
        return err
    }
    bgmCache["title"] = title

    retry, err := loadMP3("assets/ost/nikke_retry_ost.mp3")
    if err != nil {
        return err
    }
    bgmCache["retry"] = retry

    return nil
}

// Accessors
func LoadTitleBGM() *audio.Player {
    return bgmCache["title"]
}

func LoadRetryBGM() *audio.Player {
    return bgmCache["retry"]
}

func LoadRandomBGM() *audio.Player {
    idx := rand.Intn(len(playlist))
    return bgmCache[playlist[idx]]
}

// Stage 1: Only first 3 gameplay tracks (fast startup)
func LoadRandomGameplayStage1() *audio.Player {
    idx := rand.Intn(3) // first 3 tracks
    return bgmCache[playlist[idx]]
}

// Stage 2: Full playlist (after background preload finishes)
func LoadRandomGameplayFull() *audio.Player {
    idx := rand.Intn(len(playlist))
    return bgmCache[playlist[idx]]
}