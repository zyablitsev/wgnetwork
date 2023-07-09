<script>
  import { afterUpdate, onMount } from 'svelte';

  import { copyToClipboard } from './func.js';
  import { moveBack } from '../../../state';

  import CardHeading from '../../../Shared/Components/CardHeading.svelte';
  import Link from '../../../Shared/Components/Link.svelte';
  import LinkButton from '../../../Shared/Components/Button/LinkButton.svelte';
  import DeviceInformationRow from './DeviceInformationRow.svelte';

  export let cfg;
  let router = cfg.router;

  export let device = {
    ip: '',
    label: '',
    wanForward: false,
    user: {
      uuid: '',
      name: ''
    }
  };
  export let isLoading = false;

  onMount(() => {
    let elBtnBack = document.getElementById('btn_back');
    elBtnBack.addEventListener('click', moveBack);

    return () => {
      elBtnBack.removeEventListener('click', moveBack);
    }
  });

  let cbElems = [];
  afterUpdate(() => {
    if (isLoading) {
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].removeEventListener('click', copyToClipboard);
      }
    } else {
      cbElems = document.querySelectorAll('[data-cb="deviceinfo"]');
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].addEventListener('click', copyToClipboard);
      }
    }
  });
</script>

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading} title='device information' description='main parameters.'>
    <LinkButton id='btn_back' cssLayout='float-right'>back</LinkButton>
  </CardHeading>

  <div class="p-0">
    <dl class="divide-y divide-gray-200">
      <DeviceInformationRow key='ip' value={device.ip} {isLoading} clipboard='deviceinfo' />
      <DeviceInformationRow key='label' value={device.label} {isLoading} clipboard='deviceinfo' />
      <DeviceInformationRow key='wan forward' value={device.wanForward ? 'on' : 'off'} {isLoading} />
      <DeviceInformationRow key='user' {isLoading}>
        <Link href="{router.ReverseURI('user', {'uuid': device.user.uuid})}"
              css=''>
        {device.user.name} ({device.user.uuid})
        </Link>
      </DeviceInformationRow>
    </dl>
  </div>

</div>

