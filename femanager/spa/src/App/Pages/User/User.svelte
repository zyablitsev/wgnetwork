<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { setBodyClass, setHtmlClass } from '../../util/dom';
  import { ValidationError, RPCError } from '../../../lib/rpcapi';
  import { getUserWithDevices } from './func.js';
  import { authToken, moveTo } from '../../state';

  import Header from '../../Shared/Components/Header.svelte';

  import UserInformation from './Components/UserInformation.svelte';
  import UserManagerInformation from './Components/UserManagerInformation.svelte';
  import UserDevices from './Components/UserDevices.svelte';
  import UserActionButton from './Components/UserActionButton.svelte';
  import UserModalRemove from './Components/UserModalRemove.svelte';

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

  let user = {};
  let title = 'user';

  function handleActionButton(event) {
    event.preventDefault;
    if (event.detail.action === 'remove') {
      showModal = true;
    } else if (event.detail.action === 'edit') {
      let path = router.ReverseURI('user_edit', {'uuid': user.uuid});
      moveTo(path);
    }
  }

  function handleModalButton(event) {
    event.preventDefault;
    if (event.detail.action === 'remove') {
      removeUser();
    } else if (event.detail.action === 'cancel') {
      showModal = false;
    }
  }

  function removeUser() {
    isLoading = true;
    isDisabled = true;

    let params = {'uuid': user.uuid};
    client.Fetch('manager/user/remove', params, session)
      .then(result => {
        isLoading = false;
        isDisabled = false;
        let path = router.ReverseURI('users');
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
    getUserWithDevices(client, session, params['uuid'])
      .then(result => {
        user = result['user'];
        title = 'user - ' + user['name'];
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

  <UserInformation {user} {isLoading} />
  <UserActionButton {isDisabled} on:message={handleActionButton} />
  {#if user.isManager}
  <UserManagerInformation {user} {isLoading} />
  {/if}
  <UserDevices {cfg} devices={user.devices} userUUID={user.uuid} {isLoading} />

  {#if showModal}
    <UserModalRemove {isDisabled} on:message={handleModalButton} />
  {/if}

</div>
