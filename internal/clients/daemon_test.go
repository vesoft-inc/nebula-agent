package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

func TestDaemon(t *testing.T) {
	assert := assert.New(t)
	rootDir := "/tmp/nebula-install"
	startReq := &pb.StartServiceRequest{
		Role: pb.ServiceRole_STORAGE,
		Dir:  rootDir,
	}
	s := FromStartReq(startReq)
	assert.Equal(s.name, ServiceName_Storaged)
	assert.Equal(s.dir, rootDir)
	d, err := NewDaemon(s)
	assert.Nil(err)
	assert.NotNil(d.Start()) // root dir has no binary

	// bad role format
	startReq = &pb.StartServiceRequest{
		Role: pb.ServiceRole(100),
	}
	s = FromStartReq(startReq)
	assert.Equal(s.name, ServiceName_Unknown)

	// Empty dir
	_, err = NewDaemon(s)
	assert.NotNil(err)
}
