interface "task-store" {

  method "describe-tasks" {
    output "tasks" "[]TaskDescription" {}
  }

  method "create-task" {
    input "parent-id" "*string" {}
    input "name" "string" {}

    output "id" "string" {}
  }

  method "update-task" {
    input "id" "string" {}
    input "status" "*string" {}
    input "message" "*string" {}
    input "started" "*string" {}
    input "finished" "*string" {}
  }

  method "evict-tasks" {}

}

struct "task-description" {
  field "id" "string" {}
  field "job-id" "string" {
    doc = "ID of root task in this tree."
  }
  field "parent-id" "*string" {}
  field "name" "string" {}
  field "status" "string" {}
  field "message" "string" {}
  field "created" "string" {}
  field "updated" "string" {}
  field "started" "*string" {}
  field "finished" "*string" {}
}
