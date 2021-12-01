package pipelinerun

import (
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	"kubesphere.io/devops/pkg/api/devops/v1alpha3"
	"kubesphere.io/devops/pkg/apiserver/query"
	resourcesV1alpha3 "kubesphere.io/devops/pkg/models/resources/v1alpha3"
)

type pipelineRunListHandler struct{}

var _ resourcesV1alpha3.ListHandler = pipelineRunListHandler{}

// pipelineRunCompare returns greater than 0, indicating that left PipelineRun is ahead of right.
// pipelineRunCompare returns less than 0, indicating that left comes after right.
// pipelineRunCompare returns equal to 0, indicating that left and right positions remain unchanged.
type pipelineRunCompare func(left, right *v1alpha3.PipelineRun) int

func creationTimeComapre(left, right *v1alpha3.PipelineRun) int {
	if left.CreationTimestamp.IsZero() && right.CreationTimestamp.IsZero() ||
		left.CreationTimestamp.Equal(&right.CreationTimestamp) {
		return 0
	}
	if left.CreationTimestamp.IsZero() {
		return -1
	}
	if right.CreationTimestamp.IsZero() {
		return 1
	}
	if left.CreationTimestamp.Before(&right.CreationTimestamp) {
		return -1
	}
	return 1
}

func startTimeCompare(left, right *v1alpha3.PipelineRun) int {
	if left.Status.StartTime.IsZero() && right.Status.StartTime.IsZero() ||
		left.Status.StartTime.Equal(right.Status.StartTime) {
		return 0
	}
	if left.Status.StartTime.IsZero() {
		return -1
	}
	if right.Status.StartTime.IsZero() {
		return 1
	}
	if left.Status.StartTime.Before(right.Status.StartTime) {
		return -1
	}
	return 1
}

func scmRefNameCompare(left, right *v1alpha3.PipelineRun) int {
	if left.Spec.SCM == nil && right.Spec.SCM == nil {
		return 0
	}
	if left.Spec.SCM == nil {
		return -1
	}
	if right.Spec.SCM == nil {
		return 1
	}
	leftRefName := left.Spec.SCM.RefName
	rightRefName := right.Spec.SCM.RefName
	if leftRefName == rightRefName {
		return 0
	}
	if leftRefName < rightRefName {
		return -1
	}
	return 1
}

func runIDCompare(left, right *v1alpha3.PipelineRun) int {
	leftRunID := left.GetAnnotations()[v1alpha3.JenkinsPipelineRunIDAnnoKey]
	rightRunID := right.GetAnnotations()[v1alpha3.JenkinsPipelineRunIDAnnoKey]
	if leftRunID == "" && rightRunID == "" {
		return 0
	}
	if leftRunID == "" {
		return -1
	}
	if rightRunID == "" {
		return 1
	}
	if leftRunID == rightRunID {
		return 0
	}
	leftID, err := strconv.Atoi(leftRunID)
	if err != nil {
		return -1
	}
	rightID, err := strconv.Atoi(rightRunID)
	if err != nil {
		return 1
	}
	if leftID < rightID {
		return -1
	}
	if leftID > rightID {
		return 1
	}
	return 0
}

func metaNameCompare(left, right *v1alpha3.PipelineRun) int {
	if left.GetName() == right.GetName() {
		return 0
	}
	if left.GetName() < right.GetName() {
		return -1
	}
	return 1
}

func (h pipelineRunListHandler) Comparator() resourcesV1alpha3.CompareFunc {
	return func(left runtime.Object, right runtime.Object, _ query.Field) bool {
		leftPipelineRun, ok := left.(*v1alpha3.PipelineRun)
		if !ok {
			return true
		}
		rightPipelineRun, ok := right.(*v1alpha3.PipelineRun)
		if !ok {
			return true
		}
		if leftPipelineRun == nil || rightPipelineRun == nil {
			return true
		}

		compareChain := []pipelineRunCompare{}
		compareChain = append(compareChain, creationTimeComapre)
		// compareChain = append(compareChain, startTimeCompare)
		// compareChain = append(compareChain, scmRefNameCompare)
		// compareChain = append(compareChain, runIDCompare)
		compareChain = append(compareChain, metaNameCompare)

		comp := 0
		for _, compare := range compareChain {
			comp = compare(leftPipelineRun, rightPipelineRun)
			if comp != 0 {
				break
			}
			// continue next comparison if comp == 0
		}
		// make sure the default order is descending
		return comp >= 0
	}
}

func (h pipelineRunListHandler) Filter() resourcesV1alpha3.FilterFunc {
	return resourcesV1alpha3.DefaultFilter()
}

func (h pipelineRunListHandler) Transformer() resourcesV1alpha3.TransformFunc {
	return resourcesV1alpha3.NoTransformFunc()
}
