<script>
  import { afterUpdate, onMount } from 'svelte';

  import { ipsetAdd, ipsetRemove } from './func.js';

  import CardHeading from '../../Shared/Components/CardHeading.svelte';
  import LockClosed from '../../Shared/Components/Icon/Outline/LockClosed.svelte';
  import Spinner from '../../Shared/Components/Spinner.svelte';
  import LinkButton from '../../Shared/Components/Button/LinkButton.svelte';
  import PrimaryButton from '../../Shared/Components/Button/PrimaryButton.svelte';

  export let cfg;
  export let session = '';

  export let ipSet = [];
  export let isLoading = false;
  export let isDisabled = false;

  let client = cfg.client;

  function handleIPSetAdd(event) {
    event.preventDefault();
    isLoading = true;
    isDisabled = true;

    let ip = document.getElementById('ip').value;
    ipsetAdd(client, session, ip)
      .then(result => {
        ipSet = result;
        isLoading = false;
        isDisabled = false;
      })
      .catch(err => {
        console.error(err);

        isLoading = false;
        isDisabled = false;
      });
  }

  function handleIPSetRemove(event) {
    event.preventDefault();
    isLoading = true;
    isDisabled = true;

    let ip = this.parentElement.dataset['ip'];
    ipsetRemove(client, session, ip)
      .then(result => {
        ipSet = result;
        isLoading = false;
        isDisabled = false;
      })
      .catch(err => {
        console.error(err);

        isLoading = false;
        isDisabled = false;
      });
  }

  onMount(() => {
    let elBtnSubmit = document.getElementById('btn_ipsetadd');
    elBtnSubmit.addEventListener('click', handleIPSetAdd);

    return () => {
      elBtnSubmit.removeEventListener('click', handleIPSetAdd);
    }
  });

  let ipElems = [];
  afterUpdate(() => {
    if (isLoading) {
      for (let i = 0; i < ipElems.length; i++) {
        let btn = ipElems[i].firstChild;
        btn.removeEventListener('click', handleIPSetRemove);
      }
    } else {
      ipElems = document.querySelectorAll('[data-ip]');
      for (let i = 0; i < ipElems.length; i++) {
        let btn = ipElems[i].firstChild;
        btn.addEventListener('click', handleIPSetRemove);
      }
    }
  });
</script>

<div class="mx-auto max-w-3xl bg-white shadow sm:rounded-lg mt-5">

  <CardHeading {isLoading}
    title='trusted ipset'
    description='list of ip-addresses allowed for remote access to the server via ssh' />

  {#if isLoading}
    <ul class="divide-y divide-gray-200">
      <li>
        <div class="flex flex-row items-center justify-between px-4 py-4">
          <div class="h-5 w-40 overflow-hidden relative bg-gray-200"></div>
        </div>
      </li>
      <li>
        <div class="flex flex-row items-center justify-between px-4 py-4">
          <div class="h-5 w-40 overflow-hidden relative bg-gray-200"></div>
        </div>
      </li>
    </ul>
  {:else}
    {#if ipSet.length > 0}
      <ul class="divide-y divide-gray-200">
      {#each ipSet as ip}
        <li>
          <div class="flex flex-row items-center justify-between px-6 py-4">
            <p class="text-sm font-medium text-gray-600">{ip}</p>
            <div data-ip="{ip}">
              <LinkButton
                      cssFlex='inline-flex items-center'
                      cssBorders='rounded-full border border-gray-300'
                      cssSpacing='px-2.5 py-0.5'
                      cssTypography='text-sm font-medium leading-5'
                      cssTypographyColor='text-gray-700'
                      cssFocus=''
                      isDisabled='{isDisabled}'>
                {#if isLoading}
                  <Spinner css='text-black mx-auto' />
                {:else}
                  remove
                {/if}
              </LinkButton>
            </div>
          </div>
        </li>
      {/each}
      </ul>
    {:else}
      <div class="text-center py-5">
        <LockClosed />
        <h3 class="mt-2 text-sm font-medium text-gray-900">no one allowed</h3>
        <p class="mt-1 text-sm text-gray-500">add one or more ip addresses to whitelist.</p>
      </div>
    {/if}
  {/if}

  <form class="flex flex-row items-center justify-center px-4 pb-5">
    <div class="mx-2">
      <input type="text" name="ip" id="ip" required="true" class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm touch-none" placeholder="10.0.0.1" />
    </div>
    <PrimaryButton id='btn_ipsetadd' type='submit' isDisabled={isDisabled} >
      {#if isLoading}
        <Spinner css='text-white mx-auto' />
      {:else}
        add
      {/if}
    </PrimaryButton>
  </form>

</div>
