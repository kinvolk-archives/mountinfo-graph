package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Node struct {
	*MountInfoLine
	Children []*Node `json:"children,omitempty"`
}

type MountInfoLine struct {
	MountID        string   `json:"-"`
	ParentID       string   `json:"-"`
	MajorMinor     string   `json:"-"`
	Root           string   `json:"-"`
	MountPoint     string   `json:"name"`
	MountOptions   string   `json:"-"`
	OptionalFields []string `json:"-"`
	FileSystemType string   `json:"-"`
	MountSource    string   `json:"-"`
	SuperOptions   string   `json:"-"`
}

func extractMountinfo(flp string) ([]*MountInfoLine, error) {
	var mi []*MountInfoLine

	mifl, err := os.Open(flp)
	if err != nil {
		return nil, err
	}
	defer mifl.Close()

	sc := bufio.NewScanner(mifl)
	for sc.Scan() {
		line := strings.Split(sc.Text(), " ")
		lineLen := len(line)
		if lineLen < 9 {
			return nil, fmt.Errorf("not enough fields in the mountinfo file: %v", line)
		}
		mountinfo := &MountInfoLine{
			// TODO: find a cleaner way of doing this
			MountID:        line[0],
			ParentID:       line[1],
			MajorMinor:     line[2],
			Root:           line[3],
			MountPoint:     line[4],
			MountOptions:   line[5],
			FileSystemType: line[(lineLen-10)+7],
			MountSource:    line[(lineLen-10)+8],
			SuperOptions:   line[(lineLen-10)+9],
		}

		switch {
		case lineLen > 10:
			for i := 0; i < (lineLen - 10); i++ {
				mountinfo.OptionalFields = append(mountinfo.OptionalFields, line[6+i])
			}
		default:
			mountinfo.OptionalFields = append(mountinfo.OptionalFields, "")
		}
		mi = append(mi, mountinfo)
	}
	return mi, nil
}

func generateD3Tree(fln string) (*Node, error) {
	const topNodeID = "0"
	var node *Node
	graph := map[string]*Node{}
	mi, err := extractMountinfo(fln)
	if err != nil {
		return nil, fmt.Errorf("can't extract mountinfo: %v", err)
	}

	for _, mountinfo := range mi {
		graph[mountinfo.MountID] = &Node{
			MountInfoLine: mountinfo,
		}
	}
	for _, n := range graph {
		if n.ParentID == topNodeID {
			node = n
			continue
		}
		graph[n.ParentID].Children = append(graph[n.ParentID].Children, n)
	}
	return node, nil
}

func main() {
	fln := "mi"

	d3Tree, err := generateD3Tree(fln)
	if err != nil {
		log.Fatalf("problem generating D3 tree: %v", err)
	}

	d3json, err := json.Marshal(d3Tree)
	if err != nil {
		log.Fatalf("problem converting to json: %v", err)
	}

	fmt.Println(string(d3json))
}
