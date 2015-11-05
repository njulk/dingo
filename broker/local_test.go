package broker

import (
	"testing"
	"time"

	"github.com/mission-liao/dingo/common"
	"github.com/mission-liao/dingo/meta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestLocalSend(t *testing.T) {
	ass := assert.New(t)
	ivk := meta.NewDefaultInvoker()

	var (
		b1, b2 Broker
		err    error
	)
	cfg := Default()
	cfg.Local.Bypass(false)
	b1, err = New("local", cfg)
	ass.Nil(err)

	cfg.Local.Bypass(true)
	b2, err = New("local", cfg)
	ass.Nil(err)

	for _, v := range []interface{}{b1, b2} {
		sender, receiver := v.(Producer), v.(Consumer)
		rpt := make(chan Receipt, 10)

		// prepare consumer
		tasks, err := receiver.AddListener(rpt)
		ass.Nil(err)

		// wait for 1 seconds,
		// make sure mux accept that receipt channel
		<-time.After(1 * time.Second)

		// composing a task
		// note: when converting to/from json, only type of float64
		// would be unchanged.
		tk, err := ivk.ComposeTask("test", "param#1", float64(123))
		ass.NotNil(tk)
		ass.Nil(err)

		// send it
		err = sender.Send(tk)
		ass.Nil(err)

		select {
		case expected, ok := <-tasks:
			if !ok {
				ass.Fail("tasks channel is closed")
			} else {
				ass.NotNil(expected)
				if expected != nil {
					ass.True(expected.Equal(tk))
				}
				rpt <- Receipt{
					Id:     expected.GetId(),
					Status: Status.OK,
				}

			}
		}

		// done
		ass.Nil(receiver.(common.Object).Close())
	}
}

func TestLocalConsumeReceipt(t *testing.T) {
	ass := assert.New(t)
	ivk := meta.NewDefaultInvoker()
	rpt := make(chan Receipt, 10)

	cfg := Default()
	cfg.Local.Bypass(false)

	v, err := New("local", cfg)
	ass.Nil(err)
	sender, receiver := v.(Producer), v.(Consumer)
	tasks, err := receiver.AddListener(rpt)
	ass.Nil(err)

	// wait for 1 seconds,
	// make sure mux accept that receipt channel
	<-time.After(1 * time.Second)

	// compose a task
	tk, err := ivk.ComposeTask("test", "test#1")
	ass.NotNil(tk)
	ass.Nil(err)

	err = sender.Send(tk)
	ass.Nil(err)

	select {
	case expected, ok := <-tasks:
		if !ok {
			ass.Fail("tasks channel is closed")
		} else {
			ass.NotNil(expected)

			// There should be an monitored element
			{
				val := v.(*_local)
				_, ok := val.unhandled[expected.GetId()]
				ass.True(ok)
			}

			rpt <- Receipt{
				Id:     expected.GetId(),
				Status: Status.OK,
			}
		}
	}

	// done
	ass.Nil(receiver.(common.Object).Close())
}

//
// generic suite for Brokers
//

type LocalBrokerTestSuite struct {
	BrokerTestSuite
}

func (me *LocalBrokerTestSuite) SetupSuite() {
	var err error

	me.BrokerTestSuite.SetupSuite()
	me._broker, err = New("local", Default())
	me.Nil(err)
}

func (me *LocalBrokerTestSuite) TearDownSuite() {
	me.BrokerTestSuite.TearDownSuite()
}

func TestLocalBrokerSuite(t *testing.T) {
	suite.Run(t, &LocalBrokerTestSuite{})
}
