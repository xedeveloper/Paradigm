package core

import (
	"html/template"
	"net/http"
	"sync"

	lifecycle "github.com/xedeveloper/Paradigm/engine/core/life_cycle"
)

type ComponentInterface interface {
	GetTemplate() *template.Template
	GetState() interface{}
	GetLifeCycle() *lifecycle.LifeCycleHooks
	Render(w http.ResponseWriter) error
}

type BaseComponent struct {
	template  *template.Template
	state     interface{}
	lifecycle *lifecycle.LifeCycleHooks
	mutex     sync.RWMutex
}

type ComponentConfig struct {
	Selector string
	Template string
	Style    string
}

func NewComponent(config ComponentConfig) *BaseComponent {
	tmpl := template.Must(template.New(config.Selector).Parse(config.Template))
	return &BaseComponent{
		template:  tmpl,
		lifecycle: &lifecycle.LifeCycleHooks{},
	}
}

func (c *BaseComponent) SetState(newState interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.state = newState
	if c.lifecycle.OnUpdate != nil {
		c.lifecycle.OnUpdate()
	}
}

func (c *BaseComponent) GetState() interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.state
}

func (c *BaseComponent) GetLifeCycle() *lifecycle.LifeCycleHooks {
	return c.lifecycle
}

func (c *BaseComponent) Render(w http.ResponseWriter) error {
	return c.template.Execute(w, c.state)
}

func (c *BaseComponent) GetTemplate() *template.Template {
	return c.template
}
