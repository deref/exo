exo = "0.1"

component "echo" "process" {
  # TODO: Support inline object syntax for specs.
  spec = jsonencode({
    command = "socat"
    arguments = ["TCP4-LISTEN:2000,fork", "EXEC:cat"]
  })
}
