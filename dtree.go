package dtree

import (
  "path/filepath"
)

const (
  isNode = iota
  isLeaf
)

// Dtree directory tree
struct Dtree {
  Root string
  
}

func (t *Dtree) NewDtree(directory string) *Dtree {
  
}

func (t *Dtree) Insert(path string, value string) {
}

func (t *Dtree) Delete(key string) error {
}

func (t *Dtree) Search(key string) string {
}

func (t *Dtree) GetValues(key string) []string {
}

struct TreeNode {
  Path string
  Name string
  NodeType int32
  Childrens []string
}

func (m *TreeNode) GetChildrens() []string{
  if m != nil {
    return m.Childrens
  }
  return nil
}

func (m *TreeNode) GetPath() string {
  if m != nil && m.Path != nil {
    return *m.Path
  }
  return nil
}

func (m *TreeNode) GetNodeType() int32 {
}
