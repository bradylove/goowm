package windowmanager

import (
	"goowm/broadcaster"
	"goowm/config"

	"github.com/BurntSushi/xgbutil"
)

type WorkspaceManager struct {
	Workspaces  []*Workspace
	activeIndex int
}

func NewWorkspaceManager(initialCount int) *WorkspaceManager {
	return &WorkspaceManager{
		Workspaces: make([]*Workspace, 0, initialCount),
	}
}

func (wm *WorkspaceManager) Add(x *xgbutil.XUtil, wc ...*config.WorkspaceConfig) {
	for _, conf := range wc {
		wm.Workspaces = append(wm.Workspaces, NewWorkspace(x, conf))
	}
}

func (wm *WorkspaceManager) Names() []string {
	names := make([]string, 0, len(wm.Workspaces))
	for _, ws := range wm.Workspaces {
		names = append(names, ws.Name())
	}

	return names
}

func (wm *WorkspaceManager) NextIndex() int {
	var index int
	if wm.activeIndex != len(wm.Workspaces)-1 {
		index = wm.activeIndex + 1
	}

	return index
}

func (wm *WorkspaceManager) PreviousIndex() int {
	index := wm.activeIndex - 1

	if index == -1 {
		index = len(wm.Workspaces) - 1
	}

	return index
}

func (wm *WorkspaceManager) Activate(index int) {
	wm.Workspaces[wm.activeIndex].Deactivate()
	wm.Workspaces[index].Activate()
	wm.activeIndex = index

	broadcaster.Trigger(broadcaster.EventWorkspaceChanged)
}

func (wm *WorkspaceManager) ActiveIndex() int {
	return wm.activeIndex
}
