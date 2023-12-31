package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// TransactionManagerMock implements service.TransactionManager
type TransactionManagerMock struct {
	t minimock.Tester

	funcRunRepeatableRead          func(ctx context.Context, fn func(ctxTx context.Context) error) (err error)
	inspectFuncRunRepeatableRead   func(ctx context.Context, fn func(ctxTx context.Context) error)
	afterRunRepeatableReadCounter  uint64
	beforeRunRepeatableReadCounter uint64
	RunRepeatableReadMock          mTransactionManagerMockRunRepeatableRead

	funcSerializable          func(ctx context.Context, fn func(ctxTx context.Context) error) (err error)
	inspectFuncSerializable   func(ctx context.Context, fn func(ctxTx context.Context) error)
	afterSerializableCounter  uint64
	beforeSerializableCounter uint64
	SerializableMock          mTransactionManagerMockSerializable
}

// NewTransactionManagerMock returns a mock for service.TransactionManager
func NewTransactionManagerMock(t minimock.Tester) *TransactionManagerMock {
	m := &TransactionManagerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RunRepeatableReadMock = mTransactionManagerMockRunRepeatableRead{mock: m}
	m.RunRepeatableReadMock.callArgs = []*TransactionManagerMockRunRepeatableReadParams{}

	m.SerializableMock = mTransactionManagerMockSerializable{mock: m}
	m.SerializableMock.callArgs = []*TransactionManagerMockSerializableParams{}

	return m
}

type mTransactionManagerMockRunRepeatableRead struct {
	mock               *TransactionManagerMock
	defaultExpectation *TransactionManagerMockRunRepeatableReadExpectation
	expectations       []*TransactionManagerMockRunRepeatableReadExpectation

	callArgs []*TransactionManagerMockRunRepeatableReadParams
	mutex    sync.RWMutex
}

// TransactionManagerMockRunRepeatableReadExpectation specifies expectation struct of the TransactionManager.RunRepeatableRead
type TransactionManagerMockRunRepeatableReadExpectation struct {
	mock    *TransactionManagerMock
	params  *TransactionManagerMockRunRepeatableReadParams
	results *TransactionManagerMockRunRepeatableReadResults
	Counter uint64
}

// TransactionManagerMockRunRepeatableReadParams contains parameters of the TransactionManager.RunRepeatableRead
type TransactionManagerMockRunRepeatableReadParams struct {
	ctx context.Context
	fn  func(ctxTx context.Context) error
}

// TransactionManagerMockRunRepeatableReadResults contains results of the TransactionManager.RunRepeatableRead
type TransactionManagerMockRunRepeatableReadResults struct {
	err error
}

// Expect sets up expected params for TransactionManager.RunRepeatableRead
func (mmRunRepeatableRead *mTransactionManagerMockRunRepeatableRead) Expect(ctx context.Context, fn func(ctxTx context.Context) error) *mTransactionManagerMockRunRepeatableRead {
	if mmRunRepeatableRead.mock.funcRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("TransactionManagerMock.RunRepeatableRead mock is already set by Set")
	}

	if mmRunRepeatableRead.defaultExpectation == nil {
		mmRunRepeatableRead.defaultExpectation = &TransactionManagerMockRunRepeatableReadExpectation{}
	}

	mmRunRepeatableRead.defaultExpectation.params = &TransactionManagerMockRunRepeatableReadParams{ctx, fn}
	for _, e := range mmRunRepeatableRead.expectations {
		if minimock.Equal(e.params, mmRunRepeatableRead.defaultExpectation.params) {
			mmRunRepeatableRead.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRunRepeatableRead.defaultExpectation.params)
		}
	}

	return mmRunRepeatableRead
}

// Inspect accepts an inspector function that has same arguments as the TransactionManager.RunRepeatableRead
func (mmRunRepeatableRead *mTransactionManagerMockRunRepeatableRead) Inspect(f func(ctx context.Context, fn func(ctxTx context.Context) error)) *mTransactionManagerMockRunRepeatableRead {
	if mmRunRepeatableRead.mock.inspectFuncRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("Inspect function is already set for TransactionManagerMock.RunRepeatableRead")
	}

	mmRunRepeatableRead.mock.inspectFuncRunRepeatableRead = f

	return mmRunRepeatableRead
}

// Return sets up results that will be returned by TransactionManager.RunRepeatableRead
func (mmRunRepeatableRead *mTransactionManagerMockRunRepeatableRead) Return(err error) *TransactionManagerMock {
	if mmRunRepeatableRead.mock.funcRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("TransactionManagerMock.RunRepeatableRead mock is already set by Set")
	}

	if mmRunRepeatableRead.defaultExpectation == nil {
		mmRunRepeatableRead.defaultExpectation = &TransactionManagerMockRunRepeatableReadExpectation{mock: mmRunRepeatableRead.mock}
	}
	mmRunRepeatableRead.defaultExpectation.results = &TransactionManagerMockRunRepeatableReadResults{err}
	return mmRunRepeatableRead.mock
}

// Set uses given function f to mock the TransactionManager.RunRepeatableRead method
func (mmRunRepeatableRead *mTransactionManagerMockRunRepeatableRead) Set(f func(ctx context.Context, fn func(ctxTx context.Context) error) (err error)) *TransactionManagerMock {
	if mmRunRepeatableRead.defaultExpectation != nil {
		mmRunRepeatableRead.mock.t.Fatalf("Default expectation is already set for the TransactionManager.RunRepeatableRead method")
	}

	if len(mmRunRepeatableRead.expectations) > 0 {
		mmRunRepeatableRead.mock.t.Fatalf("Some expectations are already set for the TransactionManager.RunRepeatableRead method")
	}

	mmRunRepeatableRead.mock.funcRunRepeatableRead = f
	return mmRunRepeatableRead.mock
}

// When sets expectation for the TransactionManager.RunRepeatableRead which will trigger the result defined by the following
// Then helper
func (mmRunRepeatableRead *mTransactionManagerMockRunRepeatableRead) When(ctx context.Context, fn func(ctxTx context.Context) error) *TransactionManagerMockRunRepeatableReadExpectation {
	if mmRunRepeatableRead.mock.funcRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("TransactionManagerMock.RunRepeatableRead mock is already set by Set")
	}

	expectation := &TransactionManagerMockRunRepeatableReadExpectation{
		mock:   mmRunRepeatableRead.mock,
		params: &TransactionManagerMockRunRepeatableReadParams{ctx, fn},
	}
	mmRunRepeatableRead.expectations = append(mmRunRepeatableRead.expectations, expectation)
	return expectation
}

// Then sets up TransactionManager.RunRepeatableRead return parameters for the expectation previously defined by the When method
func (e *TransactionManagerMockRunRepeatableReadExpectation) Then(err error) *TransactionManagerMock {
	e.results = &TransactionManagerMockRunRepeatableReadResults{err}
	return e.mock
}

// RunRepeatableRead implements service.TransactionManager
func (mmRunRepeatableRead *TransactionManagerMock) RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) (err error) {
	mm_atomic.AddUint64(&mmRunRepeatableRead.beforeRunRepeatableReadCounter, 1)
	defer mm_atomic.AddUint64(&mmRunRepeatableRead.afterRunRepeatableReadCounter, 1)

	if mmRunRepeatableRead.inspectFuncRunRepeatableRead != nil {
		mmRunRepeatableRead.inspectFuncRunRepeatableRead(ctx, fn)
	}

	mm_params := &TransactionManagerMockRunRepeatableReadParams{ctx, fn}

	// Record call args
	mmRunRepeatableRead.RunRepeatableReadMock.mutex.Lock()
	mmRunRepeatableRead.RunRepeatableReadMock.callArgs = append(mmRunRepeatableRead.RunRepeatableReadMock.callArgs, mm_params)
	mmRunRepeatableRead.RunRepeatableReadMock.mutex.Unlock()

	for _, e := range mmRunRepeatableRead.RunRepeatableReadMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation.Counter, 1)
		mm_want := mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation.params
		mm_got := TransactionManagerMockRunRepeatableReadParams{ctx, fn}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRunRepeatableRead.t.Errorf("TransactionManagerMock.RunRepeatableRead got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation.results
		if mm_results == nil {
			mmRunRepeatableRead.t.Fatal("No results are set for the TransactionManagerMock.RunRepeatableRead")
		}
		return (*mm_results).err
	}
	if mmRunRepeatableRead.funcRunRepeatableRead != nil {
		return mmRunRepeatableRead.funcRunRepeatableRead(ctx, fn)
	}
	mmRunRepeatableRead.t.Fatalf("Unexpected call to TransactionManagerMock.RunRepeatableRead. %v %v", ctx, fn)
	return
}

// RunRepeatableReadAfterCounter returns a count of finished TransactionManagerMock.RunRepeatableRead invocations
func (mmRunRepeatableRead *TransactionManagerMock) RunRepeatableReadAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunRepeatableRead.afterRunRepeatableReadCounter)
}

// RunRepeatableReadBeforeCounter returns a count of TransactionManagerMock.RunRepeatableRead invocations
func (mmRunRepeatableRead *TransactionManagerMock) RunRepeatableReadBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunRepeatableRead.beforeRunRepeatableReadCounter)
}

// Calls returns a list of arguments used in each call to TransactionManagerMock.RunRepeatableRead.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRunRepeatableRead *mTransactionManagerMockRunRepeatableRead) Calls() []*TransactionManagerMockRunRepeatableReadParams {
	mmRunRepeatableRead.mutex.RLock()

	argCopy := make([]*TransactionManagerMockRunRepeatableReadParams, len(mmRunRepeatableRead.callArgs))
	copy(argCopy, mmRunRepeatableRead.callArgs)

	mmRunRepeatableRead.mutex.RUnlock()

	return argCopy
}

// MinimockRunRepeatableReadDone returns true if the count of the RunRepeatableRead invocations corresponds
// the number of defined expectations
func (m *TransactionManagerMock) MinimockRunRepeatableReadDone() bool {
	for _, e := range m.RunRepeatableReadMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunRepeatableReadMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunRepeatableRead != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		return false
	}
	return true
}

// MinimockRunRepeatableReadInspect logs each unmet expectation
func (m *TransactionManagerMock) MinimockRunRepeatableReadInspect() {
	for _, e := range m.RunRepeatableReadMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactionManagerMock.RunRepeatableRead with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunRepeatableReadMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		if m.RunRepeatableReadMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TransactionManagerMock.RunRepeatableRead")
		} else {
			m.t.Errorf("Expected call to TransactionManagerMock.RunRepeatableRead with params: %#v", *m.RunRepeatableReadMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunRepeatableRead != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		m.t.Error("Expected call to TransactionManagerMock.RunRepeatableRead")
	}
}

type mTransactionManagerMockSerializable struct {
	mock               *TransactionManagerMock
	defaultExpectation *TransactionManagerMockSerializableExpectation
	expectations       []*TransactionManagerMockSerializableExpectation

	callArgs []*TransactionManagerMockSerializableParams
	mutex    sync.RWMutex
}

// TransactionManagerMockSerializableExpectation specifies expectation struct of the TransactionManager.Serializable
type TransactionManagerMockSerializableExpectation struct {
	mock    *TransactionManagerMock
	params  *TransactionManagerMockSerializableParams
	results *TransactionManagerMockSerializableResults
	Counter uint64
}

// TransactionManagerMockSerializableParams contains parameters of the TransactionManager.Serializable
type TransactionManagerMockSerializableParams struct {
	ctx context.Context
	fn  func(ctxTx context.Context) error
}

// TransactionManagerMockSerializableResults contains results of the TransactionManager.Serializable
type TransactionManagerMockSerializableResults struct {
	err error
}

// Expect sets up expected params for TransactionManager.Serializable
func (mmSerializable *mTransactionManagerMockSerializable) Expect(ctx context.Context, fn func(ctxTx context.Context) error) *mTransactionManagerMockSerializable {
	if mmSerializable.mock.funcSerializable != nil {
		mmSerializable.mock.t.Fatalf("TransactionManagerMock.Serializable mock is already set by Set")
	}

	if mmSerializable.defaultExpectation == nil {
		mmSerializable.defaultExpectation = &TransactionManagerMockSerializableExpectation{}
	}

	mmSerializable.defaultExpectation.params = &TransactionManagerMockSerializableParams{ctx, fn}
	for _, e := range mmSerializable.expectations {
		if minimock.Equal(e.params, mmSerializable.defaultExpectation.params) {
			mmSerializable.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSerializable.defaultExpectation.params)
		}
	}

	return mmSerializable
}

// Inspect accepts an inspector function that has same arguments as the TransactionManager.Serializable
func (mmSerializable *mTransactionManagerMockSerializable) Inspect(f func(ctx context.Context, fn func(ctxTx context.Context) error)) *mTransactionManagerMockSerializable {
	if mmSerializable.mock.inspectFuncSerializable != nil {
		mmSerializable.mock.t.Fatalf("Inspect function is already set for TransactionManagerMock.Serializable")
	}

	mmSerializable.mock.inspectFuncSerializable = f

	return mmSerializable
}

// Return sets up results that will be returned by TransactionManager.Serializable
func (mmSerializable *mTransactionManagerMockSerializable) Return(err error) *TransactionManagerMock {
	if mmSerializable.mock.funcSerializable != nil {
		mmSerializable.mock.t.Fatalf("TransactionManagerMock.Serializable mock is already set by Set")
	}

	if mmSerializable.defaultExpectation == nil {
		mmSerializable.defaultExpectation = &TransactionManagerMockSerializableExpectation{mock: mmSerializable.mock}
	}
	mmSerializable.defaultExpectation.results = &TransactionManagerMockSerializableResults{err}
	return mmSerializable.mock
}

// Set uses given function f to mock the TransactionManager.Serializable method
func (mmSerializable *mTransactionManagerMockSerializable) Set(f func(ctx context.Context, fn func(ctxTx context.Context) error) (err error)) *TransactionManagerMock {
	if mmSerializable.defaultExpectation != nil {
		mmSerializable.mock.t.Fatalf("Default expectation is already set for the TransactionManager.Serializable method")
	}

	if len(mmSerializable.expectations) > 0 {
		mmSerializable.mock.t.Fatalf("Some expectations are already set for the TransactionManager.Serializable method")
	}

	mmSerializable.mock.funcSerializable = f
	return mmSerializable.mock
}

// When sets expectation for the TransactionManager.Serializable which will trigger the result defined by the following
// Then helper
func (mmSerializable *mTransactionManagerMockSerializable) When(ctx context.Context, fn func(ctxTx context.Context) error) *TransactionManagerMockSerializableExpectation {
	if mmSerializable.mock.funcSerializable != nil {
		mmSerializable.mock.t.Fatalf("TransactionManagerMock.Serializable mock is already set by Set")
	}

	expectation := &TransactionManagerMockSerializableExpectation{
		mock:   mmSerializable.mock,
		params: &TransactionManagerMockSerializableParams{ctx, fn},
	}
	mmSerializable.expectations = append(mmSerializable.expectations, expectation)
	return expectation
}

// Then sets up TransactionManager.Serializable return parameters for the expectation previously defined by the When method
func (e *TransactionManagerMockSerializableExpectation) Then(err error) *TransactionManagerMock {
	e.results = &TransactionManagerMockSerializableResults{err}
	return e.mock
}

// Serializable implements service.TransactionManager
func (mmSerializable *TransactionManagerMock) Serializable(ctx context.Context, fn func(ctxTx context.Context) error) (err error) {
	mm_atomic.AddUint64(&mmSerializable.beforeSerializableCounter, 1)
	defer mm_atomic.AddUint64(&mmSerializable.afterSerializableCounter, 1)

	if mmSerializable.inspectFuncSerializable != nil {
		mmSerializable.inspectFuncSerializable(ctx, fn)
	}

	mm_params := &TransactionManagerMockSerializableParams{ctx, fn}

	// Record call args
	mmSerializable.SerializableMock.mutex.Lock()
	mmSerializable.SerializableMock.callArgs = append(mmSerializable.SerializableMock.callArgs, mm_params)
	mmSerializable.SerializableMock.mutex.Unlock()

	for _, e := range mmSerializable.SerializableMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmSerializable.SerializableMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSerializable.SerializableMock.defaultExpectation.Counter, 1)
		mm_want := mmSerializable.SerializableMock.defaultExpectation.params
		mm_got := TransactionManagerMockSerializableParams{ctx, fn}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSerializable.t.Errorf("TransactionManagerMock.Serializable got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSerializable.SerializableMock.defaultExpectation.results
		if mm_results == nil {
			mmSerializable.t.Fatal("No results are set for the TransactionManagerMock.Serializable")
		}
		return (*mm_results).err
	}
	if mmSerializable.funcSerializable != nil {
		return mmSerializable.funcSerializable(ctx, fn)
	}
	mmSerializable.t.Fatalf("Unexpected call to TransactionManagerMock.Serializable. %v %v", ctx, fn)
	return
}

// SerializableAfterCounter returns a count of finished TransactionManagerMock.Serializable invocations
func (mmSerializable *TransactionManagerMock) SerializableAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSerializable.afterSerializableCounter)
}

// SerializableBeforeCounter returns a count of TransactionManagerMock.Serializable invocations
func (mmSerializable *TransactionManagerMock) SerializableBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSerializable.beforeSerializableCounter)
}

// Calls returns a list of arguments used in each call to TransactionManagerMock.Serializable.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSerializable *mTransactionManagerMockSerializable) Calls() []*TransactionManagerMockSerializableParams {
	mmSerializable.mutex.RLock()

	argCopy := make([]*TransactionManagerMockSerializableParams, len(mmSerializable.callArgs))
	copy(argCopy, mmSerializable.callArgs)

	mmSerializable.mutex.RUnlock()

	return argCopy
}

// MinimockSerializableDone returns true if the count of the Serializable invocations corresponds
// the number of defined expectations
func (m *TransactionManagerMock) MinimockSerializableDone() bool {
	for _, e := range m.SerializableMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SerializableMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSerializableCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSerializable != nil && mm_atomic.LoadUint64(&m.afterSerializableCounter) < 1 {
		return false
	}
	return true
}

// MinimockSerializableInspect logs each unmet expectation
func (m *TransactionManagerMock) MinimockSerializableInspect() {
	for _, e := range m.SerializableMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactionManagerMock.Serializable with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SerializableMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSerializableCounter) < 1 {
		if m.SerializableMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TransactionManagerMock.Serializable")
		} else {
			m.t.Errorf("Expected call to TransactionManagerMock.Serializable with params: %#v", *m.SerializableMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSerializable != nil && mm_atomic.LoadUint64(&m.afterSerializableCounter) < 1 {
		m.t.Error("Expected call to TransactionManagerMock.Serializable")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TransactionManagerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockRunRepeatableReadInspect()

		m.MinimockSerializableInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TransactionManagerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *TransactionManagerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockRunRepeatableReadDone() &&
		m.MinimockSerializableDone()
}
