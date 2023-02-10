<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { setBodyClass, setHtmlClass } from '../../util/dom';
  import { ValidationError, RPCError } from '../../../lib/rpcapi';
  import { getDevice } from './func.js';
  import { authToken, moveBack, moveTo } from '../../state';

  import PrimaryButton from '../../Shared/Components/Button/PrimaryButton.svelte';
  import LinkButton from '../../Shared/Components/Button/LinkButton.svelte';
  import CardHeading from '../../Shared/Components/CardHeading.svelte';
  import Header from '../../Shared/Components/Header.svelte';
  import Spinner from '../../Shared/Components/Spinner.svelte';

  import DeviceFormEdit from './Components/DeviceFormEdit.svelte';

  export let params = {};
  export let cfg;
  let client = cfg.client;
  let router = cfg.router;

  setHtmlClass('h-full', 'bg-gray-100');
  setBodyClass('h-full');

  let session = get(authToken);
  let isLoading = true;

  let device = {
    ip: '',
    label: '',
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
  let title = 'edit device';

  function handleResult(event) {
    event.preventDefault;
    let result = event.detail.result;

    let ip = result['wg_device_inet'].split('/');
    ip = (ip.length > 0) ? ip[0] : '';

    if (ip.length < 1) {
      // TODO: handle as error
      return;
    }

    let path = router.ReverseURI('device', {'ip': ip});
    moveTo(path);
  }

  onMount(() => {
    getDevice(client, session, params['ip'])
      .then(result => {
        device = result['device'];
        wgcfg = result['wgcfg'];
        authToken.set(result['session']);

        isLoading = false;
      })
      .catch(err => {
        console.error(err);

        // TODO: rework
        if (err instanceof RPCError) {
          isLoading = false;
        } else {
          isLoading = false;
        }
      });

    return () => {
    }
  });
</script>

<svelte:head>
  <title>{title}</title>
</svelte:head>

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading} title='edit device' description='the "wan forward" option allows the user device to access the Internet via a server' />

  <div class="p-0">
    <DeviceFormEdit {cfg} {isLoading} {session} {device} {wgcfg} on:message={handleResult} />
  </div>

</div>
