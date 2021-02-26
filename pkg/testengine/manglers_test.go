package testengine

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/mirbft/pkg/pb/msgs"
	"github.com/IBM/mirbft/pkg/pb/recording"
	"github.com/IBM/mirbft/pkg/pb/state"
)

var _ = Describe("Matchers", func() {
	Describe("initializeMatching", func() {
		It("binds the fields to the underlying implementations", func() {
			mm := &MsgMatching{}
			initializeMatching(mm)
			Expect(mm.FromNode).NotTo(BeNil())
			Expect(mm.FromNodes).NotTo(BeNil())
			Expect(mm.ToNode).NotTo(BeNil())
			Expect(mm.ToNodes).NotTo(BeNil())
			Expect(mm.AtPercent).NotTo(BeNil())

			nmm := mm.AtPercent(30)
			Expect(nmm.Filters).To(HaveLen(1))
			Expect(nmm.Filters[0].apply(5, nil)).To(BeTrue())

			nmm = nmm.ToNodes(3, 5)
			Expect(nmm.Filters).To(HaveLen(2))
			Expect(nmm.Filters[1].apply(0, &recording.Event{
				NodeId: 2,
			})).To(BeFalse())
			Expect(nmm.Filters[1].apply(0, &recording.Event{
				NodeId: 3,
			})).To(BeTrue())
		})
	})

	Describe("MsgTypeMangling", func() {
		It("matches the message type", func() {
			mtm := MatchMsgs().OfTypePreprepare()
			Expect(mtm.Filters).To(HaveLen(2))
			Expect(mtm.Filters[1].apply(0, &recording.Event{
				StateEvent: &state.Event{
					Type: &state.Event_Step{
						Step: &state.EventStep{
							Msg: &msgs.Msg{
								Type: &msgs.Msg_Preprepare{
									Preprepare: &msgs.Preprepare{},
								},
							},
						},
					},
				},
			},
			)).To(BeTrue())
		})

		It("does not match the wrong message type", func() {
			mtm := MatchMsgs().OfTypePreprepare()
			Expect(mtm.Filters).To(HaveLen(2))
			Expect(mtm.Filters[1].apply(0, &recording.Event{
				StateEvent: &state.Event{
					Type: &state.Event_Step{
						Step: &state.EventStep{
							Msg: &msgs.Msg{
								Type: &msgs.Msg_Commit{
									Commit: &msgs.Commit{},
								},
							},
						},
					},
				},
			},
			)).To(BeFalse())
		})
	})
})
