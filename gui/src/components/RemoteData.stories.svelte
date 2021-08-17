<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import { errorResponse, pendingRequest, successResponse } from '../lib/api';
  import RemoteData from './RemoteData.svelte';
</script>

<Meta title="RemoteData" component={RemoteData} />

<Template let:args>
    <RemoteData data={args.data} let:data let:error>
        <span slot="success">Success: {data}</span>
        <span slot="pending">Loading (custom)</span>
        <span slot="error">Error: {error}</span>
    </RemoteData>
</Template>

<Story
  name="Pending Request"
  args={{
    data: pendingRequest(),
  }}
/>

<Story
  name="Successful Response"
  args={{
    data: successResponse("API data"),
  }}
/>

<Story
  name="Error Response"
  args={{
    data: errorResponse("Something went wrong"),
  }}
/>

