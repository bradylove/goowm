package windowmanager

import (
	"goowm/config"

	"github.com/BurntSushi/xgbutil"
)

type WorkspaceManager struct {
	workspaces  []*Workspace
	activeIndex int
}

func NewWorkspaceManager(initialCount int) *WorkspaceManager {
	return &WorkspaceManager{
		workspaces: make([]*Workspace, 0, initialCount),
	}
}

func (wm *WorkspaceManager) Add(x *xgbutil.XUtil, wc ...*config.WorkspaceConfig) {
	for _, conf := range wc {
		wm.workspaces = append(wm.workspaces, NewWorkspace(x, conf))
	}
}

func (wm *WorkspaceManager) Names() []string {
	names := make([]string, 0, len(wm.workspaces))
	for _, ws := range wm.workspaces {
		names = append(names, ws.Name())
	}

	return names
}

func (wm *WorkspaceManager) NextIndex() int {
	var index int
	if wm.activeIndex != len(wm.workspaces)-1 {
		index = wm.activeIndex + 1
	}

	return index
}

func (wm *WorkspaceManager) PreviousIndex() int {
	index := wm.activeIndex - 1

	if index == -1 {
		index = len(wm.workspaces) - 1
	}

	return index
}

func (wm *WorkspaceManager) Activate(index int) {
	wm.workspaces[wm.activeIndex].Deactivate()
	wm.workspaces[index].Activate()
	wm.activeIndex = index
}
