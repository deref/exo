<script lang="ts">
  import * as router from 'svelte-spa-router';
  import VersionInfo from './VersionInfo.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';
  import FeedbackSVG from './mono/FeedbackSVG.svelte';
  import PreferencesSVG from './mono/PreferencesSVG.svelte';
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
        <img src="/deref-rounded-icon.png" alt="Deref" height="24px" />
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
      <div class:devmode={import.meta.env.MODE === 'development'}>
        <VersionInfo />
      </div>
    </footer>
  </nav>
  <div>
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
    z-index: 3;
    display: grid;
    grid-template-rows: auto 1fr max-content;
    grid-auto-flow: column;
    gap: 1px;
    background: var(--nav-bg-color);
  }

  .navbar-wrapper {
    width: 48px;
    overflow-y: auto;
    border-top: 1px solid var(--layout-bg-color);
    border-bottom: 1px solid var(--layout-bg-color);
  }

  nav ::-webkit-scrollbar {
    display: none;
  }

  div {
    position: relative;
    z-index: 2;
    overflow: hidden;
  }

  footer {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: 11px;
    width: 100%;
    color: var(--grey-5-color);
    overflow: hidden;
  }

  div.devmode {
    background: var(--dev-mode-bg-color);
    color: var(--dev-mode-color);
  }
</style>
