// TODO: This is now a misnomer, as it receives inbound events instead of polling for them.
interface "log-collector" {

  doc = "Manages a set of logs. Collects and stores events from them."

  method "clear-events" {
    input "logs" "[]string" {}
  }

  method "describe-logs" {
    input "names" "[]string" {}
    output "logs" "[]LogDescription" {}
  }
  
  method "add-event" {
    input "log" "string" {}
    input "timestamp" "string" {}
    input "message" "string" {}
  }

  method "get-events" {
    doc = "Returns pages of log events for some set of logs. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the log."

    # TODO: Replace this with some filter expression.
    input "logs" "[]string" {}

    input "cursor" "*string" {}
    input "filterStr" "*string" {}
    input "prev" "*int" {}
    input "next" "*int" {}

    output "items" "[]Event" {}
    output "prevCursor" "string" {}
    output "nextCursor" "string" {}
  }

  method "remove-old-events" {}

}

struct "log-description" {
  field "name" "string" {}
  field "last-event-at" "*string" {}
}

struct "event" {
  field "id" "string" {}
  field "log" "string" {}
  field "timestamp" "string" {}
  field "message" "string" {}
}
