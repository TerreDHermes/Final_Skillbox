package voice_call_handler

import (
	"aggregator"
)

type SliceVoiceCall []aggregator.VoiceCallData

func NewSliceVoiceCall() SliceVoiceCall {
	return make(SliceVoiceCall, 0)
}

func (VoiceCallSlice *SliceVoiceCall) AddVoiceCall(sms aggregator.VoiceCallData) {
	*VoiceCallSlice = append(*VoiceCallSlice, sms)
}
