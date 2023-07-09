<script>
  import Button from './Button.svelte';

  export let type = 'submit';
  export let id = '';
  export let isDisabled = false;

  export let cssFlex = '';
  export let cssSizing = '';
  export let cssBorders = 'rounded-md border border-transparent';
  export let cssSpacing = 'py-2 px-2';
  export let cssBackgrounds = 'bg-indigo-600';
  export let cssTypography = 'text-sm font-medium text-white';
  export let cssEffects = 'shadow-sm';

  export let cssHover = 'hover:bg-indigo-700';
  export let cssFocus = 'focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2';
  export let cssDisabled = 'disabled:cursor-not-allowed disabled:border-gray-200 disabled:bg-indigo-500';
</script>

<Button type={type}
        id={id}
        isDisabled={isDisabled}
        cssFlex={cssFlex}
        cssSizing={cssSizing}
        cssBorders={cssBorders}
        cssSpacing={cssSpacing}
        cssBackgrounds={cssBackgrounds}
        cssTypography={cssTypography}
        cssEffects={cssEffects}
        cssHover={cssHover}
        cssFocus={cssFocus}
        cssDisabled={cssDisabled}>
  <slot></slot>
</Button>
