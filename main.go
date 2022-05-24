package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type codeWorkSpace struct {
	Folders []folder `json:"folders"`
}

type folder struct {
	Path string `json:"path"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("empty dst")
		os.Exit(0)
	}
	for _, v := range os.Args[1:] {
		items, err := os.ReadDir(v)
		if err != nil {
			fmt.Printf("os.ReadDir: %s\n", err)
			os.Exit(0)
		}
		cws := &codeWorkSpace{
			Folders: make([]folder, 0),
		}
		cwsName := v + ".code-workspace"
		m := map[string]struct{}{}
		if FileExists(cwsName) {
			data, err := os.ReadFile(cwsName)
			if err != nil {
				fmt.Printf("os.ReadFile: %s\n", err)
				os.Exit(0)
			}
			if err := json.Unmarshal(data, cws); err != nil {
				fmt.Printf("json.Unmarshal: %s\n", err)
				os.Exit(0)
			}
			for _, v := range cws.Folders {
				m[v.Path] = struct{}{}
			}
		}
		for i := range items {
			if !items[i].IsDir() {
				continue
			}
			path := filepath.Join(v, items[i].Name())
			if _, ok := m[path]; ok {
				continue
			}
			cws.Folders = append(cws.Folders, folder{Path: path})
		}
		data, err := json.MarshalIndent(cws, "", "    ")
		if err != nil {
			fmt.Printf("json.MarshalIndent: %s\n", err)
			os.Exit(0)
		}
		if err := os.WriteFile(cwsName, data, 0644); err != nil {
			fmt.Printf("os.WriteFile: %s\n", err)
			os.Exit(0)
		}
		fmt.Printf("code-workspace: %s sycned\n", cwsName)
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}