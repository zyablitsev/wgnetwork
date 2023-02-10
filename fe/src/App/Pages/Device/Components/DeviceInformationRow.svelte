<script>
  import ArchiveBoxArrowDown from '../../../Shared/Components/Icon/Mini/ArchiveBoxArrowDown.svelte';
  import ClipboardDocument from '../../../Shared/Components/Icon/Mini/ClipboardDocument.svelte';

  export let key = '';
  export let value = '';
  export let loadingWidth = 40;

  export let isLoading = false;
  export let clipboard = '';
  export let download = '';
</script>

<div class="grid grid-cols-5 gap-4 py-5 px-6">
  <dt class="text-sm font-medium text-gray-500">{key}:</dt>
    <dd class="col-span-4 mt-0 text-sm text-gray-900">
    {#if isLoading}
      <div class="h-5 w-{loadingWidth} overflow-hidden relative bg-gray-200"></div>
    {:else}
      {#if clipboard.length > 0 && download.length > 0}
      <div class="relative float-right ml-3">
        <div class="absolute right-5 top-0 invisible transition duration-150 opacity-0 ease-in-out">
          <div class="rounded-md bg-green-400 px-2 py-0 float-left mr-2">
            <p class="text-sm font-light text-white">copied</p>
          </div>
        </div>
        <div class='float-right'>
          <button type="button" data-cb="{clipboard}" class='block text-gray-400 hover:text-gray-500'>
          <ClipboardDocument css='' />
          </button>
          <button type="button" data-cb="{download}" class='block mt-5 text-gray-400 hover:text-gray-500'>
          <ArchiveBoxArrowDown css='' />
          </button>
        </div>
      </div>
      {:else if clipboard.length > 0}
      <div class="relative float-right ml-3">
        <div class="absolute right-5 top-0 invisible transition duration-150 opacity-0 ease-in-out">
          <div class="rounded-md bg-green-400 px-2 py-0 float-left mr-2">
            <p class="text-sm font-light text-white">copied</p>
          </div>
        </div>
        <div class='float-right'>
          <button type="button" data-cb="{clipboard}" class='float-right text-gray-400 hover:text-gray-500'>
          <ClipboardDocument css='' />
          </button>
        </div>
      </div>
      {/if}
      <div class="break-all">
      <slot>{value}</slot>
      </div>
    {/if}
    </dd>
</div>
