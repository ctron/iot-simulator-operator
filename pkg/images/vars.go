/*******************************************************************************
 * Copyright (c) 2019 Red Hat Inc
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0
 *
 * SPDX-License-Identifier: EPL-2.0
 *******************************************************************************/

package images

import (
	"os"
)

var SimulatorImage string
var ConsoleImage string

func init() {
	base := os.Getenv("IMAGE_BASE")

	if base == "" {
		base = "docker.io/ctron"

	}

	tag := os.Getenv("IMAGE_TAG")
	if tag == "" {
		tag = ":latest"
	}

	SimulatorImage = base + "/iot-hono-simulator" + tag
	ConsoleImage = base + "/iot-simulator-console" + tag
}
