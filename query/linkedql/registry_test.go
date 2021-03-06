package linkedql

import (
	"testing"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/query"
	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/cayley/query/path"

	"github.com/stretchr/testify/require"
)

func init() {
	Register(&TestStep{})
}

var unmarshalCases = []struct {
	name string
	data string
	exp  Step
}{
	{
		name: "simple",
		data: `{
	"@type": "cayley:TestStep",
	"limit": 10
}`,
		exp: &TestStep{Limit: 10},
	},
	{
		name: "simple",
		data: `{
	"@type": "cayley:TestStep",
	"tags": ["a", "b"]
}`,
		exp: &TestStep{Tags: []string{"a", "b"}},
	},
	{
		name: "nested",
		data: `{
	"@type": "cayley:TestStep",
	"limit": 10,
	"main": {
		"@type": "cayley:TestStep",
		"limit": 15,
		"main": {
			"@type": "cayley:TestStep",
			"limit": 20
		}
	}
}`,
		exp: &TestStep{
			Limit: 10,
			Main: &TestStep{
				Limit: 15,
				Main: &TestStep{
					Limit: 20,
				},
			},
		},
	},
	{
		name: "nested slice",
		data: `{
	"@type": "cayley:TestStep",
	"limit": 10,
	"sub": [
		{
			"@type": "cayley:TestStep",
			"limit": 15
		},
		{
			"@type": "cayley:TestStep",
			"limit": 20
		}
	]
}`,
		exp: &TestStep{
			Limit: 10,
			Sub: []PathStep{
				&TestStep{
					Limit: 15,
				},
				&TestStep{
					Limit: 20,
				},
			},
		},
	},
}

type TestStep struct {
	Limit int        `json:"limit"`
	Tags  []string   `json:"tags"`
	Main  PathStep   `json:"main"`
	Sub   []PathStep `json:"sub"`
}

func (s *TestStep) Type() quad.IRI {
	return "cayley:TestStep"
}

func (s *TestStep) Description() string {
	return "A TestStep for checking the registry"
}

func (s *TestStep) BuildIterator(qs graph.QuadStore) (query.Iterator, error) {
	panic("Can't build iterator for TestStep")
}

func (s *TestStep) BuildPath(qs graph.QuadStore) (*path.Path, error) {
	panic("Can't build path for TestStep")
}

func TestUnmarshalStep(t *testing.T) {
	for _, c := range unmarshalCases {
		t.Run(c.name, func(t *testing.T) {
			s, err := Unmarshal([]byte(c.data))
			require.NoError(t, err)
			require.Equal(t, c.exp, s)
		})
	}
}
