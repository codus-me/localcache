package localcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type localcacheTestSuite struct {
	suite.Suite
	cache *cacheImpl
}

func (suite *localcacheTestSuite) SetupTest() {
	suite.cache = New()
}

func (suite *localcacheTestSuite) TestLocalcacheGetNil() {
	suite.Require().Equal(nil, suite.cache.Get("not exist"))
}
func (suite *localcacheTestSuite) TestLocalcacheSet() {
	suite.cache.Set("mykey", 1)
	suite.Require().Equal(1, suite.cache.hashMap["mykey"].data)
}
func (suite *localcacheTestSuite) TestLocalcacheGetOutdatedData() {
	suite.cache.Set("mykey", 1)
	suite.cache.hashMap["mykey"].createdAt = time.Time{}
	suite.Require().Equal(nil, suite.cache.Get("mykey"))
}

func TestLocalcacheTestSuite(t *testing.T) {
	suite.Run(t, new(localcacheTestSuite))
}
