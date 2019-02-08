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

package v1alpha1

type SimulatorComponent interface {
	GetCommonSpec() CommonSpec
}

var _ SimulatorComponent = &SimulatorConsumer{}
var _ SimulatorComponent = &SimulatorProducer{}

func (s *SimulatorConsumer) GetCommonSpec() CommonSpec {
	return s.Spec.CommonSpec
}

func (s *SimulatorProducer) GetCommonSpec() CommonSpec {
	return s.Spec.CommonSpec
}
