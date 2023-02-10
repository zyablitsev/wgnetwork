<script>
  import ClipboardDocument from '../../Shared/Components/Icon/Mini/ClipboardDocument.svelte';

  export let key = '';
  export let value = '';
  export let loadingWidth = 40;

  export let isLoading = false;
</script>

<div class="flex flex-row items-start justify-start px-6 py-4">

  <div class="basis-1/5"><p class="text-sm font-medium text-gray-600">{key}:</p></div>

  <div class="basis-4/5 ml-3">
  {#if isLoading}
    <div class="h-5 w-{loadingWidth} overflow-hidden relative bg-gray-200"></div>
  {:else}
    <div class="relative float-right ml-3">
      <div class="absolute right-5 top-0 invisible transition duration-150 opacity-0 ease-in-out">
        <div class="rounded-md bg-green-400 px-2 py-0 float-left mr-2">
          <p class="text-sm font-light text-white">copied</p>
        </div>
      </div>
      <button type="button" data-cb class='float-right text-gray-400 hover:text-gray-500'>
      <ClipboardDocument css='' />
      </button>
    </div>

    <p class="text-sm font-medium text-gray-600 break-all">{value}</p>
  {/if}
  </div>

</div>
