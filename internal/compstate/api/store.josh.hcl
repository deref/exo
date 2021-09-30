interface "store" {
  method "set-state" {
    input "component-id" "string" {}
    input "type" "string" {}
    input "content" "string" {}
    input "tags" "map[string]string" {}
    input "timestamp" "string" {}
    
    output "version" "int" {}
  }
  
  method "get-states" {
    input "component-id" "string" {}
    input "version" "int" {
      doc = "If not specified, begins history with most recent."
    }
    input "history" "int" {
      doc = "Limit of historical states to return per component. Defaults to 1."
    }
    
    output "states" "[]State" {
      doc = "With descending version numbers."
    }
  }
  
  // TODO: describe-components?
  // TODO: garbage collection?
}

struct "state" {
  field "component-id" "string" {}
  field "version" "int" {}
  field "type" "string" {}
  field "content" "string" {}
  field "tags" "map[string]string" {}
  field "timestamp" "string" {}
}
