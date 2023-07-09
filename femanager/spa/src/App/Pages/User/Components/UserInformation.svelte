<script>
  import { afterUpdate, onMount } from 'svelte';

  import { copyToClipboard } from './func.js';
  import { moveBack } from '../../../state';

  import CardHeading from '../../../Shared/Components/CardHeading.svelte';
  import LinkButton from '../../../Shared/Components/Button/LinkButton.svelte';

  import UserInformationRow from './UserInformationRow.svelte';

  export let user = {
    uuid: '',
    name: '',
    role: ''
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
      cbElems = document.querySelectorAll('[data-cb="userinfo"]');
      for (let i = 0; i < cbElems.length; i++) {
        cbElems[i].addEventListener('click', copyToClipboard);
      }
    }
  });
</script>

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading} title='user information' description='main attributes.'>
    <LinkButton id='btn_back' cssLayout='float-right'>back</LinkButton>
  </CardHeading>

  <div class="p-0">
    <dl class="divide-y divide-gray-200">
      <UserInformationRow key='uuid' value={user.uuid} {isLoading} clipboard='userinfo' />
      <UserInformationRow key='name' value={user.name} {isLoading} clipboard='userinfo' />
      <UserInformationRow key='role' value={user.role} {isLoading} clipboard='userinfo' />
    </dl>
  </div>

</div>
