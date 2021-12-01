package pipelinerun

import (
	"reflect"
	"testing"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubesphere.io/devops/pkg/api/devops/v1alpha3"
)

func Test_startTimeCompare(t *testing.T) {
	createPipelineRun := func(name string, startTime *v1.Time) *v1alpha3.PipelineRun {
		return &v1alpha3.PipelineRun{
			ObjectMeta: v1.ObjectMeta{
				Name: name,
			},
			Status: v1alpha3.PipelineRunStatus{
				StartTime: startTime,
			},
		}
	}
	now := v1.Now()
	tomorrow := v1.Time{Time: now.Add(time.Hour * 24)}
	type args struct {
		left  *v1alpha3.PipelineRun
		right *v1alpha3.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "Compare with nil start time",
		args: args{
			left:  createPipelineRun("a", nil),
			right: createPipelineRun("b", nil),
		},
		want: 0,
	}, {
		name: "Compare with same start time",
		args: args{
			left:  createPipelineRun("a", &now),
			right: createPipelineRun("b", &now),
		},
		want: 0,
	}, {
		name: "Compare with where first start time is nil",
		args: args{
			left:  createPipelineRun("a", nil),
			right: createPipelineRun("b", &now),
		},
		want: -1,
	}, {
		name: "Compare with where second start time is nil",
		args: args{
			left:  createPipelineRun("a", &now),
			right: createPipelineRun("b", nil),
		},
		want: 1,
	}, {
		name: "Compare with where first start time is earlier than second",
		args: args{
			left:  createPipelineRun("a", &now),
			right: createPipelineRun("b", &tomorrow),
		},
		want: -1,
	}, {
		name: "Compare with where first start time is later than second",
		args: args{
			left:  createPipelineRun("a", &tomorrow),
			right: createPipelineRun("b", &now),
		},
		want: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startTimeCompare(tt.args.left, tt.args.right); got != tt.want {
				t.Errorf("startTimeCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scmRefNameCompare(t *testing.T) {
	aSCM := &v1alpha3.SCM{
		RefName: "a",
	}
	bSCM := &v1alpha3.SCM{
		RefName: "b",
	}
	createPipelineRun := func(name string, scm *v1alpha3.SCM) *v1alpha3.PipelineRun {
		return &v1alpha3.PipelineRun{
			ObjectMeta: v1.ObjectMeta{
				Name: name,
			},
			Spec: v1alpha3.PipelineRunSpec{
				SCM: scm,
			},
		}
	}
	type args struct {
		left  *v1alpha3.PipelineRun
		right *v1alpha3.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "Compare with nil SCM",
		args: args{
			left:  createPipelineRun("a", nil),
			right: createPipelineRun("b", nil),
		},
		want: 0,
	}, {
		name: "Compare with where first SCM is nil",
		args: args{
			left:  createPipelineRun("a", nil),
			right: createPipelineRun("b", bSCM),
		},
		want: -1,
	}, {
		name: "Compare with where second SCM is nil",
		args: args{
			left:  createPipelineRun("a", aSCM),
			right: createPipelineRun("b", nil),
		},
		want: 1,
	}, {
		name: "Compare with where first SCM is a but second SCM is b",
		args: args{
			left:  createPipelineRun("a", aSCM),
			right: createPipelineRun("b", bSCM),
		},
		want: -1,
	}, {
		name: "Compare with where first SCM is b but second SCM is a",
		args: args{
			left:  createPipelineRun("b", bSCM),
			right: createPipelineRun("a", aSCM),
		},
		want: 1,
	}, {
		name: "Compare with same SCM",
		args: args{
			left:  createPipelineRun("a", aSCM),
			right: createPipelineRun("b", aSCM),
		},
		want: 0,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scmRefNameCompare(tt.args.left, tt.args.right); got != tt.want {
				t.Errorf("scmRefNameCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runIDCompare(t *testing.T) {
	createPipelineRunByAnnotation := func(name string, annotations map[string]string) *v1alpha3.PipelineRun {
		return &v1alpha3.PipelineRun{
			ObjectMeta: v1.ObjectMeta{
				Name:        name,
				Annotations: annotations,
			},
		}
	}

	createPipelineRunByRunID := func(name, runID string) *v1alpha3.PipelineRun {
		return createPipelineRunByAnnotation(name, map[string]string{
			v1alpha3.JenkinsPipelineRunIDAnnoKey: runID,
		})
	}

	type args struct {
		left  *v1alpha3.PipelineRun
		right *v1alpha3.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "Compare with nil annotations",
		args: args{
			left:  createPipelineRunByAnnotation("a", nil),
			right: createPipelineRunByAnnotation("b", nil),
		},
		want: 0,
	}, {
		name: "Compare with where left annotations is nil only",
		args: args{
			left:  createPipelineRunByAnnotation("a", nil),
			right: createPipelineRunByRunID("b", "1"),
		},
		want: -1,
	}, {
		name: "Compare with where right annotations is nil only",
		args: args{
			left:  createPipelineRunByRunID("a", "1"),
			right: createPipelineRunByAnnotation("b", nil),
		},
		want: 1,
	}, {
		name: "Compare with where left run id is smaller than right",
		args: args{
			left:  createPipelineRunByRunID("a", "2"),
			right: createPipelineRunByRunID("b", "11"),
		},
		want: -1,
	}, {
		name: "Compare wih where left run id is greater than right",
		args: args{
			left:  createPipelineRunByRunID("a", "11"),
			right: createPipelineRunByRunID("b", "2"),
		},
		want: 1,
	}, {
		name: "Compare with where left run id is equal to right",
		args: args{
			left:  createPipelineRunByRunID("a", "123"),
			right: createPipelineRunByRunID("b", "123"),
		},
		want: 0,
	}, {
		name: "Compare with where type of left run id is not int",
		args: args{
			left:  createPipelineRunByRunID("a", "invalid_int_id_for_a"),
			right: createPipelineRunByRunID("b", "123"),
		},
		want: -1,
	}, {
		name: "Compare with tpye of right run id is not int",
		args: args{
			left:  createPipelineRunByRunID("a", "123"),
			right: createPipelineRunByRunID("b", "invalid_int_id_for_b"),
		},
		want: 1,
	}, {
		name: "Compare with both types of run id are not int",
		args: args{
			left:  createPipelineRunByRunID("a", "invalid_int_id_for_a"),
			right: createPipelineRunByRunID("b", "invalid_int_id_for_b"),
		},
		want: -1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runIDCompare(tt.args.left, tt.args.right); got != tt.want {
				t.Errorf("runIDCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaNameCompare(t *testing.T) {
	createPipelineRun := func(name string) *v1alpha3.PipelineRun {
		return &v1alpha3.PipelineRun{
			ObjectMeta: v1.ObjectMeta{
				Name: name,
			},
		}
	}
	type args struct {
		left  *v1alpha3.PipelineRun
		right *v1alpha3.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "Compare with where left name is smaller than right",
		args: args{
			left:  createPipelineRun("a"),
			right: createPipelineRun("b"),
		},
		want: -1,
	}, {
		name: "Comapre with where left name is greater than right",
		args: args{
			left:  createPipelineRun("b"),
			right: createPipelineRun("a"),
		},
		want: 1,
	}, {
		name: "Compare with where left name is equal to right",
		args: args{
			left:  createPipelineRun("same"),
			right: createPipelineRun("same"),
		},
		want: 0,
	}, {
		name: "Compare with wher length of left name is longer than right but alphabetical order is lower",
		args: args{
			left:  createPipelineRun("abcdefg"),
			right: createPipelineRun("hijk"),
		},
		want: -1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metaNameCompare(tt.args.left, tt.args.right); got != tt.want {
				t.Errorf("metaNameCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pipelineRunListHandler_Comparator(t *testing.T) {
	now := &v1.Time{Time: time.Now()}
	tomorrow := &v1.Time{Time: now.Add(24 * time.Hour)}
	createPipelineRun := func(creationTime *v1.Time, startTime *v1.Time, refName, id, name string) *v1alpha3.PipelineRun {
		if creationTime == nil {
			creationTime = &v1.Time{}
		}
		pipelineRun := &v1alpha3.PipelineRun{
			ObjectMeta: v1.ObjectMeta{
				Name:              name,
				CreationTimestamp: *creationTime,
				Annotations:       map[string]string{},
			},
			Status: v1alpha3.PipelineRunStatus{
				StartTime: startTime,
			},
		}
		if refName != "" {
			pipelineRun.Spec.SCM = &v1alpha3.SCM{
				RefName: refName,
			}
		}
		if id != "" {
			pipelineRun.GetAnnotations()[v1alpha3.JenkinsPipelineRunIDAnnoKey] = id
		}
		return pipelineRun
	}
	type args struct {
		left  *v1alpha3.PipelineRun
		right *v1alpha3.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "Left creation time is different from right",
		args: args{
			left:  createPipelineRun(now, nil, "", "", ""),
			right: createPipelineRun(tomorrow, nil, "", "", ""),
		},
		want: false,
	}, {
		name: "Left start time is different from right",
		args: args{
			left:  createPipelineRun(now, now, "main", "", ""),
			right: createPipelineRun(now, tomorrow, "dev", "", ""),
		},
		want: false,
	}, {
		name: "Left reference name is different from right",
		args: args{
			left:  createPipelineRun(now, now, "main", "1", ""),
			right: createPipelineRun(now, now, "dev", "2", ""),
		},
		want: true,
	}, {
		name: "Left run ID is different from right",
		args: args{
			left:  createPipelineRun(now, now, "main", "2", "a"),
			right: createPipelineRun(now, now, "main", "1", "b"),
		},
		want: true,
	}, {
		name: "Left name is different from right",
		args: args{
			left:  createPipelineRun(now, now, "main", "1", "b"),
			right: createPipelineRun(now, now, "main", "1", "a"),
		},
		want: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := pipelineRunListHandler{}
			if got := h.Comparator()(tt.args.left, tt.args.right, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pipelineRunListHandler.Comparator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_creationTimeComapre(t *testing.T) {
	createPipelineRun := func(creationTime v1.Time) *v1alpha3.PipelineRun {
		return &v1alpha3.PipelineRun{
			ObjectMeta: v1.ObjectMeta{
				CreationTimestamp: creationTime,
			},
		}
	}
	now := v1.Now()
	tomorrow := v1.Time{Time: now.Add(24 * time.Hour)}

	type args struct {
		left  *v1alpha3.PipelineRun
		right *v1alpha3.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "Compare with zero times",
		args: args{
			left:  createPipelineRun(v1.Time{}),
			right: createPipelineRun(v1.Time{}),
		},
		want: 0,
	}, {
		name: "Compare with same time",
		args: args{
			left:  createPipelineRun(now),
			right: createPipelineRun(now),
		},
		want: 0,
	}, {
		name: "Left time is earlier than right",
		args: args{
			left:  createPipelineRun(now),
			right: createPipelineRun(tomorrow),
		},
		want: -1,
	}, {
		name: "Left time is later than right",
		args: args{
			left:  createPipelineRun(tomorrow),
			right: createPipelineRun(now),
		},
		want: 1,
	}, {
		name: "Left time is zero only",
		args: args{
			left:  createPipelineRun(v1.Time{}),
			right: createPipelineRun(now),
		},
		want: -1,
	}, {
		name: "Right time is zero only",
		args: args{
			left:  createPipelineRun(now),
			right: createPipelineRun(v1.Time{}),
		},
		want: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := creationTimeComapre(tt.args.left, tt.args.right); got != tt.want {
				t.Errorf("creationTimeComapre() = %v, want %v", got, tt.want)
			}
		})
	}
}
