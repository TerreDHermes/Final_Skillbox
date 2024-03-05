package incident_handler

import (
	"aggregator"
	"sort"
)

type SliceIncident []aggregator.IncidentData

func NewSliceIncident() SliceIncident {
	return make(SliceIncident, 0)
}

func (IncidentSlice *SliceIncident) AddIncident(block aggregator.IncidentData) {
	*IncidentSlice = append(*IncidentSlice, block)
}

func (IncidentSlice SliceIncident) SortByStatus() SliceIncident {
	sortedSlice := make(SliceIncident, len(IncidentSlice))
	copy(sortedSlice, IncidentSlice)

	sort.Slice(sortedSlice, func(i, j int) bool {
		// Сначала сравниваем статус
		statusI := sortedSlice[i].Status
		statusJ := sortedSlice[j].Status

		// Если статус "active", то он идет выше в списке
		if statusI == "active" && statusJ != "active" {
			return true
		} else if statusI != "active" && statusJ == "active" {
			return false
		}

		// В остальных случаях порядок не важен, используйте, например, сравнение по Topic
		return sortedSlice[i].Topic < sortedSlice[j].Topic
	})

	return sortedSlice
}
