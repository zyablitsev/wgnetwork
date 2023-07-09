<script>
  import { afterUpdate } from 'svelte';

  import { handleClipboard } from './func.js';

  import CardHeading from '../../Shared/Components/CardHeading.svelte';
  import WireguardInterfaceRow from './WireguardInterfaceRow.svelte';

  export let wgCfg = {
    wanIP: '',
    wgInet: '',
    wgPort: '',
    wgPubKey: ''
  };
  export let isLoading = false;

  let cbElems = [];
  afterUpdate(() => {
    if (isLoading) {
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].removeEventListener('click', handleClipboard);
      }
    } else {
      cbElems = document.querySelectorAll('[data-cb]');
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].addEventListener('click', handleClipboard);
      }
    }
  });
</script>

<div class="mx-auto max-w-3xl bg-white shadow sm:rounded-lg my-5">

  <CardHeading {isLoading}
    title='wireguard interface configuration'
    description='server peer parameters' />

  <div class="p-0">
    <dl class="divide-y divide-gray-200">
      <WireguardInterfaceRow key='ip' value={wgCfg.wanIP} {isLoading} />
      <WireguardInterfaceRow key='port' value={wgCfg.wgPort} {isLoading} loadingWidth=20 />
      <WireguardInterfaceRow key='cidr' value={wgCfg.wgInet} {isLoading} />
      <WireguardInterfaceRow key='public key' value={wgCfg.wgPubKey} {isLoading} loadingWidth=80 />
    </dl>
  </div>

</div>
