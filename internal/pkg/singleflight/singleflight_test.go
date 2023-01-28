package singleflight

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroup_Do(t *testing.T) {
	wg := new(sync.WaitGroup)
	nums := int64(0)
	cnt := 100
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func(n int) {
			defer wg.Done()
			_, err := Group.Do("redis", func() (interface{}, error) {
				atomic.AddInt64(&nums, 1)
				return n, nil
			})
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()
	t.Log(cnt, nums)
	require.True(t, nums <= int64(cnt))
}

func TestGroup_WithoutDo(t *testing.T) {
	wg := new(sync.WaitGroup)
	cnt := 100
	nums := int64(0)
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func(n int) {
			defer wg.Done()
			_, err := func() (interface{}, error) {
				atomic.AddInt64(&nums, 1)
				return n, nil
			}()
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()
	require.EqualValues(t, nums, cnt)
}
