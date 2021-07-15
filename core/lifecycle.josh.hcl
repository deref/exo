interface "lifecycle" {
  
  method "initialize" {
    input "id" "string" {}
    input "spec" "string" {}

    output "state" "string" {}
  }
  
  method "update" {
    input "id" "string" {}
    input "old-spec" "string" {}
    input "new-spec" "string" {}
    input "state" "string" {}

    output "state" "string" {}
  }
  
  method "refresh" {
    input "id" "string" {}
    input "spec" "string" {}
    input "state" "string" {}

    output "state" "string" {}
  }
  
  method "dispose" {
    input "id" "string" {}
    input "spec" "string" {}
    input "state" "string" {}
    
    output "state" "string" {}
    // TODO: output promise for awaiting synchronous deletes.
  }

}
