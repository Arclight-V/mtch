package notification

import (
	"testing"

	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

func TestProtoContactsToUserContact(t *testing.T) {
	proto := []*notificationservicepb.Contact{
		{Chanel: notificationservicepb.Channel_ChannelEmail, Value: "a@b.com"},
		{Chanel: notificationservicepb.Channel_ChannelPush, Value: "dev-token"},
	}

	uc := protoContactsToUserContacts("1", proto)

	if uc.UserID != "1" {
		t.Fatalf("uc.UserID = %q, want %q", uc.UserID, "1")
	}

	if len(uc.Contacts) != 2 {
		t.Fatalf("len(uc.Contacts) = %d, want %d", len(uc.Contacts), 2)
	}

	if uc.Contacts[0].Value != "a@b.com" {
		t.Fatal("Wrong mapping")
	}
}
