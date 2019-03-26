package util

import (
	"github.com/cdyfng/go_binance_demo/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)
func TestNames(t *testing.T) {
  	name := ""
  	for i:=0; i<101; i++ {
		name = util.MakeRandomStrLower(12)
		fmt.Printf("name: %v \n", name)
  	}
	//t.Logf("name:  %v \n", name)
	assert.Equal(t, len(name), 12)
}
