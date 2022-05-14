// TODO: Explicitly include manifest.cue?

#EnvironmentValue: { [string]: #EnvironmentVariableValue }

// TODO: What, if anything, should manifests be allowed to say about clusters?
#Cluster: {
  id: string
  name: string
  environment: #EnvironmentValue
}

// TODO: How do manifests specify adoption of resources?
#Resource: {
  type: string
  id: string
  iri: string
}

#Component: {
  id: string
  spec: _
  model: spec // Effectively spec & state.
  environment: #EnvironmentValue
  resources: #ResourcesById
  children: #ComponentsByName
}

#ComponentsById: {
  [Id=string]: #Component & {
    id: Id
  }
}

#ResourcesById: {
  [Id=string]: #Resource & {
    id: Id
  }
}

#Stack: {
  id: string
  name: string
  environment: #EnvironmentValue
  components: #ComponentsByName
  detachedResources: #ResourcesById
}

// TODO: Move to os package, prefix type string, etc.
#Daemon: #Component & {
  type: "daemon"
  spec: {
    program: string
    arguments: [...string] | *[]
    "environment": environment
  }
  environment: _
}

// TODO: Move to os package, prefix type string, etc.
#Process: #Component & {
  type: "process"
  spec: {
    program: string
    arguments: [...string] | *[]
    "environment": environment
  }
  environment: _
}
