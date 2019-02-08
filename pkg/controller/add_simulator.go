/*******************************************************************************
 * Copyright (c) 2018 Red Hat Inc
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

package controller

import (
	"github.com/ctron/iot-simulator-operator/pkg/controller/config"
	"github.com/ctron/iot-simulator-operator/pkg/controller/consumer"
	"github.com/ctron/iot-simulator-operator/pkg/controller/producer"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, consumer.Add)
	AddToManagerFuncs = append(AddToManagerFuncs, producer.Add)
	AddToManagerFuncs = append(AddToManagerFuncs, config.Add)
}
