<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { setBodyClass, setHtmlClass } from '../../util/dom';
  import { ValidationError, RPCError } from '../../../lib/rpcapi';
  import { getDevices } from './func.js';
  import { authToken } from '../../state';

  import CardHeading from '../../Shared/Components/CardHeading.svelte';
  import ChevronRight from '../../Shared/Components/Icon/Mini/ChevronRight.svelte';
  import PrimaryButton from '../../Shared/Components/Button/PrimaryButton.svelte';
  import Header from '../../Shared/Components/Header.svelte';
  import Link from '../../Shared/Components/Link.svelte';

  export const params = {};
  export let cfg;
  let client = cfg.client;
  let router = cfg.router;

  setHtmlClass('h-full', 'bg-gray-100');
  setBodyClass('h-full');

  let session = get(authToken);
  let devices = [];

  let isDisabled = true;
  let isLoading = true;

  onMount(() => {
    getDevices(client, session)
      .then(result => {
        devices = result['devices'];
        authToken.set(result['session']);

        isLoading = false;
        isDisabled = false;
      })
      .catch(err => {
        console.error(err);

        // TODO: rework
        if (err instanceof RPCError) {
          isLoading = false;
          isDisabled = true;
        } else {
          isLoading = false;
          isDisabled = false;
        }

        authToken.set('');
      });

    return () => {
    }
  });
</script>

<svelte:head>
  <title>devices</title>
</svelte:head>

<Header {cfg} />

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading} title='devices' description='choose device to view/edit/remove or create a new one'>
    <Link href="{router.ReverseURI('device_create', {})}"
          css='float-right'>
      <PrimaryButton cssSpacing='py-2 px-4'>create</PrimaryButton>
    </Link>
  </CardHeading>

  <div class="p-0">
    {#if isLoading}
    <ul class="divide-y divide-gray-200">
      <li>
        <a href="/" class="block hover:bg-gray-50">
          <div class="flex flex-row items-center justify-between px-6 py-4">
            <div>
              <div class="h-10 w-40 overflow-hidden relative bg-gray-200"></div>
            </div>
          </div>
        </a>
      </li>
    </ul>
    {:else}
      <ul class="divide-y divide-gray-200">
      {#each devices as device}
        <li>
          <Link css="block hover:bg-gray-50"
                href="{router.ReverseURI('device', {'ip': device.ip})}">
            <div class="flex flex-row items-center justify-between px-6 py-4">
              <div>
                <p class="text-sm font-medium text-indigo-600">{device.label} <span class="font-normal text-sm text-gray-500">({device.username})</span></p>
                <p class="text-sm text-gray-500">cidr: {device.ipnetwork}</p>
              </div>
              <div>
                <ChevronRight />
              </div>
            </div>
          </Link>
        </li>
      {/each}
      </ul>
    {/if}
  </div>

</div>
