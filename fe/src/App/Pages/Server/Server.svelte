<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { setBodyClass, setHtmlClass } from '../../util/dom.js';
  import { ValidationError, RPCError } from '../../../lib/rpcapi.js';
  import { getData } from './func.js';
  import { authToken } from '../../state';

  import Header from '../../Shared/Components/Header.svelte';
  import WireguardInterface from './WireguardInterface.svelte';
  import TrustedIPSet from './TrustedIPSet.svelte';

  export let cfg;
  let client = cfg.client;

  setHtmlClass('h-full', 'bg-gray-100');
  setBodyClass('h-full');

  // state
  let session = get(authToken);
  let wgCfg = {
    wanIP: '',
    wgInet: '',
    wgPort: '',
    wgPubKey: ''
  };
  let ipSet = [];

  let isDisabled = true;
  let isLoading = true;

  onMount(() => {
    getData(client, session)
      .then(result => {
        wgCfg = {
          wanIP: result['wgCfg']['wanip'],
          wgInet: result['wgCfg']['wg_inet'],
          wgPort: result['wgCfg']['wg_port'],
          wgPubKey: result['wgCfg']['wg_pubkey']
        };
        ipSet = result['ipSet'];
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
      })

    return () => {
    }
  });
</script>

<svelte:head>
  <title>server</title>
</svelte:head>

<Header {cfg} />

<WireguardInterface {wgCfg} {isLoading} />
<TrustedIPSet {cfg} {session} {ipSet} {isLoading} {isDisabled} />
