package examples

type MyStruct struct{}

func (ms *MyStruct) A(a int, b uint) (c string, err error) {}

func (ms *MyStruct) B(a int, b uint) {}

func (ms *MyStruct) C() (c string, err error) {}

func (ms *MyStruct) D() error {}

func (ms *MyStruct) E(a int) (string, error) {}

func (ms *MyStruct) F() {}
