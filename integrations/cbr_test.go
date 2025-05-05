package integrations

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCBRIntegration(t *testing.T) {
	rate, err := GetCentralBankRate()
	assert.NoError(t, err)
	assert.Greater(t, rate, 0.0)
	fmt.Printf("Ключевая ставка с маржой банка: %.2f%%\n", rate)
}
