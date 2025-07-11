package mockutil

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
)

// Helper contains helper functions for working with gomock
type Helper struct {
	controller *gomock.Controller
	ctx        context.Context
}

// Controller returns a gomock Controller that should be used
// when creating new mock objects
func (p *Helper) Controller() *gomock.Controller {
	return p.controller
}

// Context returns a Context that is unique and will be canceled on
// test failure
func (p *Helper) Context() context.Context {
	return p.ctx
}

// MockableFunc represents a ginkgo Describe closure function that wishes
// to have mocking enabled.
type MockableFunc func(helper *Helper)

// Mockable allows Ginkgo tests to be easily written with
// gomock support.
//
// You will need to wrap your top-level Describe as follows:
//
//	    var _ = Describe("SomethingTest", mockutil.Mockable(func(helper *mockutil.Helper) {
//		       BeforeEach(func() {
//	            dependency := NewMockDependency(helper.Controller())
//	        })
//	    }))
func Mockable(mockableClosure MockableFunc) func() {
	return func() {
		gomockHelper := &Helper{}

		BeforeEach(func() {
			gomockHelper.controller, gomockHelper.ctx = gomock.WithContext(context.Background(), GinkgoT())
		})

		AfterEach(func() {
			gomockHelper.controller.Finish()
		})

		Describe("", func() {
			mockableClosure(gomockHelper)
		})
	}
}
