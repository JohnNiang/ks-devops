/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PipelineRunSpec defines the desired state of PipelineRun
type PipelineRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	PipelineRef *PipelineRef `json:"pipelineRef" description:"pipeline reference"`

	Parameters []Parameter `json:"parameters,omitempty" description:"parameters"`
}

// PipelineRunStatus defines the observed state of PipelineRun
type PipelineRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	PipelineRunStatusField `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PipelineRun is the Schema for the pipelineruns API
type PipelineRun struct {
	metav1.TypeMeta   `jsPipelineRunStatusFieldson:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineRunSpec   `json:"spec,omitempty"`
	Status PipelineRunStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineRunList contains a list of PipelineRun
type PipelineRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineRun `json:"items"`
}

type PipelineRunStatusField struct {
	Organization              string              `json:"organization" description:"name of the organization"`
	Id                        string              `json:"id" description:"pipeline id - unique within a pipeline"`
	Pipeline                  string              `json:"pipeline" description:"pipeline name - unique within a pipeline"`
	Name                      string              `json:"name,omitempty" description:"name of the run"`
	Description               string              `json:"description,omitempty" description:"description of the run"`
	ChangeSet                 []ChangeEntry       `json:"changeSet,omitempty" description:"change set of the run"`
	StartTime                 string              `json:"startTime" description:"start time of the run"`
	EnQueueTime               string              `json:"enQueueTime" description:"enqueue time of the run"`
	EndTime                   string              `json:"endTime" description:"end time of the run"`
	DurationInMillis          int64               `json:"durationInMillis" description:"build duration in milli seconds"`
	EstimatedDurationInMillis int64               `json:"estimatedDurationInMillis" description:"Estimated build duration in milli seconds"`
	State                     RunState            `json:"state" description:"the state of the run"`
	Result                    RunResult           `json:"result" description:"the result state of the job (e.g. unstable)"`
	RunSummary                string              `json:"runSummary" description:"build summary"`
	Type                      string              `json:"type" description:"type of the run"`
	ArtifactsZipFile          string              `json:"artifactsZipFile,omitempty" description:"uri of artifacts zip file"`
	Causes                    []map[string]string `json:"causes" description:"cause of the run being created"`
	CauseOfBlockage           string              `json:"causeOfBlockage" description:"cause of what is blocking this run"`
	Replayable                bool                `json:"replayable" description:"if the run will allow a replay"`
}

type ChangeEntry struct {
	Author        *Author  `json:"author"`
	CommitId      string   `json:"commit,omitempty" description:"a human readable display name of the commit number, revision number, and such thing that identifies this entry. null if such a concept doesn't make sense for the implementation. For example, in CVS there's no single identifier for commits. Each file gets a different revision number."`
	Timestamp     string   `json:"timestamp" description:"the timestamp of this commit"`
	Msg           string   `json:"msg" description:"commit message"`
	AffectedPaths []string `json:"affectedPaths" description:"a set of paths in the workspace that was affected by this change"`
	Url           string   `json:"url,omitempty" description:"a browser friendly url to the commit"`
	Issues        []Issue  `json:"issues" description:"related Issues, like GitHub issues or Jira issues"`
	CheckoutCount int      `json:"checkoutCount,omitempty" description:"checkout count of current branch"`
}

type Author struct {
	Id       string `json:"id" description:"the id of the author"`
	FullName string `json:"fullName" description:"the name of the author e.g. John Smith"`
	Email    string `json:"email,omitempty" description:"email address of this author"`
	Avatar   string `json:"avatar" description:"avatar of this author"`
}

type Parameter struct {
	Name  string `json:"name" description:"parameter name"`
	Value string `json:"value" description:"parameter value"`
}

// PipelineRef contains enough information to let you identify the referred resource.
type PipelineRef struct {
	// Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds"
	Kind string
	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string
	// API version of the referent
	// +optional
	APIVersion string
}

type RunState string

const (
	Queued        RunState = "QUEUED"
	Running       RunState = "RUNNING"
	Paused        RunState = "PAUSED"
	Skipped       RunState = "SKIPPED"
	NotBuiltState RunState = "NOT_BUILT"
	Finished      RunState = "FINISHED"
)

type RunResult string

const (
	Success        RunResult = "SUCCESS"
	Unstable       RunResult = "UNSTABLE"
	Failure        RunResult = "FAILURE"
	NotBuiltResult RunResult = "NOT_BUILT"
	Unknown        RunResult = "UNKNOWN"
	Aborted        RunResult = "ABORTED"
)

type Issue struct {
	Id  string `json:"id" description:"issue identifier"`
	Url string `json:"url" description:"issue URL"`
}

func init() {
	SchemeBuilder.Register(&PipelineRun{}, &PipelineRunList{})
}
