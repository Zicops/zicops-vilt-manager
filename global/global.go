package global

import (
	"context"
	"sync"

	"github.com/zicops/zicops-cass-pool/cassandra"
	cry "github.com/zicops/zicops-vilt-manager/lib/crypto"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
	"github.com/zicops/zicops-vilt-manager/lib/sendgrid"
)

var (
	IDP             *identity.IDP
	CTX             context.Context
	Cancel          context.CancelFunc
	CryptSession    *cry.Cryptography
	CassPool        *cassandra.CassandraPool
	SGClient        *sendgrid.ClientSendGrid
	WaitGroupServer sync.WaitGroup
)
