package core

import (
	"html/template"
	"sync"

	lifecycle "github.com/xedeveloper/Paradigm/engine/core/life_cycle"
)

type Component struct {
	Template  *template.Template
	State     interface{}
	Lifecycle *lifecycle.LifeCycleHooks
	Router    *Router
	mutex     sync.RWMutex
}

type ComponentConfig struct {
	Selector string
	Template string
	Style    string
}
