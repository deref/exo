<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import DirectoryBrowser from './DirectoryBrowser.svelte';
  import { nonNull } from '../lib/util';

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

  let nodes = new Map<string, Node>();

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
        name: 'users',
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
  const homeDir = '/users/alice/';

  let directory = nodes.get(homeDir);

  const setDirectory = (path: string) => {
    directory = nodes.get(path);
  };
</script>

<Meta title="DirectoryBrowser" component={DirectoryBrowser} />

<Template let:args>
  <DirectoryBrowser {...args} />
</Template>

<Story
  name="Default"
  args={{
    homePath: '/home/user',
    directory,
    setDirectory,
  }}
/>
