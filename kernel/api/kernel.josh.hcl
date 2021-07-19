interface "kernel" {
  
  method "create-workspace" {
    input "name" "string" {}
    input "path" "string" {}
  }
  
  method "forget-workspace" {
    input "ref" "string" {}
  }
  
  method "describe-workspaces" {
    output "workspaces" "[]WorkspaceDescription" {}
  }
  
  method "find-workspace" {
    input "path" "string" {}

    output "id" "*string" {}
  }

}

struct "workspace-description" {
  field "id" "string" {}
  field "name" "string" {}
  field "root" "string" {}
  field "project" "string" {}
}
