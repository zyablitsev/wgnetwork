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
  export let device = {
    ip: '',
    label: '',
    wanForward: false,
    user: {
      uuid: '',
      name: ''
    }
  };
  export let wgcfg = {
    wgDeviceInet: '',
    wgDevicePort: 0,
    wgDevicePubKey: '',
    wgDeviceAllowedIPs: [],
    wgDeviceDNS: [],

    wgServerInet: '',
    wgServerIPNet: '',
    wgServerPort: 0,
    wgServerPubKey: '',

    serverWanIP: '',
  };
  let client = cfg.client;

  let isDisabled = false;
  let labelChanged = false;
  let wanForwardChanged = false;
  let wgDevicePubKeyChanged = false;

  function formHandleLabel(event) {
    event.preventDefault;

    device.label = event.target.value;
    labelChanged = true;
    isDisabled = !formValidate(device.label, device.user.uuid);
  }

  function formToggleWanForward(event) {
    event.preventDefault;

    device.wanForward = event.detail.isChecked;
    wanForwardChanged = true;
  }

  function formHandleWGPubKey(event) {
    event.preventDefault;

    wgcfg.wgDevicePubKey = event.target.value;
    wgDevicePubKeyChanged = true;
  }

  const dispatch = createEventDispatcher();

  function handleDeviceEdit(event) {
    event.preventDefault();

    isLoading = true;
    isDisabled = true;

    let params = {'ip': device['ip']}
    if (labelChanged) {
      params['label'] = device.label;
    }
    if (wanForwardChanged) {
      params['wan_forward'] = device.wanForward;
    }
    if (wgDevicePubKeyChanged) {
      params['wg_public_key'] = wgcfg.wgDevicePubKey;
    }
    client.Fetch('manager/device/edit', params, session)
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
    let elInputLabel = document.getElementById('label');
    elInputLabel.addEventListener('input', formHandleLabel);
    let elInputWGPubKey = document.getElementById('wgpubkey');
    elInputWGPubKey.addEventListener('input', formHandleWGPubKey);
    let elBtnCancel = document.getElementById('btn_cancel');
    elBtnCancel.addEventListener('click', moveBack);
    let elBtnSave = document.getElementById('btn_submit');
    elBtnSave.addEventListener('click', handleDeviceEdit);

    return () => {
      elBtnSave.removeEventListener('click', handleDeviceEdit);
      elBtnCancel.removeEventListener('click', moveBack);
      elInputWGPubKey.removeEventListener('input', formHandleWGPubKey);
      elInputLabel.removeEventListener('input', formHandleLabel);
    }
  });

  beforeUpdate(() => {
    isDisabled = !formValidate(device.label, device.user.uuid);
  });
</script>

<form>

  <dl class="divide-y divide-gray-200">

    <div class="grid grid-cols-5 gap-4 py-5 px-6">
      <dt class="text-sm font-medium text-gray-500">ip:</dt>
      <dd class="col-span-4 mt-0 text-sm text-gray-900">
          <div>{device.ip}</div>
      </dd>
    </div>

    <div class="grid grid-cols-5 gap-4 py-5 px-6">
      <dt class="text-sm font-medium text-gray-500">label:</dt>
      <dd class="col-span-4 mt-0 text-sm text-gray-900">
        <input type="text"
               name="label"
               id="label"
               class="block w-full max-w-lg rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:max-w-xs sm:text-sm touch-none"
               value="{device.label}">
      </dd>
    </div>

    <div class="grid grid-cols-5 gap-4 py-5 px-6">
      <dt class="text-sm font-medium text-gray-500">wan forward:</dt>
      <dd class="col-span-4 mt-0 text-sm text-gray-900">
        {#if device.wanForward}
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
               class="block w-full max-w-lg rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:max-w-xs sm:text-sm touch-none"
               value="{wgcfg.wgDevicePubKey}">
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
        save
      {/if}
    </PrimaryButton>

  </div>

</form>
