<script lang="ts">
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import IconButton from '../components/IconButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { theme, themeOptions } from '../lib/theme';
  import type { Preferences } from '../lib/preferences';
  import { preferences } from '../lib/preferences';
  import { onMount } from 'svelte';
  import InputQuantity from '../components/form/InputQuantity.svelte';
  import InputSelection from '../components/form/InputSelection.svelte';

  type PreferenceName = keyof Preferences;

  interface QuantityPreference {
    name: PreferenceName;
    type: 'quantity';
    units?: string[];
  }

  interface SelectPreference {
    name: PreferenceName;
    type: 'select';
    options: string[];
  }

  type Preference = QuantityPreference | SelectPreference;

  interface PreferenceGroup {
    title: string;
    preferences: Preference[];
  }

  const groups: PreferenceGroup[] = [
    {
      title: 'Typography',
      preferences: [
        {
          name: 'main-font-size',
          type: 'quantity',
          units: ['px', 'em', 'rem', 'ex', '%'],
        },
        {
          name: 'log-font-size',
          type: 'quantity',
          units: ['px', 'em', 'rem', 'ex', '%'],
        },
        {
          name: 'ligatures-logs',
          type: 'select',
          options: ['none', 'normal', 'common-ligatures'],
        },
        {
          name: 'ligatures-code',
          type: 'select',
          options: ['none', 'normal', 'common-ligatures'],
        },
      ],
    },
  ];

  let dirtyPrefs: Preferences = {};

  onMount(() => {
    dirtyPrefs = $preferences;
  });

  let testBindQVal = '215px';
</script>

<Layout>
  <CenterFormPanel title="Preferences">
    <div slot="actions">
      <IconButton
        glyph="Reset"
        tooltip="Reset to defaults"
        on:click={() => {
          theme.apply('auto');
          preferences.reset();
          setTimeout(() => {
            dirtyPrefs = $preferences;
          }, 50);
        }}
      />
    </div>
    <div>
      <div class="group">
        <div class="group-header">
          <h2>{'Theme & GUI'}</h2>
        </div>
        <div class="button-row">
          {#each themeOptions as option}
            <Button
              on:click={() => theme.apply(option)}
              inset={$theme === option}
            >
              <code>{option}</code>
            </Button>
          {/each}
        </div>
      </div>

      {#each groups as group}
        <div class="group">
          <div class="group-header">
            <h2>{group.title}</h2>
          </div>
          {#each group.preferences as preference}
            <div class="input-row">
              <code>{preference.name}</code>
              {#if preference.type === 'quantity'}
                <InputQuantity
                  bind:value={dirtyPrefs[preference.name]}
                  on:input={() => preferences.apply({ ...dirtyPrefs })}
                  unitOptions={preference.units}
                />
              {:else if preference.type === 'select'}
                <InputSelection
                  bind:value={dirtyPrefs[preference.name]}
                  on:input={() => preferences.apply({ ...dirtyPrefs })}
                  options={preference.options}
                />
              {/if}
            </div>
          {/each}
        </div>
      {/each}
    </div>
    <hr />
    <div>
      <p>
        testBindQVal = {testBindQVal}
      </p>
      <InputQuantity bind:value={testBindQVal} unitOptions={['px', 'rem']} />
    </div>
  </CenterFormPanel>
</Layout>

<style>
  .input-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    align-items: center;
    margin-bottom: 12px;
  }

  .group:not(:last-child) {
    margin-bottom: 40px;
  }

  .group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .button-row {
    display: grid;
    grid-auto-flow: column;
    gap: 12px;
  }
</style>
