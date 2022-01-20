## Hierarchy

```mermaid
erDiagram
    Project ||--o{ Workspace : "checkouts"
    Workspace |o--o| Stack : "current"
    Stack ||--o{ Component : "contains"
    Stack }o--|| Project : "instances"
    Component |o--o{ Resource : "owns"
```

## Ownership

The hierarchy diagram shows the ideal, where resources are owned by components.
However, components may also be owned by stacks or projects. Additionally, they
may be ownerless - aka "orphans".
