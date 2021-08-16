interface "lifecycle" {

  method "initialize" {
    input "exo-labels" "map[string]string" {}
  }

  method "refresh" {}

  method "dispose" {
    // TODO: output promise for awaiting synchronous deletes.
  }

}
