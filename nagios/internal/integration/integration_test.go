package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ulumuri/go-nagios/nagios"
)

func TestArchive(t *testing.T) {
	if len(testInstanceAddress) == 0 {
		t.Skip()
	}

	c, err := nagios.NewClient(http.DefaultClient, testInstanceAddress)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	now := time.Now()
	req := nagios.NotificationListRequest{
		GeneralNotificationRequest: nagios.GeneralNotificationRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
			Count:     10,
			StartTime: 0,
			EndTime:   now.Unix(),
		},
	}

	var list nagios.NotificationList

	if err := c.Query(req, &list); err != nil {
		t.Fatalf("Query: %v", err)
	}

	spew.Dump(list)
}
