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

package utils

import (
	"time"
)

func MaxInt64(x, y int64) int64 {
	if x < y {
		return y
	} else {
		return x
	}
}

func MaxDuration(x, y time.Duration) time.Duration {
	return time.Duration(MaxInt64(int64(x), int64(y)))
}
