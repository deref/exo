<script lang="ts">
  import Textbox from '../Textbox.svelte';
</script>

<hr />
<div class="group">
  <div class="group-header">
    <h2>Typography</h2>
    <div class="add-new">
      <span>+ Add new</span>
      <ul class="add-new-menu">
        <li>code_ligatures</li>
        <li>font_code</li>
        <li>font_sans</li>
      </ul>
    </div>
  </div>
  <div class="input-text">
    <h3 title="code_ligatures">code_ligatures</h3>
    <Textbox value="false" --input-width="100%" />
  </div>
</div>
<hr />
<div class="group">
  <div class="group-header">
    <h2>Telemetry</h2>
    <div class="add-new">
      <span>+ Add new</span>
      <ul class="add-new-menu">
        <li>disable_telemetry</li>
      </ul>
    </div>
  </div>
  <div class="input-text">
    <h3 title="disable_telemetry">disable_telemetry</h3>
    <Textbox value="true" --input-width="100%" />
  </div>
</div>

<style>
  .group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .add-new {
    position: relative;
  }

  .add-new > span {
    padding: 8px 16px;
    margin: -8px -16px;
    cursor: pointer;
  }

  .add-new:not(:hover) .add-new-menu {
    display: none;
  }

  .add-new:not(:hover) > span {
    color: var(--grey-7-color);
  }

  .add-new-menu {
    position: absolute;
    top: 100%;
    right: -16px;
    margin: 0;
    padding: 8px 0;
    border-radius: 6px;
    background: var(--primary-bg-color);
    box-shadow: var(--button-hover-shadow);
    list-style: none;
  }

  .add-new-menu li {
    padding: 6px 16px;
    font-family: var(--font-mono);
    color: var(--grey-5-color);
    font-size: 15px;
    font-weight: 450;
    cursor: pointer;
  }

  .add-new-menu li:hover {
    background: var(--grey-e-color);
    color: var(--strong-color);
  }

  .add-new-menu li:active {
    background: var(--grey-d-color);
    color: var(--strong-color);
  }

  .group h3 {
    font-family: var(--font-mono);
    font-size: 15px;
    font-weight: 450;
    color: var(--grey-5-color);
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .input-text {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: 40px;
  }
</style>
