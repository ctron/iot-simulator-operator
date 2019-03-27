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

package recon

import (
	"context"
	"time"

	"github.com/ctron/operator-tools/pkg/utils"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"go.uber.org/multierr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type reconcileContext struct {
	ctx     context.Context
	request reconcile.Request
	client  client.Client
	scheme  *runtime.Scheme

	requeue      bool
	requeueAfter time.Duration
	error        error
}

var _ InstallContext = &reconcileContext{}
var _ Reconcile = &reconcileContext{}

func (r *reconcileContext) GetClient() client.Client {
	return r.client
}

func (r *reconcileContext) GetContext() context.Context {
	return r.ctx
}

func (r *reconcileContext) GetScheme() *runtime.Scheme {
	return r.scheme
}

func (r *reconcileContext) GetRequest() reconcile.Request {
	return r.request
}

func NewContext(ctx context.Context, request reconcile.Request, client client.Client, scheme *runtime.Scheme) *reconcileContext {
	return &reconcileContext{
		ctx:     ctx,
		request: request,
		client:  client,
		scheme:  scheme,
	}
}

type InstallContext interface {
	GetClient() client.Client
	GetContext() context.Context
	GetScheme() *runtime.Scheme
	GetRequest() reconcile.Request
}

type Reconcile interface {
	Process(processor Processor)
	ProcessSimple(processor SimpleProcessor)

	NeedRequeue() bool
	Error() error
	Result() (reconcile.Result, error)
}

type Processor func(ctx InstallContext) (reconcile.Result, error)
type SimpleProcessor func(ctx InstallContext) error

func Simple(target SimpleProcessor) Processor {
	return func(ctx InstallContext) (reconcile.Result, error) {
		err := target(ctx)
		return reconcile.Result{}, err
	}
}

func (r *reconcileContext) Process(processor Processor) {
	result, err := processor(r)

	if err != nil {
		r.error = multierr.Append(r.error, err)
	} else {
		r.requeue = r.requeue || result.Requeue
		r.requeueAfter = utils.MaxDuration(r.requeueAfter, result.RequeueAfter)
	}

}

func (r *reconcileContext) ProcessSimple(processor SimpleProcessor) {
	r.Process(Simple(processor))
}

func (r *reconcileContext) NeedRequeue() bool {
	return r.requeue
}

func (r *reconcileContext) Error() error {
	return r.error
}

func (r *reconcileContext) Result() (reconcile.Result, error) {
	return reconcile.Result{
		Requeue:      r.requeue,
		RequeueAfter: r.requeueAfter,
	}, r.error
}
