interface "kernel" {
  
  method "create-workspace" {
    input "root" "string" {}
    
    output "id" "string" {}
  }
  
  method "describe-workspaces" {
    output "workspaces" "[]WorkspaceDescription" {}
  }
  
  method "find-workspace" {
    input "path" "string" {}

    output "id" "*string" {}
  }
  
  method "panic" {
    doc = "Debug method to test what happens when the service panics."

    input "message" "string" {}
  }

}
