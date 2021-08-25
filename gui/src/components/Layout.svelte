<script lang="ts">
  import VersionInfo from './VersionInfo.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';
  import FeedbackSVG from './mono/FeedbackSVG.svelte';
  import PreferencesSVG from './mono/PreferencesSVG.svelte';
  import { theme, applyBodyTheme } from '../lib/gui-theme';
  import { onMount } from 'svelte';
  import * as router from 'svelte-spa-router';

  onMount(() => {
    applyBodyTheme($theme);
  });
</script>

<main>
  <nav>
    <header>
      <NavbarButton
        title="Home"
        on:click={() => {
          router.push('#/');
        }}
      >
        {#if import.meta.env.MODE === 'development'}
          <img src="/deref-rounded-icon-dev.png" alt="Deref" height="24px" />
        {:else}
          <img src="/deref-rounded-icon.png" alt="Deref" height="24px" />
        {/if}
      </NavbarButton>
    </header>
    <div class="navbar-wrapper">
      <slot name="navbar" />
    </div>
    <footer>
      <NavbarButton
        title="Preferences"
        on:click={() => {
          router.push('#/preferences');
        }}
      >
        <PreferencesSVG />
      </NavbarButton>
      <NavbarButton
        title="Give feedback on GitHub"
        on:click={() => {
          window.location.href = 'https://github.com/deref/exo/discussions';
        }}
      >
        <FeedbackSVG />
      </NavbarButton>
      <VersionInfo />
    </footer>
  </nav>
  <div class="content-wrapper">
    <slot />
  </div>
</main>

<style>
  :global(body) {
    overscroll-behavior: none;
  }

  main {
    display: grid;
    grid-template-columns: max-content 1fr;
    gap: 1px;
    height: 100vh;
    overflow: hidden;
    background: var(--layout-bg-color);
  }

  nav {
    position: relative;
    width: 48px;
    height: 100vh;
    z-index: 4;
    display: grid;
    grid-template-rows: auto 1fr max-content;
    grid-auto-flow: column;
    gap: 1px;
    background: var(--nav-bg-color);
  }

  .navbar-wrapper {
    width: 48px;
    border-top: 1px solid var(--layout-bg-color);
    border-bottom: 1px solid var(--layout-bg-color);
  }

  nav ::-webkit-scrollbar {
    display: none;
  }

  .content-wrapper {
    position: relative;
    z-index: 2;
    overflow: hidden;
  }
</style>
