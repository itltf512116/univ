package resolve_test

import (
	"testing"

	"github.com/itltf512116/univ/resolve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var client resolve.UnivClient

func TestMain(m *testing.M) {
	client = resolve.NewUnivClient()

	m.Run()
}

func TestResolveAllCountries(t *testing.T) {
	cs, err := client.ResolveAllCountries()

	require.NoError(t, err)

	t.Logf("the count of countries :%d", len(cs))
}

func TestResolveUniversityByCountry(t *testing.T) {
	cty := resolve.Country{Name: "Australia", Code: "au"} // 52 university
	us, err := client.ResolveUniversityByCountry(cty)

	require.NoError(t, err)

	assert.Equal(t, 52, len(us))
	// fmt.Printf("%+v", us)
}
