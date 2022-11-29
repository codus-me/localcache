package localcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type LocalcacheTestSuite struct {
	suite.Suite
	cache                         Cache
	VariableThatShouldStartAtFive int
}

func (suite *LocalcacheTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
	suite.cache = New()
}

func (suite *LocalcacheTestSuite) TestLocalcacheGetNil() {
	suite.Equal(nil, suite.cache.Get("not exist"))
}
func (suite *LocalcacheTestSuite) TestLocalcacheSetThenGet() {
	suite.cache.Set("mykey", 1)
	suite.Equal(1, suite.cache.Get("mykey"))
}
func (suite *LocalcacheTestSuite) TestLocalcacheGetOutdatedData() {
	suite.cache.Set("mykey", 1)
	impl := suite.cache.(*cacheImpl)
	impl.createdAt["mykey"] = &time.Time{}
	suite.Equal(nil, suite.cache.Get("mykey"))
}

func TestLocalcacheTestSuite(t *testing.T) {
	suite.Run(t, new(LocalcacheTestSuite))
}
