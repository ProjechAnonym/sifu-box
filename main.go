package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sifu-box/models"

	"gopkg.in/yaml.v3"
)

var workDir string

func init() {
	workDir = getWorkDir()
	file,_ := os.Open(filepath.Join(workDir, "static", "default.template.yaml"))
	defer file.Close()
	content, _ := io.ReadAll(file)
	var template models.Template
	yaml.Unmarshal(content, &template)
	
	
	a,_ := json.MarshalIndent(template, "", "  ")
	fmt.Println(string(a))
}

func main() {

}

func getWorkDir() string {
	// workDir := filepath.Dir(os.Args[0])
	workDir := "E:/Myproject/sifu-box@1.1.0/bin"
	// base_dir := "/root/sifu-clash"
	return filepath.Dir(workDir)
}