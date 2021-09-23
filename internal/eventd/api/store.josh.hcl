interface "store" {

  doc = "Database of event organized into streams."

  method "clear-events" {
    input "streams" "[]string" {}
  }

  method "describe-streams" {
    input "names" "[]string" {}
    output "streams" "[]StreamDescription" {}
  }
  
  method "add-event" {
    input "log" "string" {}
    input "timestamp" "string" {}
    input "message" "string" {}
  }

  method "get-events" {
    doc = "Returns pages of events for some set of streams. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the stream."

    # TODO: Replace this with some filter expression.
    input "streams" "[]string" {}

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

struct "stream-description" {
  field "name" "string" {}
  field "last-event-at" "*string" {}
}

struct "event" {
  field "id" "string" {}
  field "stream" "string" {}
  field "timestamp" "string" {}
  field "message" "string" {}
}
