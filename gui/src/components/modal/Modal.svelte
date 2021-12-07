<!--
  See https://svelte.dev/repl/033e824fad0a4e34907666e7196caec4?version=3.43.1
-->
<script context="module" lang="ts">
  // TODO: Improve `any` types.
  export function bind(Component: any, props = {}) {
    return function ModalComponent(options: any) {
      return new Component({
        ...options,
        props: {
          ...props,
          ...options.props,
        },
      });
    };
  }
</script>

<script lang="ts">
  import * as svelte from 'svelte';
  import { fade } from 'svelte/transition';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  const baseSetContext = svelte.setContext;

  export let show: any = null;

  export let key = 'simple-modal';
  export let closeButton = true;
  export let closeOnEsc = true;
  export let closeOnOuterClick = true;
  export let styleBg = {};
  export let styleWindowWrap = {};
  export let styleWindow = {};
  export let styleContent = {};
  export let styleCloseButton = {};
  export let setContext = baseSetContext;
  export let transitionBg = fade;
  export let transitionBgProps = { duration: 125 };
  export let transitionWindow = transitionBg;
  export let transitionWindowProps = transitionBgProps;
  export let disableFocusTrap = false;

  const defaultState = {
    closeButton,
    closeOnEsc,
    closeOnOuterClick,
    styleBg,
    styleWindowWrap,
    styleWindow,
    styleContent,
    styleCloseButton,
    transitionBg,
    transitionBgProps,
    transitionWindow,
    transitionWindowProps,
    disableFocusTrap,
  };
  let state = { ...defaultState };

  // TODO: Improve `any` types.

  let Component: any = null;

  let background: any;
  let wrap: any;
  let modalWindow: any;
  let scrollY: any;
  let cssBg: any;
  let cssWindowWrap: any;
  let cssWindow: any;
  let cssContent: any;
  let cssCloseButton: any;
  let currentTransitionBg: any;
  let currentTransitionWindow: any;
  let prevBodyPosition: any;
  let prevBodyOverflow: any;
  let prevBodyWidth: any;
  let outerClickTarget: any;

  const camelCaseToDash = (str: string) =>
    str.replace(/([a-zA-Z])(?=[A-Z])/g, '$1-').toLowerCase();

  const toCssString = (props: Record<string, any> | undefined | null) =>
    props
      ? Object.keys(props).reduce(
          (str, key) => `${str}; ${camelCaseToDash(key)}: ${props[key]}`,
          '',
        )
      : '';

  const isFunction = (f: any) => !!(f && f.constructor && f.call && f.apply);

  const updateStyleTransition = () => {
    cssBg = toCssString(
      Object.assign(
        {},
        {
          width: window.innerWidth,
          height: window.innerHeight,
        },
        state.styleBg,
      ),
    );
    cssWindowWrap = toCssString(state.styleWindowWrap);
    cssWindow = toCssString(state.styleWindow);
    cssContent = toCssString(state.styleContent);
    cssCloseButton = toCssString(state.styleCloseButton);
    currentTransitionBg = state.transitionBg;
    currentTransitionWindow = state.transitionWindow;
  };

  // TODO: Improve `any` types.

  const toVoid = () => {};
  let onOpen: any = toVoid;
  let onClose: any = toVoid;
  let onOpened: any = toVoid;
  let onClosed: any = toVoid;

  const open = (
    NewComponent: any,
    newProps = {},
    options = {},
    callback: Record<string, any> = {},
  ) => {
    Component = bind(NewComponent, newProps);
    state = { ...defaultState, ...options };
    updateStyleTransition();
    disableScroll();
    onOpen = (event: Event) => {
      if (callback.onOpen) callback.onOpen(event);
      dispatch('open');
      dispatch('opening'); // Deprecated. Do not use!
    };
    onClose = (event: Event) => {
      if (callback.onClose) callback.onClose(event);
      dispatch('close');
      dispatch('closing'); // Deprecated. Do not use!
    };
    onOpened = (event: Event) => {
      if (callback.onOpened) callback.onOpened(event);
      dispatch('opened');
    };
    onClosed = (event: Event) => {
      if (callback.onClosed) callback.onClosed(event);
      dispatch('closed');
    };
  };

  const close = (callback: Record<string, any> = {}) => {
    if (!Component) return;
    onClose = callback.onClose || onClose;
    onClosed = callback.onClosed || onClosed;
    Component = null;
    enableScroll();
  };

  const handleKeydown = (event: KeyboardEvent) => {
    if (state.closeOnEsc && Component && event.key === 'Escape') {
      event.preventDefault();
      close();
    }

    if (Component && event.key === 'Tab' && !state.disableFocusTrap) {
      // trap focus
      const nodes = modalWindow.querySelectorAll('*');
      const tabbable: any = Array.from(nodes).filter(
        (node: any) => node.tabIndex >= 0,
      );

      let index = tabbable.indexOf(document.activeElement);
      if (index === -1 && event.shiftKey) index = 0;

      index += tabbable.length + (event.shiftKey ? -1 : 1);
      index %= tabbable.length;

      tabbable[index].focus();
      event.preventDefault();
    }
  };

  const handleOuterMousedown = (event: MouseEvent) => {
    if (
      state.closeOnOuterClick &&
      (event.target === background || event.target === wrap)
    )
      outerClickTarget = event.target;
  };

  const handleOuterMouseup = (event: MouseEvent) => {
    if (state.closeOnOuterClick && event.target === outerClickTarget) {
      event.preventDefault();
      close();
    }
  };

  const disableScroll = () => {
    scrollY = window.scrollY;
    prevBodyPosition = document.body.style.position;
    prevBodyOverflow = document.body.style.overflow;
    prevBodyWidth = document.body.style.width;
    document.body.style.position = 'fixed';
    document.body.style.top = `-${scrollY}px`;
    document.body.style.overflow = 'hidden';
    document.body.style.width = '100%';
  };

  const enableScroll = () => {
    document.body.style.position = prevBodyPosition || '';
    document.body.style.top = '';
    document.body.style.overflow = prevBodyOverflow || '';
    document.body.style.width = prevBodyWidth || '';
    window.scrollTo(0, scrollY);
  };

  setContext(key, { open, close });

  let isMounted = false;

  $: {
    if (isMounted) {
      if (isFunction(show)) {
        open(show);
      } else {
        close();
      }
    }
  }

  svelte.onDestroy(() => {
    if (isMounted) close();
  });

  svelte.onMount(() => {
    isMounted = true;
  });
</script>

<svelte:window on:keydown={handleKeydown} />

{#if Component}
  <div
    class="bg"
    on:mousedown={handleOuterMousedown}
    on:mouseup={handleOuterMouseup}
    bind:this={background}
    transition:currentTransitionBg={state.transitionBgProps}
    style={cssBg}
  >
    <div class="window-wrap" bind:this={wrap} style={cssWindowWrap}>
      <div
        class="window"
        role="dialog"
        aria-modal="true"
        bind:this={modalWindow}
        transition:currentTransitionWindow={state.transitionWindowProps}
        on:introstart={onOpen}
        on:outrostart={onClose}
        on:introend={onOpened}
        on:outroend={onClosed}
        style={cssWindow}
      >
        {#if state.closeButton}
          {#if isFunction(state.closeButton)}
            <svelte:component this={state.closeButton} onClose={close} />
          {:else}
            <button on:click={close} class="close" style={cssCloseButton} />
          {/if}
        {/if}
        <div class="content" style={cssContent}>
          <svelte:component this={Component} />
        </div>
      </div>
    </div>
  </div>
{/if}
<slot />

<style>
  * {
    box-sizing: border-box;
  }

  .bg {
    position: fixed;
    z-index: 1000;
    top: 0;
    left: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
    width: 100vw;
    height: 100vh;
    background: var(--window-modal-bg-color);
  }

  .window-wrap {
    position: relative;
    margin: 2rem;
    max-height: 100%;
  }

  .window {
    position: relative;
    width: 40rem;
    max-width: 100%;
    max-height: 100%;
    margin: 2rem auto;
    color: var(--primary-color);
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
    border-radius: 6px;
  }

  .content {
    position: relative;
    padding: 2rem;
    max-height: calc(100vh - 4rem);
    overflow: auto;
  }

  .close {
    display: block;
    box-sizing: border-box;
    position: absolute;
    z-index: 1000;
    top: 1rem;
    right: 1rem;
    margin: 0;
    padding: 0;
    width: 1.5rem;
    height: 1.5rem;
    border: 0;
    color: var(--strong-color);
    border-radius: 1.5rem;
    background: var(--primary-bg-color);
    -webkit-appearance: none;
    outline: none;
  }

  .close:before,
  .close:after {
    content: '';
    display: block;
    box-sizing: border-box;
    position: absolute;
    top: 50%;
    width: 1rem;
    height: 1px;
    background: var(--strong-color);
    transform-origin: center;
  }

  .close:before {
    -webkit-transform: translate(0, -50%) rotate(45deg);
    -moz-transform: translate(0, -50%) rotate(45deg);
    transform: translate(0, -50%) rotate(45deg);
    left: 0.25rem;
  }

  .close:after {
    -webkit-transform: translate(0, -50%) rotate(-45deg);
    -moz-transform: translate(0, -50%) rotate(-45deg);
    transform: translate(0, -50%) rotate(-45deg);
    left: 0.25rem;
  }

  .close:hover {
    background-color: var(--grey-d-color);
  }

  .close:focus {
    border-color: var(--link-color);
    box-shadow: 0 0 0 2px var(--link-color);
  }

  .close:active {
    transform: scale(0.9);
  }
</style>
