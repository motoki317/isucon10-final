package xsuportal

import (
	"encoding/base64"
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/golang/protobuf/proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/isucon/isucon10-final/webapp/golang/proto/xsuportal/resources"
)

type Notifier struct {
}

var options = webpush.Options{
	Subscriber:      "xsuportal@example.com",
	VAPIDPrivateKey: "8Hhzlr3izBRZ0RWKXraDpk42blfsZbUnVmy1NyniZKk",
	VAPIDPublicKey:  "BC7mQPMOgmwiJYTQyswmsRHLzpGVhd07HYSXtRT9EDgIf-0QMWOzYpGRGdelgT8MmOPxqtjtv4eSexJxJX8oZKc",
}

func (n *Notifier) VAPIDKey() *webpush.Options {
	return &options
}

type notifiableContestant struct {
	ID       string `db:"id"`
	TeamID   int64  `db:"team_id"`
	Endpoint string `db:"endpoint"`
	P256dh   string `db:"p256dh"`
	Auth     string `db:"auth"`
}

func (n *Notifier) NotifyClarificationAnswered(db sqlx.Ext, c *Clarification, updated bool) error {
	var contestants []notifiableContestant

	if c.Disclosed.Valid && c.Disclosed.Bool {
		err := sqlx.Select(
			db,
			&contestants,
			"SELECT c.id AS `id`, c.team_id AS `team_id`, s.endpoint AS `endpoint`, s.p256dh AS `p256dh`, s.auth AS `auth` FROM `contestants` AS c JOIN `push_subscriptions` AS s ON c.id = s.contestant_id WHERE `team_id` IS NOT NULL",
		)
		if err != nil {
			return fmt.Errorf("select all contestants: %w", err)
		}
	} else {
		err := sqlx.Select(
			db,
			&contestants,
			"SELECT `id`, `team_id` FROM `contestants` WHERE `team_id` = ?",
			c.TeamID,
		)
		if err != nil {
			return fmt.Errorf("select contestants(team_id=%v): %w", c.TeamID, err)
		}
	}
	for _, contestant := range contestants {
		notificationPB := &resources.Notification{
			Content: &resources.Notification_ContentClarification{
				ContentClarification: &resources.Notification_ClarificationMessage{
					ClarificationId: c.ID,
					Owned:           c.TeamID == contestant.TeamID,
					Updated:         updated,
				},
			},
		}
		notification, err := n.notify(db, notificationPB, contestant.ID)
		if err != nil {
			return fmt.Errorf("notify: %w", err)
		}
		if n.VAPIDKey() != nil {
			notificationPB.Id = notification.ID
			notificationPB.CreatedAt = timestamppb.New(notification.CreatedAt)
			// TODO: Web Push IIKANJI NI SHITE
			n.notifyProto(contestant, notificationPB)
		}
	}
	return nil
}

func (n *Notifier) NotifyBenchmarkJobFinished(db sqlx.Ext, job *BenchmarkJob) error {
	var contestants []notifiableContestant

	err := sqlx.Select(
		db,
		&contestants,
		"SELECT c.id AS `id`, c.team_id AS `team_id`, s.endpoint AS `endpoint`, s.p256dh AS `p256dh`, s.auth AS `auth` FROM `contestants` AS c JOIN `push_subscriptions` AS s ON c.id = s.contestant_id WHERE `team_id` = ?",
		job.TeamID,
	)
	if err != nil {
		return fmt.Errorf("select contestants(team_id=%v): %w", job.TeamID, err)
	}
	for _, contestant := range contestants {
		notificationPB := &resources.Notification{
			Content: &resources.Notification_ContentBenchmarkJob{
				ContentBenchmarkJob: &resources.Notification_BenchmarkJobMessage{
					BenchmarkJobId: job.ID,
				},
			},
		}
		notification, err := n.notify(db, notificationPB, contestant.ID)
		if err != nil {
			return fmt.Errorf("notify: %w", err)
		}
		if n.VAPIDKey() != nil {
			notificationPB.Id = notification.ID
			notificationPB.CreatedAt = timestamppb.New(notification.CreatedAt)
			// TODO: Web Push IIKANJI NI SHITE
			n.notifyProto(contestant, notificationPB)
		}
	}
	return nil
}

func (n *Notifier) notifyProto(c notifiableContestant, m proto.Message) error {
	res, _ := proto.Marshal(m)
	encRes := base64.StdEncoding.EncodeToString(res)
	var s webpush.Subscription
	s.Endpoint = c.Endpoint
	s.Keys.P256dh = c.P256dh
	s.Keys.Auth = c.Auth
	resp, err := webpush.SendNotification([]byte(encRes), &s, n.VAPIDKey())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (n *Notifier) notify(db sqlx.Ext, notificationPB *resources.Notification, contestantID string) (*Notification, error) {
	return nil, nil
}
