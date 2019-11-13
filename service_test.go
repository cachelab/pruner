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

	os.Setenv("DELAY", "FAIL")
	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("DELAY", "43200")
	os.Setenv("DAYS_TO_KEEP", "FAIL")
	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("DAYS_TO_KEEP", "7")
	os.Setenv("MAX_RETRIES", "FAIL")
	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("MAX_RETRIES", "1")
	os.Setenv("RUN_ONCE", "FAIL")
	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("RUN_ONCE", "true")
	os.Setenv("PREFIXES", "logs")
	err = svc.Init()
	assert.Equal(t, err, nil)

	os.Setenv("ELASTICSEARCH_URL", "http://127.0.0.1:9200")
	err = svc.Init()
	assert.Equal(t, err, nil)
}
