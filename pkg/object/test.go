package object

type Test struct {
	Kind string `json:"kind" yaml:"kind"`
	Name string `json:"name" yaml:"name"`
	Spec []Inf  `json:"spec" yaml:"spec"`
}

type Test2 struct {
	Kind string `json:"kind" yaml:"kind"`
	Name string `json:"name" yaml:"name"`
	Spec []Inf2 `json:"spec" yaml:"spec"`
}

type Inf struct {
	Home string `json:"home" yaml:"home"`
	Age  int    `json:"age" yaml:"age"`
}

type Inf2 struct {
	Id      int    `json:"id" yaml:"id"`
	Class   string `json:"class" yaml:"class"`
	Teacher string `json:"teacher" yaml:"teacher"`
}
