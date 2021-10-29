package session

import (
	"fmt"
	"testing"
)

func TestAESTicketProvider_DecodeTicket(t *testing.T) {
	p := NewAESTicketProvider(nil)
	var uin uint64 = 1008611
	ticket, err := p.EncodeTicket(uin)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("ticket", ticket)
	tt, err := p.DecodeTicket(ticket)
	if err != nil {
		t.Error(err)
	}
	if tt.Uin != uin {
		t.Error("uin error", uin, tt.Uin)
	}
	fmt.Println("ticket", tt.Uin, tt.CreateTime)
}
