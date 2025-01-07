package messageprocessor

import (
	"fmt"
	"sync"
)

// DbConn this only contains info about the dtabase connection and connection Id
// this will be replaced by the connection object from sql package
type DbConn struct {
	ConnectionString string
	ConnectionId     int
}

// countOfConnectionObj is just used to maintan a count
// of number of times the New function is called
var countOfConnectionObj int

// DatabaseMessageProcessor is the actual DB message processor
// that implements the messageProcessor interface
type DatabaseMessageProcessor struct {
	connPool *sync.Pool
}

// this function calls the Get function from the pool
// and returns it back as soon as the work is done
func (processor *DatabaseMessageProcessor) Push(data string) {
	conn := processor.connPool.Get().(*DbConn)
	defer func(c *DbConn) {
		processor.connPool.Put(c)
	}(conn)
	fmt.Printf("used connectionId: %d to insert %s\n", conn.ConnectionId, data)
}

// InitDatabaseMessageProcessor will initialize the DatabaseMessageProcessor having a pool of 2 DbConn obj
// this will be configured based on the number of connect we need at start of app
func InitDatabaseMessageProcessor(connString string) *DatabaseMessageProcessor {
	connPool := &sync.Pool{
		New: func() any {
			countOfConnectionObj++
			fmt.Println("New DB Connection created")
			return &DbConn{
				ConnectionString: connString,
				ConnectionId:     countOfConnectionObj,
			}
		},
	}
	for i := 0; i < 2; i++ {
		// force create a new object and put it in the pool
		connPool.Put(connPool.New())
	}

	return &DatabaseMessageProcessor{
		connPool: connPool,
	}
}
