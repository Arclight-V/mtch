package notification

import (
	"fmt"
	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

// protoContactToUserContact mapper for *notificationservicepb.Contact to *domain.UserContact
func protoContactToUserContact(pc *notificationservicepb.Contact) *domain.UserContact {
	uc := domain.UserContact{
		Channel: domain.Channel(pc.Chanel),
		Value:   pc.Value,
		Meta:    pc.Meta,
	}

	return &uc
}

// protoContactsToUserContacts mapper for []*notificationservicepb.Contact to *domain.UserContacts
func protoContactsToUserContacts(userID string, protoUC []*notificationservicepb.Contact) *domain.UserContacts {
	n := len(protoUC)
	if n == 0 {
		return nil
	}

	fmt.Println(protoUC)

	ucs := domain.UserContacts{
		UserID: userID,
	}

	contacts := make([]*domain.UserContact, n)
	for i := range n {
		contacts[i] = protoContactToUserContact(protoUC[i])
	}

	ucs.Contacts = contacts

	return &ucs
}
