package vagrant

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMachineReadableItemsFromString(t *testing.T) {
	emptyItms := []MachineReadableItem{}
	require.Equal(t, emptyItms, ParseMachineReadableItemsFromString(``, "", ""))
	require.Equal(t, emptyItms, ParseMachineReadableItemsFromString(`abc`, "", ""))

	// a complex 'status' example
	inStr := `some gibberish

here

1443644808,default,provider-name,example-provider
and there
1443644808,default,state,stopped
1443644808,default,state-human-short,stopped
1443644808,default,state-human-long,The VM is stopped. To start the VM%!(VAGRANT_COMMA) simply ...`

	expectedItems := []MachineReadableItem{
		MachineReadableItem{1443644808, "default", "provider-name", "example-provider"},
		MachineReadableItem{1443644808, "default", "state", "stopped"},
		MachineReadableItem{1443644808, "default", "state-human-short", "stopped"},
		MachineReadableItem{1443644808, "default", "state-human-long", "The VM is stopped. To start the VM%!(VAGRANT_COMMA) simply ..."},
	}
	// no filter
	require.Equal(t, expectedItems, ParseMachineReadableItemsFromString(inStr, "", ""))
	// target filter
	require.Equal(t, expectedItems, ParseMachineReadableItemsFromString(inStr, "default", ""))
	// target filter - no result
	require.Equal(t, emptyItms, ParseMachineReadableItemsFromString(inStr, "no-target", ""))
	// type filter - only one
	require.Equal(t, []MachineReadableItem{MachineReadableItem{1443644808, "default", "state", "stopped"}},
		ParseMachineReadableItemsFromString(inStr, "", "state"))
	// type filter - no result
	require.Equal(t, emptyItms, ParseMachineReadableItemsFromString(inStr, "", "no-type"))
}
