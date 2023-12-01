package FrontendServiceCache

import (
	"DistributedFileSystem/Metadata"
	"container/list"
)

// Cache Structure definition of LRU Cache, consist of capacity limit, linked list for Cache, Map with filename and FileMetadata element
type Cache struct {
	Capacity        int
	CacheLinkedList *list.List
	CacheMap        map[string]*list.Element
}

// NewCache Initialize a new cache using capacity limit
func (c *Cache) NewCache(capacity int) *Cache {
	return &Cache{
		Capacity:        capacity,
		CacheLinkedList: list.New(),
		CacheMap:        make(map[string]*list.Element),
	}
}

// Get Given filename return corresponding FileMetadata
func (c *Cache) Get(Filename string) (*Metadata.FileMetaData, bool) {
	element, exist := c.CacheMap[Filename]
	if exist {
		c.CacheLinkedList.MoveToFront(element)
		return element.Value.(*Metadata.FileMetaData), true
	}
	return nil, false
}

// Put Check if filename already exist, if exist change the Filename corresponding FileMetadata to this else it will push a new FileMetadata to front
func (c *Cache) Put(FileMetadata *Metadata.FileMetaData) {
	element, exist := c.CacheMap[FileMetadata.FileName]
	if exist {
		c.CacheLinkedList.MoveToFront(element)
		metadata := element.Value.(*Metadata.FileMetaData)
		*metadata = *FileMetadata
		return
	}
	element = c.CacheLinkedList.PushFront(FileMetadata)
	c.CacheMap[FileMetadata.FileName] = element
	if c.CacheLinkedList.Len() > c.Capacity {
		c.Pop()
	}
}

// Pop Delete the last FileMetadata from the linked list and corresponding map
func (c *Cache) Pop() {
	element := c.CacheLinkedList.Back()
	if element != nil {
		c.CacheLinkedList.Remove(element)
		delete(c.CacheMap, element.Value.(*Metadata.FileMetaData).FileName)
	}
}

// Del Given the Filename, delete the FileMetadata with the filename from linked list and corresponding map
func (c *Cache) Del(Filename string) {
	element, exist := c.CacheMap[Filename]
	if !exist {
		return
	}
	c.CacheLinkedList.Remove(element)
	delete(c.CacheMap, element.Value.(*Metadata.FileMetaData).FileName)
}
