package localcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type localcacheTestSuite struct {
	suite.Suite
	cache Cache
}

func (suite *localcacheTestSuite) SetupTest() {
	suite.cache = New()
}

func (suite *localcacheTestSuite) TestLocalcacheGetNil() {
	suite.Equal(nil, suite.cache.Get("not exist"))
}
func (suite *localcacheTestSuite) TestLocalcacheSetThenGet() {
	suite.cache.Set("mykey", 1)
	suite.Equal(1, suite.cache.Get("mykey"))
}
func (suite *localcacheTestSuite) TestLocalcacheGetOutdatedData() {
	suite.cache.Set("mykey", 1)
	impl := suite.cache.(*cacheImpl)
	impl.createdAt["mykey"] = &time.Time{}
	suite.Equal(nil, suite.cache.Get("mykey"))
}

func TestLocalcacheTestSuite(t *testing.T) {
	suite.Run(t, new(localcacheTestSuite))
}
