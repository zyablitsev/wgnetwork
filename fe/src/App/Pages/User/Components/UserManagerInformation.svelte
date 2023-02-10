<script>
  import { afterUpdate } from 'svelte';

  import { copyToClipboard } from './func.js';

  import CardHeading from '../../../Shared/Components/CardHeading.svelte';
  import QRCode from '../../../Shared/Components/QRCode.svelte';

  import UserInformationRow from './UserInformationRow.svelte';

  export let user = {
    key: '',
    provisionUri: '',
  };
  export let isLoading = false;

  let cbElems = [];
  afterUpdate(() => {
    if (isLoading) {
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].removeEventListener('click', copyToClipboard);
      }
    } else {
      cbElems = document.querySelectorAll('[data-cb="authinfo"]');
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].addEventListener('click', copyToClipboard);
      }
    }
  });
</script>

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading}
    title='management panel authentication information'
    description='totp secret' />

  <div class="p-0">
    <dl class="divide-y divide-gray-200">
      <UserInformationRow key='totp key' value={user.key} {isLoading} clipboard='authinfo' />
      <UserInformationRow key='provision uri' value={user.provisionUri} {isLoading} clipboard='authinfo' />
      <UserInformationRow key='qr code' {isLoading}>
        {#if isLoading}
          <img id="qrcode" alt="" width="250" height="250" src="https://media.istockphoto.com/id/1088618726/photo/mans-hands-working-on-business-data-to-discuss-information-with-calculator-laptop-cup-of.jpg?s=170667a&amp;w=0&amp;k=20&amp;c=f9vAcfHCDETwzJ7c9w3QMe9BITuYDyOrMWjZMMulkG4=" title="Man's hands working on business data to discuss information with calculator,laptop,cup of coffee,mouse,book on modern white desk table at office." class="tB6UZ a5VGX" loading="lazy">
        {:else}
          <QRCode id='qrcode' text={user['provisionUri']}>
            <div id="qrcode"></div>
          </QRCode>
          <p class="mt-2 text-sm text-gray-500">scan with authenticator app</p>
        {/if}
      </UserInformationRow>
    </dl>
  </div>

</div>
