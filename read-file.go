package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func extractMountinfo(mifl io.Reader) ([]*MountInfoLine, error) {
	var mi []*MountInfoLine
	sc := bufio.NewScanner(mifl)

	for sc.Scan() {
		line := strings.Split(sc.Text(), " ")
		numOfFields := len(line)
		if numOfFields < 9 {
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
			FileSystemType: line[(numOfFields-10)+7],
			MountSource:    line[(numOfFields-10)+8],
			SuperOptions:   line[(numOfFields-10)+9],
		}

		switch {
		case numOfFields > 10:
			for i := 0; i < (numOfFields - 10); i++ {
				mountinfo.OptionalFields = append(mountinfo.OptionalFields, line[6+i])
			}
		default:
			mountinfo.OptionalFields = append(mountinfo.OptionalFields, "")
		}
		mi = append(mi, mountinfo)
	}
	return mi, nil
}

func generateD3Tree(mi []*MountInfoLine) (*Node, error) {
	const topNodeID = "0"
	var node *Node
	graph := map[string]*Node{}

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

	mifl, err := os.Open(fln)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer mifl.Close()

	mi, err := extractMountinfo(mifl)
	if err != nil {
		log.Fatalf("can't extract mountinfo: %v", err)
	}

	d3Tree, err := generateD3Tree(mi)
	if err != nil {
		log.Fatalf("problem generating D3 tree: %v", err)
	}

	d3json, err := json.Marshal(d3Tree)
	if err != nil {
		log.Fatalf("problem converting to json: %v", err)
	}

	fmt.Println(string(d3json))
}
