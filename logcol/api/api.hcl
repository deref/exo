interface "log-collector" {

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
    input "logs" "[]string" {}
    input "before" "string" {}
    input "after" "string" {}
    output "events" "[]Event" {}
  }

}

struct "log-description" {
  field "name" "string" {}
  field "source" "string" {}
  field "last-event-at" "*string" {}
}

struct "event" {
  field "log" "string" {}
  field "sid" "string" {}
  field "timestamp" "string" {}
  field "message" "string" {}
}
