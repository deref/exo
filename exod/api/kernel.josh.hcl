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

}
