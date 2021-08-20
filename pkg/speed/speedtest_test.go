package speed

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testProvider struct {
	fn func(context.Context, chan<- Measures) error
}

func (p testProvider) Test(ctx context.Context, ch chan<- Measures) error {
	if p.fn == nil {
		return nil
	}
	return p.fn(ctx, ch)
}

func TestProviderTest_NoError(t *testing.T) {
	provider := testProvider{
		fn: func(ctx context.Context, ch chan<- Measures) error {
			close(ch)
			return nil
		},
	}
	measures, err := ProviderTest(provider)
	assert.NoError(t, err)
	assert.Equal(t, Measures{}, measures)
}

func TestProviderTest_Error(t *testing.T) {
	expectedErr := errors.New("error!!")
	provider := testProvider{
		fn: func(ctx context.Context, ch chan<- Measures) error {
			close(ch)
			return expectedErr
		},
	}
	measures, err := ProviderTest(provider)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, Measures{}, measures)
}

func TestTestFn_NoProvider(t *testing.T) {
	err := Test(context.Background(), nil, nil)
	require.Error(t, err)
}

func TestTestFn_NoResultChannel(t *testing.T) {
	err := Test(context.Background(), testProvider{}, nil)
	require.Error(t, err)
}

func TestTestFn_NoContext(t *testing.T) {
	provider := testProvider{
		fn: func(ctx context.Context, ch chan<- Measures) error {
			require.NotNil(t, ctx)
			close(ch)
			return nil
		},
	}
	err := Test(nil, provider, make(chan<- Measures, 1))
	require.NoError(t, err)
}
