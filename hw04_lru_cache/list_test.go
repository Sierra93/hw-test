package hw04lrucache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkedList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("remove front and back", func(t *testing.T) {
		l := NewList()
		for _, v := range [...]int{10, 20, 30, 40, 50, 60, 70, 80} {
			l.PushBack(v)
		}

		l.Remove(l.Front())
		l.Remove(l.Back())

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}

		require.Equal(t, []int{20, 30, 40, 50, 60, 70}, elems)
	})

	t.Run("move all elements from back to front", func(t *testing.T) {
		l := NewList()
		for _, v := range [...]int{10, 20, 30, 40, 50, 60, 70, 80} {
			l.PushBack(v)
		}

		length := l.Len()
		for i := 0; i < length; i++ {
			l.MoveToFront(l.Back())
		}

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}

		require.Equal(t, []int{10, 20, 30, 40, 50, 60, 70, 80}, elems)
	})

	t.Run("remove all elements", func(t *testing.T) {
		l := NewList()
		for _, v := range [...]int{10, 20, 30, 40, 50, 60, 70, 80} {
			l.PushBack(v)
		}

		curr := l.Front()
		for curr != nil {
			next := curr.Next
			l.Remove(curr)
			curr = next
		}

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}

func BenchmarkListPushBackAndRemove(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			list := NewList()
			for i := 0; i < size; i++ {
				list.PushFront(i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f := list.Front()
				list.PushBack(f.Value)
				list.Remove(f)
			}
		})
	}
}

func BenchmarkListPushFrontAndRemove(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			list := NewList()
			for i := 0; i < size; i++ {
				list.PushFront(i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f := list.Back()
				list.PushFront(f.Value)
				list.Remove(f)
			}
		})
	}
}

func BenchmarkListMoveToFront(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			list := NewList()
			for i := 0; i < size; i++ {
				list.PushFront(i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.MoveToFront(list.Back())
			}
		})
	}
}
