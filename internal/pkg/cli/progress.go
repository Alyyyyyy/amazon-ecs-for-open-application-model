// Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import termprogress "github.com/awslabs/amazon-ecs-for-open-application-model/internal/pkg/term/progress"

// progress is the interface to inform the user that a long operation is taking place.
type progress interface {
	// Start starts displaying progress with a label.
	Start(label string)
	// Stop ends displaying progress with a label.
	Stop(label string)
	// Events writes additional information in between the start and stop stages.
	Events([]termprogress.TabRow)
}
