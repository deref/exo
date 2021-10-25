exo = "0.1"

components {
  process "server" {
    program = "./script/dev-server.sh"
  }

  process "gui" {
    program = "./script/dev-gui.sh"
  }

  process "storybook" {
    program = "./script/storybook.sh"
  }
}
