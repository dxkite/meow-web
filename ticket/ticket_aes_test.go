package ticket

import (
	"fmt"
	"testing"
)

func TestAESTicket_Decode(t *testing.T) {
	p := NewAESTicket(nil)
	var uin uint64 = 1008611
	ticket, err := p.Encode(&SessionData{
		Uin:        uin,
		CreateTime: 0,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("ticket", ticket)
	tt, err := p.Decode(ticket)
	if err != nil {
		t.Error(err)
		return
	}
	if tt.Uin != uin {
		t.Error("uin error", uin, tt.Uin)
	}
	fmt.Println("ticket", tt.Uin, tt.CreateTime)
}
