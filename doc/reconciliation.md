# Reconciliation

Design Doc

Unorganized Notes

- scope

  - want to be able to reconcile sub-hierarchies. That is, it should be
    possible to edit a component and reconcile it without affecting other
    components.

    - Is this strictly necessary? Seems like a good idea. If some component is
      in an invalid state, don't want to have to think about it while working
      on getting another component into a valid state.

- concurrency

  - what happens when two reconciliations happen at once?
    - component level
      - should be safe for unrelated components
      - prevented for components that are already reconciling
    - what about at the stack level?
      - reconciling a component should not lock the stack
      - reconciling a stack may cause some components to reconcile, what
        happens if you add a new component and then reconcile the stack again?

- recursive

  - should sub-components be exposed to the user in the UI?
    - against: sub-components are implementation details, the consumer need not
      be aware of them.
    - in favor: some implementation details are known/public, at least for
      read-only consumption. Example: running a docker container should expose
      progress of the image download task.
    - maybe there should be some way to mark children as public vs private?

- iterative

  - after reconciling of some component ends, new outputs are available that
    other components may consume as inputs. need to iterate to a fixed point.

- transitions

  - how should "replace" vs "update" work? consider React-style "keys" as well.
  - idea: components are never replaced when their spec changes. instead, they
    are required to have logic which implements replacement for children as
    needed. keys can be used to allow multiple children to allow temporary
    name conflicts during transitions. one of the components with the same name
    is the "active" one. maybe possible to encode this as sub-sub-components
    without first-class key support.

- hooks

  - initialize
    - would it be better to eliminate this and only have an update hook?
  - update
    - is it necessary to expose the previous spec to this hook? pros/cons?
  - shutdown
    - after marked for disposal, but before delete. Forcible shutdown?

- resources
  - mutated as side-effects of initialize/update/shutdown hooks
