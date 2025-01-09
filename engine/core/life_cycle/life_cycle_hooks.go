package lifecycle

type LifeCycleHooks struct {
	OnInit      func()
	OnDestroy   func()
	OnUpdate    func()
	BeforeMount func()
	AfterMount  func()
}
