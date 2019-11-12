package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	svc := Service{}

	err := svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("ELASTICSEARCH_URL", "http://127.0.0.1:9200")
	os.Setenv("DELAY", "43200")
	os.Setenv("PREFIXES", "logs")
	os.Setenv("DAYS_TO_KEEP", "30")
	os.Setenv("RUN_ONCE", "true")

	err = svc.Init()
	assert.Equal(t, err, nil)

	os.Setenv("DELAY", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("DELAY", "43200")
	os.Setenv("DAYS_TO_KEEP", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("DELAY", "43200")
	os.Setenv("DAYS_TO_KEEP", "")
	os.Setenv("RUN_ONCE", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("DELAY", "43200")
	os.Setenv("DAYS_TO_KEEP", "")
	os.Setenv("RUN_ONCE", "")
	os.Setenv("ELASTICSEARCH_URL", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)
}
