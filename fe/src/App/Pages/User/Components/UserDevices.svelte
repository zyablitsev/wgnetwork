<script>
  import { afterUpdate, onMount } from 'svelte';

  import CardHeading from '../../../Shared/Components/CardHeading.svelte';
  import ChevronRight from '../../../Shared/Components/Icon/Mini/ChevronRight.svelte';
  import ComputerDesktop from '../../../Shared/Components/Icon/Outline/ComputerDesktop.svelte';
  import PrimaryButton from '../../../Shared/Components/Button/PrimaryButton.svelte';
  import Link from '../../../Shared/Components/Link.svelte';

  export let cfg;
  let router = cfg.router;

  export let devices = [];
  export let userUUID = '';
  export let isLoading = false;

  let deviceCreateURI = ''
  if (userUUID.length > 0) {
    let queryParams = new URLSearchParams();
    queryParams.append('user_uuid', userUUID);
    deviceCreateURI = router.ReverseURI('device_create', {}, queryParams);
  } else {
    deviceCreateURI = router.ReverseURI('device_create', {});
  }

  afterUpdate(() => {
    if (userUUID.length == 0) {
      return;
    }

    let queryParams = new URLSearchParams();
    queryParams.append('user_uuid', userUUID);
    deviceCreateURI = router.ReverseURI('device_create', {}, queryParams);
  });
</script>

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading} title='user devices' description='choose device to view/edit/remove'>
    <Link href="{deviceCreateURI}"
          css='float-right'>
      <PrimaryButton cssSpacing='py-2 px-4'>add</PrimaryButton>
    </Link>
  </CardHeading>

  <div class="p-0">
    {#if isLoading}
    <ul class="divide-y divide-gray-200">
      <li>
        <div class="block hover:bg-gray-50">
          <div class="flex flex-row items-center justify-between px-6 py-4">
            <div>
              <div class="h-10 w-40 overflow-hidden relative bg-gray-200"></div>
            </div>
          </div>
        </div>
      </li>
    </ul>
    {:else}
      {#if devices.length > 0}
      <ul class="divide-y divide-gray-200">
      {#each devices as device}
        <li>
          <Link css="block hover:bg-gray-50"
                href="{router.ReverseURI('device', {'ip': device.ip})}">
            <div class="flex flex-row items-center justify-between px-6 py-4">
              <div class="">
                <p class="text-sm font-medium text-indigo-600">{device.label}</p>
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
      {:else}
        <div class="text-center py-5">
          <ComputerDesktop />
          <h3 class="mt-2 text-sm font-medium text-gray-900">no devices yet</h3>
          <p class="mt-1 text-sm text-gray-500">add one or more user devices.</p>
        </div>
      {/if}
    {/if}
  </div>

</div>
