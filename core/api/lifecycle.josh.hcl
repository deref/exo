controller "lifecycle" {
  
  method "initialize" {}
  
  method "update" {
    input "new-spec" "string" {}
  }
  
  method "refresh" {}
  
  method "dispose" {
    // TODO: output promise for awaiting synchronous deletes.
  }

}
