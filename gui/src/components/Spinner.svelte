<div />

<style>
  div {
    width: 16px;
    height: 16px;
    border-radius: 100%;
    animation: spin 0.667s infinite linear;
    border: 2px solid var(--spinner-grey);
    border-bottom-color: var(--spinner-grey-light);
    transition: all 0.125s;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
