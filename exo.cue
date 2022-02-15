components: {

  server: {
    type: "process"
    spec: {
      program: "./script/dev-server.sh"
    }
  }
  
  gui: {
    type: "process"
    spec: {
      program: "./script/dev-gui.sh"
    }
  }
  
  storybook: {
    type: "process"
    spec: {
      program: "./script/storybook.sh"
    }
  }

}
