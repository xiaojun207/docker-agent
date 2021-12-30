package agent

import (
	"github.com/docker/docker/api/types"
	"sort"
	"strings"
)

type PortSlice struct {
	Arr []types.Port
}

func (e PortSlice) Len() int      { return len(e.Arr) }
func (e PortSlice) Swap(i, j int) { e.Arr[i], e.Arr[j] = e.Arr[j], e.Arr[i] }
func (e PortSlice) Less(i, j int) bool {
	s := e.Arr
	return s[i].PublicPort < s[j].PublicPort
	/**
	||
		s[i].PrivatePort < s[j].PrivatePort ||
		strings.Compare(s[i].IP, s[j].IP) < 0 ||
		strings.Compare(s[i].Type, s[j].Type) <0
	*/
}

type MountSlice struct {
	Arr []types.MountPoint
}

func (e MountSlice) Len() int      { return len(e.Arr) }
func (e MountSlice) Swap(i, j int) { e.Arr[i], e.Arr[j] = e.Arr[j], e.Arr[i] }
func (e MountSlice) Less(i, j int) bool {
	s := e.Arr
	return strings.Compare(s[i].Source, s[j].Source) < 0
}

type ImageSlice struct {
	Arr []types.ImageSummary
}

func (e ImageSlice) Len() int      { return len(e.Arr) }
func (e ImageSlice) Swap(i, j int) { e.Arr[i], e.Arr[j] = e.Arr[j], e.Arr[i] }
func (e ImageSlice) Sort() []types.ImageSummary {
	sort.Sort(e)
	return e.Arr
}
func (e ImageSlice) Less(i, j int) bool {
	s := e.Arr
	return strings.Compare(s[i].ID, s[j].ID) < 0
}
