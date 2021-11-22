exo = "0.1"
components {
  container "t0" {
    image = "bash"
    command  = "sleep infinity"
  }
  container "t1" {
    image = "bash"
    command  = "sleep infinity"
  }
}
