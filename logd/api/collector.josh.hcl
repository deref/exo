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
    doc = "Paginates events. Inputs before and after are mutually exclusive."
    input "logs" "[]string" {}
	  input "cursor" "string" {}
	  input "since" "string" {}
	  input "limit" "int" {}
    
    output "events" "[]Event" {}
    output "cursor" "string" {}
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
