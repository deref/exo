interface "process" {
  method "start" {
    input "id" "string" {}
    input "spec" "string" {}
    input "state" "string" {}
    output "state" "string" {}
  }
  method "stop" {
    input "id" "string" {}
    input "spec" "string" {}
    input "state" "string" {}
    output "state" "string" {}
  }
}
