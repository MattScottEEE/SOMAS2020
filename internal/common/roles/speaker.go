package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// Speaker is an interface that is implemented by baseSpeaker but can also be
// optionally implemented by individual islands.
type Speaker interface {
	PayJudge(salary shared.Resources) shared.SpeakerReturnContent
	DecideAgenda(string) shared.SpeakerReturnContent
	DecideVote(string, []shared.ClientID) shared.SpeakerReturnContent
	DecideAnnouncement(string, bool) shared.SpeakerReturnContent
}
