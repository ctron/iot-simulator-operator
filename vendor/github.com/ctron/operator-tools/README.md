# Tools for building operators

The idea of this library is to help building operators.

Creating new resources from custom resources might be tricky, as your
operator might not be the only one making changes to your newly created resources.
However, if you are watching your created resources, and reconcile them again, you
might get stuck in an endless loop of reconcile operations. This library tries to help, by
providing some non-destructive editing tools, only updating what is required.

## Reconciling

The main reconcile operation can be split into sub-operations:

~~~go
func (r *ReconcileFooBar) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	
  rec := recon.NewContext(context.TODO(), request, r.client, r.scheme)

  rec.Process(func(ctx InstallContext) (reconcile.Result, error) {
  	return reconcile.Result{}, nil
  })

  rec.ProcessSimple(func(ctx InstallContext) (reconcile.Result, error) {
    return nil
  })

  return rec.Result()

}
~~~

The `Process` function will capture the reconcile result, and aggregate
the requeue flag and time, return this in the call the `Result()`.

## Reconcile Processors

To be written …

## Mixins

Mixins get the chance to manipulate common parts of the newly created
resources, like the objects' metadata. This way is is possible to inject
e.g. an owner reference:

~~~go
rec := recon.NewContext(context.TODO(), request, r.client, r.scheme)

sharedOwnerFn := mixin.ObjectOwner(instance)

// image streams

rec.Process(imagestream.EmptyImageStream("foo", sharedOwnerFn))
rec.Process(imagestream.EmptyImageStream("bar", sharedOwnerFn))
rec.Process(imagestream.EmptyImageStream("baz",
	mixin.Mix(
		sharedOwner,
		mixin.AppendAnnotationMixin("foo", "bar"),
    )
))
~~~