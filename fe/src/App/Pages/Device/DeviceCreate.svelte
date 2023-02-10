<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { setBodyClass, setHtmlClass } from '../../util/dom';
  import { RPCError } from '../../../lib/rpcapi';
  import { getUsers } from './func.js';
  import { authToken, moveTo } from '../../state';

  import PrimaryButton from '../../Shared/Components/Button/PrimaryButton.svelte';
  import LinkButton from '../../Shared/Components/Button/LinkButton.svelte';
  import CardHeading from '../../Shared/Components/CardHeading.svelte';
  import Header from '../../Shared/Components/Header.svelte';
  import Spinner from '../../Shared/Components/Spinner.svelte';

  import DeviceFormCreate from './Components/DeviceFormCreate.svelte';
  import DeviceInformation from './Components/DeviceInformation.svelte';
  import DeviceConfigInformation from './Components/DeviceConfigInformation.svelte';

  export const params = {};
  export let query = '';
  export let cfg;
  let client = cfg.client;
  let router = cfg.router;

  setHtmlClass('h-full', 'bg-gray-100');
  setBodyClass('h-full');

  let session = get(authToken);
  let isLoading = true;
  let showInfo = false;

  let users = [];
  let user = {'name': '', 'uuid': ''};
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
    wgDevicePrivKey: '',

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

  function handleResult(event) {
    event.preventDefault;
    let result = event.detail.result;

    let ip = result['wg_device_inet'].split('/');
    ip = (ip.length > 0) ? ip[0] : '';

    if (result['wg_device_privkey']) {
      wgcfg = {
        wgDevicePrivKey: result['wg_device_privkey'],

        wgDeviceInet: result['wg_device_inet'],
        wgDevicePort: result['wg_device_port'],
        wgDevicePubKey: result['wg_device_pubkey'],
        wgDeviceAllowedIPs: result['wg_device_allowed_ips'],
        wgDeviceDNS: ['8.8.8.8', '8.8.4.4'],

        wgServerInet: result['wg_server_inet'],
        wgServerIPNet: result['wg_server_ipnet'],
        wgServerPort: result['wg_server_port'],
        wgServerPubKey: result['wg_server_pubkey'],

        serverWanIP: result['server_wanip'],
      };

      device = {
        ip: ip,
        label: result['label'],
        wanForward: result['wan_forward'],
        user: {
          uuid: result['user_uuid'],
          name: result['user_name']
        }
      }

      showInfo = true;
    } else {
      if (ip.length < 1) {
        // TODO: handle as error
        return;
      }

      let path = router.ReverseURI('device', {'ip': ip});
      moveTo(path);
    }
  }

  onMount(() => {
    getUsers(client, session)
      .then(result => {
        users = result['users'];
        authToken.set(result['session']);

        if (query.length > 0) {
          let searchParams = new URLSearchParams(query);
          if (searchParams.has('user_uuid')) {
            let uuid = searchParams.get('user_uuid');

            for (let i = 0; i < users.length; i++) {
              if (users[i]['uuid'] !== uuid) {
                continue;
              }

              user = users[i];
            }
          }
        }

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
  <title>device create</title>
</svelte:head>

<Header {cfg} />

{#if showInfo}
<div class="pb-5">

  <DeviceInformation {cfg} {device} {isLoading} />
  <DeviceConfigInformation {wgcfg} {isLoading} />

</div>
{:else}
<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading} title='create device' description='select the "wan forward" option if you want to allow the user device to access the internet through the server, if you want to generate keys for the device automatically leave the "wg pubkey" field blank' />

  <div class="p-0">
    <DeviceFormCreate {cfg} {isLoading} {session} {users} {user} on:message={handleResult} />
  </div>

</div>
{/if}
