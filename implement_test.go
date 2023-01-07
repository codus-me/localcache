package localcache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type localcacheTestSuite struct {
	suite.Suite
	cache *cacheImpl
}

func (suite *localcacheTestSuite) SetupTest() {
	suite.cache = New().(*cacheImpl)
}

func (suite *localcacheTestSuite) TestLocalcacheGet() {
	suite.cache.hashMap["mykey"] = &cachedData{
		data:      "myvalue",
		createdAt: time.Now(),
	}
	suite.Require().Equal("myvalue", suite.cache.Get("mykey"))
}
func (suite *localcacheTestSuite) TestLocalcacheGetNil() {
	suite.Require().Equal(nil, suite.cache.Get("not exist"))
}
func (suite *localcacheTestSuite) TestLocalcacheGetOutdatedData() {
	suite.cache.Set("mykey", 1)
	suite.cache.hashMap["mykey"].createdAt = time.Time{}
	suite.Require().Equal(nil, suite.cache.Get("mykey"))
}
func (suite *localcacheTestSuite) TestLocalcacheSet() {
	suite.cache.Set("mykey", 1)
	suite.Require().Equal(1, suite.cache.hashMap["mykey"].data)
}
func (suite *localcacheTestSuite) TestLocalcacheConcurrencySet() {
	pause := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-pause
			suite.cache.Set("mykey", "whatever")
		}()
	}
	close(pause)
	wg.Wait()
}
func (suite *localcacheTestSuite) TestLocalcacheFetch() {
	result := suite.cache.Fetch("key", func() interface{} {
		return "value"
	})
	suite.Require().Equal("value", result)
}
func (suite *localcacheTestSuite) TestLocalcacheConcurrencyFetch() {
	var result interface{}
	pause := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-pause
			result = suite.cache.Fetch("key", func() interface{} {
				time.Sleep(time.Second)
				return "value"
			})
		}()
	}
	close(pause)
	wg.Wait()
	suite.Require().Equal("value", result)
}

func TestLocalcacheTestSuite(t *testing.T) {
	suite.Run(t, new(localcacheTestSuite))
}
