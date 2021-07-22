interface "log-collector" {

  doc = "Manages a set of logs. Collects and stores events from them."

  # TODO: Bulk methods.

  method "add-log" {
    input "name" "string" {}
    input "source" "string" {}
  }

  method "remove-log" {
    input "name" "string" {}
  }

  method "describe-logs" {
    input "names" "[]string" {}
    output "logs" "[]LogDescription" {}
  }

  method "get-events" {
    doc = "Returns pages of log events for some set of logs. If `cursor` is spefied, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the log."

    # TODO: Replace this with some filter expression.
    input "logs" "[]string" {}

    input "cursor" "*string" {}
    input "prev" "*int" {}
    input "next" "*int" {}

    output "events" "[]Event" {}
    output "prevCursor" "string" {}
    output "nextCursor" "string" {}
  }

}

struct "log-description" {
  field "name" "string" {}
  field "source" "string" {}
  field "last-event-at" "*string" {}
}

struct "event" {
  field "id" "string" {}
  field "log" "string" {}
  field "timestamp" "string" {}
  field "message" "string" {}
}
