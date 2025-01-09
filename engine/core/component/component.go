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

func NewComponent(config ComponentConfig) *Component {
	tmpl := template.Must(template.New(config.Selector).Parse(config.Template))
	return &Component{
		Template:  tmpl,
		Lifecycle: &lifecycle.LifeCycleHooks{},
	}
}

func (c *Component) SetState(newState interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.State = newState
	if c.Lifecycle.OnUpdate != nil {
		c.Lifecycle.OnUpdate()
	}
}
