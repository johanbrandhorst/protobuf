package qunit

//GopherJS Bindings for qunitjs.com
import "github.com/gopherjs/gopherjs/js"

type QUnitAssert struct {
	*js.Object
}

type DoneCallbackObject struct {
	*js.Object
	Failed  int `js:"failed"`
	Passed  int `js:"passed"`
	Total   int `js:"total"`
	Runtime int `js:"runtime"`
}

type LogCallbackObject struct {
	*js.Object
	result   bool       `js:"result"`
	actual   *js.Object `js:"actual"`
	expected *js.Object `js:"expected"`
	message  string     `js:"message"`
	source   string     `js:"source"`
}

type ModuleStartCallbackObject struct {
	*js.Object
	name string `js:"name"`
}

type ModuleDoneCallbackObject struct {
	*js.Object
	name   string `js:"name"`
	failed int    `js:"failed"`
	passed int    `js:"passed"`
	total  int    `js:"total"`
}

type TestDoneCallbackObject struct {
	*js.Object
	name     string `js:"name"`
	module   string `js:"module"`
	failed   int    `js:"failed"`
	passed   int    `js:"passed"`
	total    int    `js:"total"`
	duration int    `js:"duration"`
}

type TestStartCallbackObject struct {
	*js.Object
	name   string `js:"name"`
	module string `js:"module"`
}

func (qa QUnitAssert) DeepEqual(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("deepEqual", actual, expected, message).Bool()
}

func (qa QUnitAssert) Equal(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("equal", actual, expected, message).Bool()
}

func (qa QUnitAssert) NotDeepEqual(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("notDeepEqual", actual, expected, message).Bool()
}

func (qa QUnitAssert) NotEqual(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("notEqual", actual, expected, message).Bool()
}

func (qa QUnitAssert) NotPropEqual(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("notPropEqual", actual, expected, message).Bool()
}

func (qa QUnitAssert) PropEqual(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("propEqual", actual, expected, message).Bool()
}

func (qa QUnitAssert) NotStrictEqual(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("notStrictEqual", actual, expected, message).Bool()
}

func (qa QUnitAssert) Ok(state interface{}, message string) bool {
	return qa.Call("ok", state, message).Bool()
}

func (qa QUnitAssert) StrictEqual(actual interface{}, expected interface{}, message string) bool {
	return qa.Call("strictEqual", actual, expected, message).Bool()
}

func (qa QUnitAssert) ThrowsExpected(block func() interface{}, expected interface{}, message string) *js.Object {
	return qa.Call("throwsExpected", block, expected, message)
}

func (qa QUnitAssert) Throws(block func() interface{}, message string) *js.Object {
	return qa.Call("throws", block, message)
}

func (qa QUnitAssert) Async() func() {
	// Use a closure to wrap around the async javascript object
	asyncObj := qa.Call("async")
	return func() {
		asyncObj.Invoke()
	}
}

//start QUnit static methods
func Test(name string, testFn func(QUnitAssert)) {
	js.Global.Get("QUnit").Call("test", name, func(e *js.Object) {
		testFn(QUnitAssert{Object: e})
	})
}

func TestExpected(title string, expected int, testFn func(assert QUnitAssert) interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("test", title, expected, func(e *js.Object) {
		testFn(QUnitAssert{Object: e})
	})
	return t
}

func Ok(state interface{}, message string) *js.Object {
	return js.Global.Get("QUnit").Call("ok", state, message)
}

func Start() *js.Object {
	return js.Global.Get("QUnit").Call("start")
}
func StartDecrement(decrement int) *js.Object {
	return js.Global.Get("QUnit").Call("start", decrement)
}
func Stop() *js.Object {
	return js.Global.Get("QUnit").Call("stop")
}
func StopIncrement(increment int) *js.Object {
	return js.Global.Get("QUnit").Call("stop", increment)
}

func Begin(callbackFn func() interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("begin", func() {
		callbackFn()
	})
	return t
}
func Done(callbackFn func(details DoneCallbackObject) interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("done", func(e *js.Object) {
		callbackFn(DoneCallbackObject{Object: e})
	})
	return t
}
func Log(callbackFn func(details LogCallbackObject) interface{}) *js.Object {
	/*
		2do:
		t := js.Global.Get("QUnit").Call("log", func(e *js.Object) {
			callbackFn(LogCallbackObject{Object: e})
		})
	*/
	t := js.Global.Get("QUnit").Call("log", callbackFn)
	return t
}
func ModuleDone(callbackFn func(details ModuleDoneCallbackObject) interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("moduleDone", func(e *js.Object) {
		callbackFn(ModuleDoneCallbackObject{Object: e})
	})
	return t
}
func ModuleStart(callbackFn func(name string) interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("moduleStart", func(e *js.Object) {
		callbackFn(e.String())
	})
	return t
}
func TestDone(callbackFn func(details TestDoneCallbackObject) interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("testDone", func(e *js.Object) {
		callbackFn(TestDoneCallbackObject{Object: e})
	})
	return t
}
func TestStart(callbackFn func(details TestStartCallbackObject) interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("testStart", func(e *js.Object) {
		callbackFn(TestStartCallbackObject{Object: e})
	})
	return t
}
func AsyncTestExpected(name string, expected interface{}, testFn func() interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("asyncTestExpected", name, expected, func() {
		testFn()
	})
	return t
}
func AsyncTest(name string, testFn func() interface{}) *js.Object {
	t := js.Global.Get("QUnit").Call("asyncTest", name, func() {
		testFn()
	})
	return t
}
func Expect(amount int) *js.Object {
	return js.Global.Get("QUnit").Call("expect", amount)
}

func Equiv(a interface{}, b interface{}) *js.Object {
	return js.Global.Get("QUnit").Call("equip", a, b)
}

func Module(name string) *js.Object {
	return js.Global.Get("QUnit").Call("module", name)
}

type Lifecycle interface {
	Setup()
	Teardown()
}

func ModuleLifecycle(name string, lc Lifecycle) *js.Object {
	o := js.Global.Get("Object").New()
	if lc.Setup != nil {
		o.Set("setup", lc.Setup)
	}
	if lc.Teardown != nil {
		o.Set("teardown", lc.Teardown)
	}
	return js.Global.Get("QUnit").Call("module", name, o)
}

type Raises struct {
	*js.Object
	Raises *js.Object `js:"raises"`
}

func Push(result interface{}, actual interface{}, expected interface{}, message string) *js.Object {
	return js.Global.Get("QUnit").Call("push", result, actual, expected, message)
}

func Reset() *js.Object {
	return js.Global.Get("QUnit").Call("reset")
}
