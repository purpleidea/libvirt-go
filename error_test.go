package libvirt

import "testing"

func TestGlobalErrorCallback(t *testing.T) {
	var nbErrors int
	errors := make([]VirError, 0, 10)
	callback := ErrorCallback(func(err VirError, f func()) {
		errors = append(errors, err)
		f()
	})
	SetErrorFunc(callback, func() {
		nbErrors++
	})
	NewVirConnection("invalid_transport:///default")
	if len(errors) == 0 {
		t.Errorf("No errors were captured")
	}
	if len(errors) != nbErrors {
		t.Errorf("Captured %d errors (%+v) but counted only %d errors",
			len(errors), errors, nbErrors)
	}
	errors = make([]VirError, 0, 10)
	SetErrorFunc(nil, nil)
	NewVirConnection("invalid_transport:///default")
	if len(errors) != 0 {
		t.Errorf("More errors have been captured: %+v", errors)
	}
}

func TestConnectionErrorCallback(t *testing.T) {
	var nbErrors int
	errors := make([]VirError, 0, 10)
	callback := ErrorCallback(func(err VirError, f func()) {
		errors = append(errors, err)
		f()
	})
	conn := buildTestConnection()
	conn.SetErrorFunc(callback, func() {
		nbErrors++
	})
	defer conn.UnsetErrorFunc()

	// To generate an error, we set memory of a domain to an insance value
	domain, err := conn.LookupDomainByName("test")
	if err != nil {
		panic(err)
	}
	err = domain.SetMemory(100000000000)
	if err == nil {
		t.Fatalf("Was expecting an error when setting memory to too high value")
	}

	if len(errors) == 0 {
		t.Errorf("No errors were captured")
	}
	if len(errors) != nbErrors {
		t.Errorf("Captured %d errors (%+v) but counted only %d errors",
			len(errors), errors, nbErrors)
	}
	errors = make([]VirError, 0, 10)
	conn.UnsetErrorFunc()
	if len(goCallbacks) != 0 {
		t.Errorf("goCallbacks entry wasn't removed: %+v", goCallbacks)
	}
	err = domain.SetMemory(100000000000)
	if err == nil {
		t.Fatalf("Was expecting an error when setting memory to too high value")
	}
	if len(errors) != 0 {
		t.Errorf("More errors have been captured: %+v", errors)
	}
}
