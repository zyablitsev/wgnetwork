<script>
  import { afterUpdate, onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { setBodyClass, setHtmlClass } from '../../util/dom';
  import { ValidationError, RPCError } from '../../../lib/rpcapi';
  import { getDevice } from './func.js';
  import { authToken, moveTo } from '../../state';

  import Header from '../../Shared/Components/Header.svelte';

  import DeviceInformation from './Components/DeviceInformation.svelte';
  import DeviceConfigInformation from './Components/DeviceConfigInformation.svelte';
  import DeviceActionButton from './Components/DeviceActionButton.svelte';
  import DeviceModalRemove from './Components/DeviceModalRemove.svelte';

  export let params = {};
  export let cfg;
  let client = cfg.client;
  let router = cfg.router;

  setHtmlClass('h-full', 'bg-gray-100');
  setBodyClass('h-full');

  let session = get(authToken);
  let isDisabled = true;
  let isLoading = true;
  let showModal = false;

  let device = {
    ip: '',
    label: '',
    wanForward: false,
    user: {
      uuid: '',
      name: ''
    }
  };
  let wgcfg = {
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
  let title = 'device';

  function handleActionButton(event) {
    event.preventDefault;
    if (event.detail.action === 'remove') {
      showModal = true;
    } else if (event.detail.action === 'edit') {
      let path = router.ReverseURI('device_edit', {'ip': device.ip});
      moveTo(path);
    }
  }

  function handleModalButton(event) {
    event.preventDefault;
    if (event.detail.action === 'remove') {
      removeDevice();
    } else if (event.detail.action === 'cancel') {
      showModal = false;
    }
  }

  function removeDevice() {
    isLoading = true;
    isDisabled = true;

    let params = {'ip': device.ip};
    client.Fetch('manager/device/remove', params, session)
      .then(result => {
        isLoading = false;
        isDisabled = false;
        let path = router.ReverseURI('devices');
        moveTo(path);
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
    getDevice(client, session, params['ip'])
      .then(result => {
        device = result['device'];
        wgcfg = result['wgcfg'];
        authToken.set(result['session']);

        isLoading = false;
        isDisabled = false;
      })
      .catch(err => {
        console.error(err);

        // TODO: rework
        if (err instanceof RPCError) {
          isLoading = false;
          isDisabled = false;
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
  <title>{title}</title>
</svelte:head>

<Header {cfg} />

<div class="pb-5">

  <DeviceInformation {cfg} {device} {isLoading} />
  <DeviceActionButton {isDisabled} on:message={handleActionButton} />
  <DeviceConfigInformation {wgcfg} {isLoading} />

  {#if showModal}
    <DeviceModalRemove {isDisabled} on:message={handleModalButton} />
  {/if}

</div>
