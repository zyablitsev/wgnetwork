<script>
  import { createEventDispatcher, beforeUpdate, onMount } from 'svelte';

  import { ValidationError, RPCError } from '../../../../lib/rpcapi';
  import { formValidate } from '../func.js';
  import { moveBack } from '../../../state';

  import PrimaryButton from '../../../Shared/Components/Button/PrimaryButton.svelte';
  import LinkButton from '../../../Shared/Components/Button/LinkButton.svelte';
  import Spinner from '../../../Shared/Components/Spinner.svelte';
  import Toggle from '../../../Shared/Components/Toggle.svelte';

  export let cfg;
  export let isLoading = false;
  export let session = '';
  export let user = {'name': '', 'uuid': ''};
  export let users = [];
  let client = cfg.client;

  let userSelected = {'name': '', 'uuid': ''};
  let isDisabled = true;
  let label = '';
  let wanForward = false;
  let wgPubKey = '';

  if (user['uuid']) {
    userSelected = user;
  }

  function formHandleUserUUID(event) {
    event.preventDefault;

    let uuid = event.target.value;
    for (let i = 0; i < users.length; i++) {
      if (users[i]['uuid'] !== uuid) {
        continue;
      }

      userSelected = users[i];
    }
    isDisabled = !formValidate(label, userSelected.uuid);
  }

  function formHandleLabel(event) {
    event.preventDefault;

    label = event.target.value;
    isDisabled = !formValidate(label, userSelected.uuid);
  }

  function formToggleWanForward(event) {
    event.preventDefault;

    wanForward = event.detail.isChecked;
  }

  function formHandleWGPubKey(event) {
    event.preventDefault;

    wgPubKey = event.target.value;
  }

  const dispatch = createEventDispatcher();

  function handleDeviceCreate(event) {
    event.preventDefault();

    isLoading = true;
    isDisabled = true;

    let params = {
      'user_uuid': userSelected['uuid'],
      'label': label,
      'wan_forward': wanForward};
    if (wgPubKey) {
      params['wg_public_key'] = wgPubKey;
    }
    client.Fetch('manager/device/create', params, session)
      .then(result => {
        isLoading = false;
        isDisabled = false;

        dispatch('message', {
          result: result
        });
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
      });
  }

  onMount(() => {
    let elSelectUserUUID = document.getElementById('user_uuid');
    elSelectUserUUID.addEventListener('change', formHandleUserUUID);
    let elInputLabel = document.getElementById('label');
    elInputLabel.addEventListener('input', formHandleLabel);
    let elInputWGPubKey = document.getElementById('wgpubkey');
    elInputWGPubKey.addEventListener('input', formHandleWGPubKey);
    let elBtnCancel = document.getElementById('btn_cancel');
    elBtnCancel.addEventListener('click', moveBack);
    let elBtnCreate = document.getElementById('btn_submit');
    elBtnCreate.addEventListener('click', handleDeviceCreate);

    return () => {
      elBtnCreate.removeEventListener('click', handleDeviceCreate);
      elBtnCancel.removeEventListener('click', moveBack);
      elInputWGPubKey.removeEventListener('input', formHandleWGPubKey);
      elInputLabel.removeEventListener('input', formHandleLabel);
      elSelectUserUUID.removeEventListener('change', formHandleUserUUID);
    }
  });

  beforeUpdate(() => {
    if (user['uuid']) {
      userSelected = user;
    }
  });
</script>

<form>

  <dl class="divide-y divide-gray-200">

    <div class="grid grid-cols-5 gap-4 py-5 px-6">
      <dt class="text-sm font-medium text-gray-500">user:</dt>
      <dd class="col-span-4 mt-0 text-sm text-gray-900">
        {#if user['uuid'].length > 0}
          <div>{user.name} ({user.uuid})</div>
        {:else}
        <select id="user_uuid" name="user_uuid" class="mt-1 block w-full rounded-md border-gray-300 py-2 pl-3 pr-10 text-base focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm">
          <option value="" selected disabled hidden>choose one</option>
          {#each users as user}
            <option value="{user.uuid}">{user.name} ({user.uuid})</option>
          {/each}
        </select>
        {/if}
      </dd>
    </div>

    <div class="grid grid-cols-5 gap-4 py-5 px-6">
      <dt class="text-sm font-medium text-gray-500">label:</dt>
      <dd class="col-span-4 mt-0 text-sm text-gray-900">
        <input type="text"
               name="label"
               id="label"
               class="block w-full max-w-lg rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:max-w-xs sm:text-sm touch-none">
      </dd>
    </div>

    <div class="grid grid-cols-5 gap-4 py-5 px-6">
      <dt class="text-sm font-medium text-gray-500">wan forward:</dt>
      <dd class="col-span-4 mt-0 text-sm text-gray-900">
        {#if wanForward}
          <Toggle isChecked={true}
                  name='wan_forward' id='wan_forward'
                  on:message={formToggleWanForward} />
        {:else}
          <Toggle isChecked={false}
                  name='wan_forward' id='wan_forward'
                  on:message={formToggleWanForward} />
        {/if}
      </dd>
    </div>

    <div class="grid grid-cols-5 gap-4 py-5 px-6">
      <dt class="text-sm font-medium text-gray-500">wg pubkey:</dt>
      <dd class="col-span-4 mt-0 text-sm text-gray-900">
        <input type="text"
               name="wgpubkey"
               id="wgpubkey"
               class="block w-full max-w-lg rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:max-w-xs sm:text-sm touch-none">
      </dd>
    </div>

  </dl>

  <div class="flex justify-end py-5 px-6">

    <LinkButton type='button'
                id='btn_cancel'
                cssFlex=''
                cssSpacing='py-2 px-4'>
        cancel
    </LinkButton>

    <PrimaryButton type='submit'
                   id='btn_submit'
                   cssFlex='inline-flex justify-center'
                   cssSpacing='py-2 px-4 ml-3'
                   isDisabled={isDisabled}>
      {#if isLoading}
        <div class="inline-flex items-center object-center">
          <Spinner css='text-white' />
          processing
        </div>
      {:else}
        create
      {/if}
    </PrimaryButton>

  </div>

</form>
