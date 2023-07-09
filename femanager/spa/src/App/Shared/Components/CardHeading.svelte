<script>
  import Spinner from './Spinner.svelte';

  export let isLoading = false;
  export let title = '';
  export let description = '';
</script>

<div class="border-b border-gray-200 px-4 py-5 sm:px-6">
  <slot></slot>
  <h3 class="text-lg font-medium leading-6 text-gray-900 inline-block">{title}</h3>
  {#if isLoading}<Spinner css='ml-1 mr-3 text-black inline-block -mt-1' />{/if}
  <p class="mt-1 text-sm text-gray-500">{description}</p>
</div>
