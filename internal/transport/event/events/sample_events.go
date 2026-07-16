package events

import entitysample "go-boilerplate-clean/internal/entity/sample"

const ResourceSample = "Sample"

type SampleEvent = Event[entitysample.Sample]

func NewSampleEvent(action Action, before, after *entitysample.Sample) SampleEvent {
	resourceID := ""
	if after != nil {
		resourceID = after.ID
	} else if before != nil {
		resourceID = before.ID
	}
	return NewEvent[entitysample.Sample](resourceID, ResourceSample, action, before, after)
}
