<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import { derived, writable } from 'svelte/store';
  import DirectoryBrowser from './DirectoryBrowser.svelte';

  type NodeIn = {
    name: string;
    children?: NodeIn[];
  };
  type Node = NodeIn & {
    path: string;
    isDirectory: boolean;
    parentPath: string | null;
    children: Node[];
  };

  const nodes = new Map<string, Node>();
  const build = (parentPath: string | null, n: NodeIn): Node => {
    const path = parentPath ? `${parentPath}${n.name}/` : '/';
    const res: Node = {
      path,
      name: n.name,
      isDirectory: n.children != null,
      parentPath,
      children: (n.children ?? []).map((child) => build(path, child)),
    };
    nodes.set(path, res);
    return res;
  };

  build(null, {
    name: '',
    children: [
      {
        name: 'home',
        children: [
          {
            name: 'alice',
            children: [
              { name: 'work', children: [] },
              { name: 'personal', children: [] },
            ],
          },
        ],
      },
    ],
  });
  const homePath = '/home/alice/';

  const directoryPath = writable(homePath);

  const store = {
    ...derived(directoryPath, ($directoryPath) => ({
      ready: true,
      homePath,
      directory: nodes.get($directoryPath),
    })),
    setDirectory: (value: string) => {
      directoryPath.set(value);
    },
  };
</script>

<Meta title="DirectoryBrowser" component={DirectoryBrowser} />

<Template let:args>
  <DirectoryBrowser {...args} />
</Template>

<Story
  name="Default"
  args={{
    store,
  }}
/>
