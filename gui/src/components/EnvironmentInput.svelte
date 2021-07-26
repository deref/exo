<script lang="ts">
  import Textbox from './Textbox.svelte';

  export let name: string;
  export let environment: Record<string, string> = {};

  interface Variable {
    id: number;
    name: string;
    value: string;
  }

  let counter = 0;
  const nextId = () => {
    counter++;
    return counter;
  };

  const newVariable = (name: string, value: string): Variable => ({
    id: nextId(),
    name,
    value,
  });
  let variables: Variable[] = Object.keys(environment).map((variableName) =>
    newVariable(variableName, environment[variableName]),
  );

  const pushBlank = () => {
    variables = [...variables, newVariable('', '')];
  };
  pushBlank();

  const removeVariable = (id: number) => {
    variables = variables.filter((v) => v.id !== id);
  };

  const isBlank = (variable: Variable) =>
    variable.name.trim() === '' && variable.value.trim() === '';

  const isLast = (variable: Variable) =>
    variables.length > 0 && variable.id === variables[variables.length - 1].id;

  $: {
    environment = {};
    for (const variable of variables) {
      const name = variable.name.trim();
      if (variable.name !== '') {
        environment[name] = variable.value.trim();
      }
    }
  }
</script>

<div class="container">
  {#each variables as variable (variable.id)}
    <div class="subcontainer">
      {#each ['name', 'value'] as field}
        <Textbox
          value={variable[field]}
          name={`${name}[${variable.id}].${field}`}
          on:blur={() => {
            if (isBlank(variable) && !isLast(variable)) {
              removeVariable(variable.id);
            }
          }}
          on:input={(e) => {
            const text = e.currentTarget.value.trim();
            const newVariable = { ...variable };
            newVariable[field] = text;
            variables = variables.map((variable) =>
              variable.id === newVariable.id ? newVariable : variable,
            );
            if (!isBlank(newVariable) && isLast(variable)) {
              pushBlank();
            }
          }}
        />
      {/each}
    </div>
  {/each}
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    gap: 24px;
    margin-bottom: 24px;
  }

  .subcontainer {
    display: flex;
    gap: 12px;
  }
</style>
