<main>
  <header>
    <img src="/deref-rounded-icon.png" alt="Deref" height="24px" />
    <h1>exo</h1>
  </header>
  <div>
    <slot/>
  </div>
</main>

<style>

:global(body) {
  overscroll-behavior: none;
}

main {
  display: grid;
  grid-auto-flow: row;
  grid-auto-rows: 48px 1fr;
  height: 100vh;
  overflow: hidden;
}

header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  box-shadow: 0px 6px 9px -6px #00000022, 0px 0.25px 0px 1px #00000022;
}

h1 {
  font-size: 20px;
  font-weight: 550;
  margin: 0;
  margin-top: -3px;
}

div {
  position: relative;
  overflow-y: auto;
  overflow-x: hidden;
}

</style>
