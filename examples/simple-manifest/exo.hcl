exo = "0.1"

environment {
  TEST_VARIABLE = "HELLO"
}

components {

  # This is the "long form"
  component "echo" {
    type = "process"
    spec = jsonencode({
      program   = "socat"
      arguments = ["TCP4-LISTEN:2000,fork", "EXEC:cat"]
    })
  }

  # This is a macro that compiles to basically what the above is.
  process "echo-short" {
    program   = "socat"
    arguments = ["TCP4-LISTEN:2000,fork", "EXEC:cat"]
  }

  # These are docker things, they expand to a longform that uses spec =
  # yamlencode(...) like the `jsonencode(` above.
  container "web" {
    build = "."
    ports = [
      "5000:5000"
    ]
  }

  volume "logvolume01" {}

}
