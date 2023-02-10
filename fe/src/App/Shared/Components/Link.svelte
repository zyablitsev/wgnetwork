<script>
  import { moveTo } from '../../state';

  export let href;
  export let css;

  function redirectTo() {
    moveTo(href);
  }
</script>

<a class={css} href={href} on:click|preventDefault={redirectTo}>
  <slot></slot>
</a>
