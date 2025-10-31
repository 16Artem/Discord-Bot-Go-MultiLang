package locale

import (
    "io/ioutil"
    "path/filepath"
    "sync"
    "gopkg.in/yaml.v3"
)

type Loader struct {
    path  string
    langs map[string]map[string]string
    mu    sync.RWMutex
}

func NewLoader(path string) (*Loader, error) {
    l := &Loader{path: path, langs: make(map[string]map[string]string)}
    if err := l.loadAll(); err != nil {
        return nil, err
    }
    return l, nil
}

func (l *Loader) loadAll() error {
    files, err := ioutil.ReadDir(l.path)
    if err != nil {
        return err
    }
    for _, f := range files {
        if f.IsDir() {
            continue
        }
        if filepath.Ext(f.Name()) != ".yml" && filepath.Ext(f.Name()) != ".yaml" {
            continue
        }
        name := f.Name()
        lang := name[:len(name)-len(filepath.Ext(name))]
        b, err := ioutil.ReadFile(filepath.Join(l.path, f.Name()))
        if err != nil {
            return err
        }
        m := make(map[string]string)
        if err := yaml.Unmarshal(b, &m); err != nil {
            return err
        }
        l.mu.Lock()
        l.langs[lang] = m
        l.mu.Unlock()
    }
    return nil
}

func (l *Loader) T(lang, key string) string {
    l.mu.RLock()
    defer l.mu.RUnlock()
    if m, ok := l.langs[lang]; ok {
        if v, ok2 := m[key]; ok2 {
            return v
        }
    }
    if m, ok := l.langs["en"]; ok {
        if v, ok2 := m[key]; ok2 {
            return v
        }
    }
    return key
}
