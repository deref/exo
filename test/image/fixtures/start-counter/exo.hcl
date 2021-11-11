exo = "0.1"

components {
  container "counter" {
    image = "python:3"
    working_dir = "/my-count"
    command = "expr $(cat count) + 1 > count ; python3 -m http.server"
    volumes = ["e2etest-start-counter:/my-count"]
    ports = ["44225:8000"]
  }

  volume "e2etest-start-counter" {}
}
