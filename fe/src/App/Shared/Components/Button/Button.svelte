<script>
  export let type = 'button';
  export let id = '';
  export let isDisabled = false;

  export let cssLayout = '';
  export let cssFlex = '';
  export let cssSizing = '';
  export let cssBorders = '';
  export let cssSpacing = '';
  export let cssBackgrounds = '';
  export let cssTypography = '';
  export let cssTypographyColor = '';
  export let cssEffects = '';

  export let cssHover = '';
  export let cssFocus = '';
  export let cssDisabled = '';
</script>

<button type="{type}"
        id="{id}"
        disabled={isDisabled}
        class="{cssLayout}
        {cssFlex}
        {cssSizing}
        {cssBorders}
        {cssSpacing}
        {cssBackgrounds}
        {cssTypography}
        {cssTypographyColor}
        {cssEffects}
        {cssHover}
        {cssFocus}
        {cssDisabled}">

  <slot></slot>
</button>
