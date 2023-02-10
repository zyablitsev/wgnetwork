<script>
  import { createEventDispatcher, onMount } from 'svelte';

  import { transitionCfg } from '../../util/transition.js';

  export let isChecked = false;
  export let name = '';
  export let id = '';
  export let cssBg = 'bg-gray-200';
  export let cssBgChecked = 'bg-indigo-600';

  let cssX = 'transform translate-x-0';
  let cssXFrom = 'transform translate-x-5';
  let bg = cssBg;
  let bgFrom = cssBgChecked;

  if (isChecked) {
    cssXFrom = 'transform translate-x-0';
    cssX = 'transform translate-x-5';
    bgFrom = cssBg;
    bg = cssBgChecked;
  }

  const dispatch = createEventDispatcher();
  function emitEvent() {
    dispatch('message', {
      isChecked: !isChecked
    });
  }

  onMount(() => {
    let el = document.getElementById(id);
    el.addEventListener('click', emitEvent);

    return () => {
      el.removeEventListener('click', emitEvent);
    }
  });
</script>

<button type="button"
        name="{name}"
        id="{id}"
        class="{bg}
        relative inline-flex
        h-6 w-11
        flex-shrink-0
        cursor-pointer
        rounded-full border-2 border-transparent
        focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        in:transitionCfg|local="{{
            duration: 200,
            base: 'transition-colors ease-in-out duration-200',
            from: bgFrom,
            to: bg,
            isOut: false,
          }}">
  <span class="sr-only">Use setting</span>
  <span class="pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow ring-0"
      in:transitionCfg|local="{{
          duration: 200,
          base: 'transition ease-in-out duration-200',
          from: cssXFrom,
          to: cssX,
          isOut: false,
        }}"
  ></span>
</button>
