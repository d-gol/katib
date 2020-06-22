/*

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

package util

import (
	"fmt"
	"math"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strconv"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var log = logf.Log.WithName("experiment-status-util")

const (
	ExperimentCreatedReason              = "ExperimentCreated"
	ExperimentRunningReason              = "ExperimentRunning"
	ExperimentRestartingReason           = "ExperimentRestarting"
	ExperimentGoalReachedReason          = "ExperimentGoalReached"
	ExperimentMaxTrialsReachedReason     = "ExperimentMaxTrialsReached"
	ExperimentSuggestionEndReachedReason = "ExperimentSuggestionEndReached"
	ExperimentFailedReason               = "ExperimentFailed"
	ExperimentKilledReason               = "ExperimentKilled"
)

func UpdateExperimentStatus(collector *ExperimentsCollector, instance *experimentsv1beta1.Experiment, trials *trialsv1beta1.TrialList) error {

	isObjectiveGoalReached := updateTrialsSummary(instance, trials)

	if !instance.IsCompleted() {
		UpdateExperimentStatusCondition(collector, instance, isObjectiveGoalReached, false)
	}
	return nil

}

func updateTrialsSummary(instance *experimentsv1beta1.Experiment, trials *trialsv1beta1.TrialList) bool {

	var bestTrialValue float64
	sts := &instance.Status
	sts.Trials = 0
	sts.RunningTrialList, sts.PendingTrialList, sts.FailedTrialList, sts.SucceededTrialList, sts.KilledTrialList = nil, nil, nil, nil, nil
	bestTrialIndex := -1
	isObjectiveGoalReached := false
	var objectiveValueGoal float64
	if instance.Spec.Objective.Goal != nil {
		objectiveValueGoal = *instance.Spec.Objective.Goal
	}
	objectiveType := instance.Spec.Objective.Type

	for index, trial := range trials.Items {
		sts.Trials++
		if trial.IsKilled() {
			sts.KilledTrialList = append(sts.KilledTrialList, trial.Name)
		} else if trial.IsFailed() {
			sts.FailedTrialList = append(sts.FailedTrialList, trial.Name)
		} else if trial.IsSucceeded() {
			sts.SucceededTrialList = append(sts.SucceededTrialList, trial.Name)
		} else if trial.IsRunning() {
			sts.RunningTrialList = append(sts.RunningTrialList, trial.Name)
		} else {
			sts.PendingTrialList = append(sts.PendingTrialList, trial.Name)
		}

		objectiveMetricValueStr := getObjectiveMetricValue(trial)
		if objectiveMetricValueStr == nil {
			continue
		}

		objectiveMetricValue, err := strconv.ParseFloat(*objectiveMetricValueStr, 64)
		// For string metrics values best trial is the latest
		if err != nil {
			bestTrialIndex = index
			continue
		}

		//initialize vars to objective metric value of the first trial
		if bestTrialIndex == -1 {
			bestTrialValue = objectiveMetricValue
			bestTrialIndex = index
		}

		if objectiveType == commonv1beta1.ObjectiveTypeMinimize {
			if objectiveMetricValue < bestTrialValue {
				bestTrialValue = objectiveMetricValue
				bestTrialIndex = index
			}
			if instance.Spec.Objective.Goal != nil && bestTrialValue <= objectiveValueGoal {
				isObjectiveGoalReached = true
			}
		} else if objectiveType == commonv1beta1.ObjectiveTypeMaximize {
			if objectiveMetricValue > bestTrialValue {
				bestTrialValue = objectiveMetricValue
				bestTrialIndex = index
			}
			if instance.Spec.Objective.Goal != nil && bestTrialValue >= objectiveValueGoal {
				isObjectiveGoalReached = true
			}
		}
	}

	sts.TrialsRunning = int32(len(sts.RunningTrialList))
	sts.TrialsPending = int32(len(sts.PendingTrialList))
	sts.TrialsSucceeded = int32(len(sts.SucceededTrialList))
	sts.TrialsFailed = int32(len(sts.FailedTrialList))
	sts.TrialsKilled = int32(len(sts.KilledTrialList))

	// if best trial is set
	if bestTrialIndex != -1 {
		bestTrial := trials.Items[bestTrialIndex]

		sts.CurrentOptimalTrial.BestTrialName = bestTrial.Name
		sts.CurrentOptimalTrial.ParameterAssignments = []commonv1beta1.ParameterAssignment{}
		for _, parameterAssigment := range bestTrial.Spec.ParameterAssignments {
			sts.CurrentOptimalTrial.ParameterAssignments = append(sts.CurrentOptimalTrial.ParameterAssignments, parameterAssigment)
		}

		sts.CurrentOptimalTrial.Observation.Metrics = []commonv1beta1.Metric{}
		for _, metric := range bestTrial.Status.Observation.Metrics {
			sts.CurrentOptimalTrial.Observation.Metrics = append(sts.CurrentOptimalTrial.Observation.Metrics, metric)
		}
	}
	return isObjectiveGoalReached
}

func getObjectiveMetricValue(trial trialsv1beta1.Trial) *string {
	if trial.Status.Observation == nil {
		return nil
	}
	var objectiveStrategy commonv1beta1.MetricStrategyType
	objectiveMetricName := trial.Spec.Objective.ObjectiveMetricName
	for _, strategy := range trial.Spec.Objective.MetricStrategies {
		if strategy.Name == objectiveMetricName {
			objectiveStrategy = strategy.Value
			break
		}
	}
	for _, metric := range trial.Status.Observation.Metrics {
		if objectiveMetricName == metric.Name {
			switch objectiveStrategy {
			case commonv1beta1.ExtractByMin:
				if math.IsNaN(metric.Min) {
					return &metric.Latest
				}
				value := fmt.Sprintf("%f", metric.Min)
				return &value
			case commonv1beta1.ExtractByMax:
				if math.IsNaN(metric.Max) {
					return &metric.Latest
				}
				value := fmt.Sprintf("%f", metric.Max)
				return &value
			case commonv1beta1.ExtractByLatest:
				return &metric.Latest
			default:
				return nil
			}
		}
	}
	return nil
}

// UpdateExperimentStatusCondition updates the experiment status.
func UpdateExperimentStatusCondition(collector *ExperimentsCollector, instance *experimentsv1beta1.Experiment, isObjectiveGoalReached bool, getSuggestionDone bool) {
	completedTrialsCount := instance.Status.TrialsSucceeded + instance.Status.TrialsFailed + instance.Status.TrialsKilled
	failedTrialsCount := instance.Status.TrialsFailed
	activeTrialsCount := instance.Status.TrialsPending + instance.Status.TrialsRunning
	now := metav1.Now()

	if isObjectiveGoalReached {
		msg := "Experiment has succeeded because Objective goal has reached"
		instance.MarkExperimentStatusSucceeded(ExperimentGoalReachedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	// First check if MaxFailedTrialCount is reached.
	if (instance.Spec.MaxFailedTrialCount != nil) && (failedTrialsCount > *instance.Spec.MaxFailedTrialCount) {
		msg := "Experiment has failed because max failed count has reached"
		instance.MarkExperimentStatusFailed(ExperimentFailedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsFailedCount(instance.Namespace)
		return
	}

	// Then Check if MaxTrialCount is reached.
	if (instance.Spec.MaxTrialCount != nil) && (completedTrialsCount >= *instance.Spec.MaxTrialCount) {
		msg := "Experiment has succeeded because max trial count has reached"
		instance.MarkExperimentStatusSucceeded(ExperimentMaxTrialsReachedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	if getSuggestionDone && activeTrialsCount == 0 {
		msg := "Experiment has succeeded because suggestion service has reached the end"
		instance.MarkExperimentStatusSucceeded(ExperimentSuggestionEndReachedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	msg := "Experiment is running"
	instance.MarkExperimentStatusRunning(ExperimentRunningReason, msg)
}

func IsCompletedExperimentRestartable(instance *experimentsv1beta1.Experiment) bool {
	if instance.IsSucceeded() && instance.IsCompletedReason(ExperimentMaxTrialsReachedReason) && instance.Spec.ResumePolicy == experimentsv1beta1.LongRunning {
		return true
	}
	return false
}