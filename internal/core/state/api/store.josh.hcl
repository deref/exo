interface "store" {

  method "describe-workspaces" {
    doc = "Returns workspace descriptions."

    input "ids" "[]string" {}

    output "workspaces" "[]WorkspaceDescription" {}
  }

  method "add-workspace" {
    input "id" "string" {}
    input "root" "string" {}
  }

  method "remove-workspace" {
    input "id" "string" {}
  }
  
  method "resolve-workspace" {
    input "ref" "string" {}

    output "id" "*string" {}
  }

  method "resolve" {
    input "workspace-id" "string" {}
    input "refs" "[]string" {}

    output "ids" "[]*string" {}
  }

  method "describe-components" {
    input "workspace-id" "string" {}
    input "refs" "[]string" {}
    input "types" "[]string" {}
    input "include-dependencies" "bool" {}
    input "include-dependents" "bool" {}

    output "components" "[]ComponentDescription" {}
  }

  method "add-component" {
    input "workspace-id" "string" {}
    input "id" "string" {}
    input "name" "string" {}
    input "type" "string" {}
    input "spec" "string" {}
    input "created" "string" {}
    input "depends-on" "[]string" {}
  }

  method "patch-component" {
	  input "id" "string" {}
	  input "state" "string" {}
	  input "initialized" "string" {}
	  input "disposed" "string" {}
	  input "depends-on" "*[]string" {}
  }

  method "remove-component" {
    input "id" "string" {}
  }

}

struct "workspace-description" {
  field "id" "string" {}
  field "root" "string" {}
  field "display-name" "string" {}
}

struct "component-description" {
	field "id" "string" {}
	field "workspace-id" "string" {}
	field "name" "string" {}
	field "type" "string" {}
	field "spec" "string" {}
	field "state" "string" {}
	field "created" "string" {}
	field "initialized" "*string" {}
	field "disposed" "*string" {}
	field "depends-on" "[]string" {}
}
