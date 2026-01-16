package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) { // nolint
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
}

func TestCache_SetAndGet_UpdateValueAndOrder(t *testing.T) { // nolint
	cache := NewCache(2)

	// Первый вызов - добавляем новый ключ
	wasInCache := cache.Set("key1", "value1")
	require.False(t, wasInCache, "Первое добавление должно возвращать false")

	// Проверяем, что получаем правильное значение
	val, ok := cache.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", val)

	// Обновляем значение того же ключа и проверяем, что оно обновилось
	wasInCache = cache.Set("key1", "new_value1")
	require.True(t, wasInCache, "Обновление существующего ключа должно возвращать true")
	val, ok = cache.Get("key1")
	require.True(t, ok)
	require.Equal(t, "new_value1", val)

	// Добавляем еще один ключ
	wasInCache = cache.Set("key2", "value2")
	require.False(t, wasInCache)

	// Порядок должен быть: key2, key1
	front := cache.(*lruCache).queue.Front()
	require.NotNil(t, front)
	require.Equal(t, "key2", front.Value.(*cacheItem).key)

	// Теперь добавим третий ключ, чтобы проверить удаление старого
	cache.Set("key3", "value3")
	require.Equal(t, 2, cache.(*lruCache).queue.Len())

	// Проверяем, что удален старый (key1)
	_, exists := cache.(*lruCache).items["key1"]
	require.False(t, exists)

	// Проверяем, что order обновлен: key3 на front
	require.Equal(t, "key3", cache.(*lruCache).queue.Front().Value.(*cacheItem).key)
	require.Equal(t, "key2", cache.(*lruCache).queue.Back().Value.(*cacheItem).key)
}

func TestCache_GetRefreshOrder(t *testing.T) { // nolint
	cache := NewCache(2)

	cache.Set("k1", "v1")
	cache.Set("k2", "v2")

	// Делается Get для k1, чтобы обновить его позицию
	val, ok := cache.Get("k1")
	require.True(t, ok)
	require.Equal(t, "v1", val)

	// После этого позиция ключа "k1" должна быть в начале
	require.Equal(t, "k1", cache.(*lruCache).queue.Front().Value.(*cacheItem).key)
	// А ключ "k2" в конце
	require.Equal(t, "k2", cache.(*lruCache).queue.Back().Value.(*cacheItem).key)
}
