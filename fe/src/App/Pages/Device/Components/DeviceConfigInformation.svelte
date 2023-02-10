<script>
  import { afterUpdate, beforeUpdate } from 'svelte';

  import { copyToClipboard, download, excludePrivateNetworks } from './func.js';

  import CardHeading from '../../../Shared/Components/CardHeading.svelte';
  import QRCodeIcon from '../../../Shared/Components/Icon/Outline/QRCode.svelte';
  import QRCode from '../../../Shared/Components/QRCode.svelte';

  import DeviceInformationRow from './DeviceInformationRow.svelte';

  export let wgcfg = {
    wgDeviceInet: '',
    wgDevicePort: '',
    wgDevicePrivKey: undefined,
    wgDevicePubKey: '',
    wgDeviceAllowedIPs: [],

    wgServerInet: '',
    wgServerIPNet: '',
    wgServerPort: 0,
    wgServerPubKey: '',

    serverWanIP: '',
  };
  export let isLoading = false;

  let devicePrivKey = (wgcfg.wgDevicePrivKey && wgcfg.wgDevicePrivKey.length > 0) ? wgcfg.wgDevicePrivKey : '*PLACEHOLDER*';
  let cfg = ''
  let allowedIPs = [];

  beforeUpdate(() => {
    devicePrivKey = (wgcfg.wgDevicePrivKey && wgcfg.wgDevicePrivKey.length > 0) ? wgcfg.wgDevicePrivKey : '*PLACEHOLDER*';
    allowedIPs = excludePrivateNetworks(wgcfg.wgServerIPNet, wgcfg.wgDeviceAllowedIPs);
    cfg = `\
[Interface]
PrivateKey = ${ devicePrivKey }
Address = ${ wgcfg.wgDeviceInet }

[Peer]
PublicKey = ${ wgcfg.wgServerPubKey }
AllowedIPs = ${ allowedIPs.join(', ') }
Endpoint = ${ wgcfg.serverWanIP }:${ wgcfg.wgServerPort }
PersistentKeepalive = 25`
  });

  let cbElems = [];
  let dElems = [];
  afterUpdate(() => {
    if (isLoading) {
      for (let i = 0; i < dElems.length; i++) {
        dElems[i].removeEventListener('click', download);
      }

      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].removeEventListener('click', copyToClipboard);
      }
    } else {
      cbElems = document.querySelectorAll('[data-cb="wgcfginfo"]');
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].addEventListener('click', copyToClipboard);
      }

      dElems = document.querySelectorAll('[data-cb="wgcfgsave"]');
      for (let i = 0; i < dElems.length; i++) {
        dElems[i].addEventListener('click', download);
      }
    }
  });
</script>

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading}
    title='wireguard interface'
    description='client peer parameters' />

  <div class="p-0">
    <dl class="divide-y divide-gray-200">
      {#if wgcfg.wgDevicePrivKey}
      <DeviceInformationRow key='privkey' value={wgcfg.wgDevicePrivKey} {isLoading} clipboard='wgcfginfo' />
      {/if}
      <DeviceInformationRow key='pubkey' value={wgcfg.wgDevicePubKey} {isLoading} clipboard='wgcfginfo' />
      <DeviceInformationRow key='addresses' value={wgcfg.wgDeviceInet} {isLoading} clipboard='wgcfginfo' />
      {#if wgcfg.DevicePort > 0}
      <DeviceInformationRow key='port' value={wgcfg.wgDevicePort} {isLoading} clipboard='wgcfginfo' />
      {/if}
      <DeviceInformationRow key='dns servers' value={wgcfg.wgDeviceDNS.join(', ')} {isLoading} clipboard='wgcfginfo' />
    </dl>
  </div>

</div>
<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading}
    title='wireguard peer'
    description='server peer parameters' />

  <div class="p-0">
    <dl class="divide-y divide-gray-200">
      <DeviceInformationRow key='pubkey' value={wgcfg.wgServerPubKey} {isLoading} clipboard='wgcfginfo' />
      <DeviceInformationRow key='endpoint' value='{wgcfg.serverWanIP}:{wgcfg.wgServerPort}' {isLoading} clipboard='wgcfginfo' />
      <DeviceInformationRow key='allowedips' value={allowedIPs.join(', ')} {isLoading} clipboard='wgcfginfo' />
      <DeviceInformationRow key='persistent keepalive' value='25' {isLoading} clipboard='wgcfginfo' />
    </dl>
  </div>

</div>
<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading}
    title='wireguard config'
    description='plain-text config' />

  <div class="p-0">
    <dl class="divide-y divide-gray-200">
      <DeviceInformationRow key='config' {isLoading} clipboard='wgcfginfo' download='wgcfgsave'>
        <code class="whitespace-pre-line">{cfg}</code>
      </DeviceInformationRow>
      {#if wgcfg.wgDevicePrivKey && wgcfg.wgDevicePrivKey.length > 0}
      <DeviceInformationRow key='qr code' {isLoading}>
        {#if isLoading}
          <QRCodeIcon />
        {:else}
          <QRCode id='qrcode' text={cfg}>
            <div id="qrcode"></div>
          </QRCode>
          <p class="mt-2 text-sm text-gray-500">scan with wireguard app</p>
        {/if}
      </DeviceInformationRow>
      {/if}
    </dl>
  </div>

</div>
